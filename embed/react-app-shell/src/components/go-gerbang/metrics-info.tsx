import { useEffect, useState, startTransition } from "react";
import useSWR from "swr";

import Highcharts from 'highcharts';

Highcharts.setOptions({
  xAxis: {
    type: "datetime"
  },
  yAxis: {
    title: {
      text: undefined
    }
  },
  time: {
    timezone: 'Asia/Jakarta'
  },
  chart: {
    height: 200,
    width: null,
  },
  credits: {
    enabled: false
  },
  legend: {
    enabled: false
  }
});

import { Chart, setHighcharts } from '@highcharts/react';
import { Area } from '@highcharts/react/series';

import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { BackendUrlBase, SWRDashboardConfig } from "@/services/baseService";

const MAX_POINTS = 100;

function formatBytes(bytes: number, decimals = 1): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function limitData(data: any): any {
  if (Array.isArray(data) && data.length > MAX_POINTS) {
    return data.slice(data.length - MAX_POINTS);
  }
  return data;
}

export default function MetricsInfo() {
  setHighcharts(Highcharts);

  const [cpuSeries, setCpuSeries] = useState<[number, number][]>([]);
  const [ramSeries, setRamSeries] = useState<{name: string, data: [number, number][]}[]>([
    {name: "ram usage", data: []},
    {name: "ram in use", data: []},
    {name: "total ram", data: []},
  ]);
  const [rtimeSeries, setRtimeSeries] = useState<[number, number][]>([]);
  const [connsSeries, setConnsSeries] = useState<[number, number][]>([]);
  const [rtime, setRtime] = useState<number>(0);

  const { data: metricsData } = useSWR(`${BackendUrlBase}/metrics`, async (url: string) => {
    const t0 = performance.now();
    const res = await fetch(url, {
      headers: { Accept: 'application/json' },
      credentials: 'same-origin'
    });
    const json = await res.json();
    const t1 = performance.now();
    setRtime(Math.round(t1 - t0));
    return json;
  }, SWRDashboardConfig);

  useEffect(() => {
    if (!metricsData) return;

    const time = Date.now();

    startTransition(() => {
      setCpuSeries(prev => limitData([...prev, [time, Number(metricsData.pid.cpu.toFixed(1))]]));

      const ramUsage = Number((metricsData.pid.ram / 1e6).toFixed(2));
      const ramIdle = Number((metricsData.os.ram / 1e6).toFixed(2));
      const ramTotal = Number((metricsData.os.total_ram / 1e6).toFixed(2));
  
      setRamSeries(prev => [
        {
          ...prev[0],
          data: limitData([...prev[0].data, [time, ramUsage]]),
        },
        {
          ...prev[1],
          data: limitData([...prev[1].data, [time, ramIdle]]),
        },
        {
          ...prev[2],
          data: limitData([...prev[2].data, [time, ramTotal]]),
        },
      ]);

      setRtimeSeries(prev => limitData([...prev, [time, rtime]]));

      setConnsSeries(prev => limitData([...prev, [time, metricsData.pid.conns]]));
    });

  }, [metricsData, rtime]);

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Memory Usage</CardDescription>
          <CardTitle className="text-2xl font-bold tabular-nums @[250px]/card:text-3xl">
            {formatBytes(metricsData?.pid?.ram)} <span className="border-s ps-1 text-sm font-semibold text-orange-400">{formatBytes(metricsData?.os?.ram)}</span> <span className="border-s ps-1 text-sm font-semibold text-red-500">{formatBytes(metricsData?.os?.total_ram)}</span>
          </CardTitle>
        </CardHeader>
        <CardFooter className="px-0">
          <Chart>
            <Area.Series data={ramSeries[0].data} />
            <Area.Series data={ramSeries[1].data} />
            <Area.Series data={ramSeries[2].data} />
          </Chart>
        </CardFooter>
      </Card>
      <Card className="@container/card pb-1">
        <CardHeader>
          <CardDescription>CPU Usage</CardDescription>
          <CardTitle className="text-2xl font-bold tabular-nums @[250px]/card:text-3xl">
            {metricsData?.pid?.cpu.toFixed(1)}% <span className="border-s ps-1 text-sm font-normal">{metricsData?.os?.cpu.toFixed(1)}%</span>
          </CardTitle>
        </CardHeader>
        <CardFooter className="px-0">
          <Chart>
            <Area.Series data={cpuSeries} />
          </Chart>
        </CardFooter>
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Response Time</CardDescription>
          <CardTitle className="text-2xl font-bold tabular-nums @[250px]/card:text-3xl">
            {rtime}ms <span className="text-sm font-semibold">clients</span>
          </CardTitle>
        </CardHeader>
        <CardFooter className="px-0">
          <Chart>
            <Area.Series data={rtimeSeries} />
          </Chart>
        </CardFooter>
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Connections</CardDescription>
          <CardTitle className="text-2xl font-bold tabular-nums @[250px]/card:text-3xl">
            {metricsData?.pid?.conns} <span className="text-sm font-semibold">{metricsData?.os?.conns}</span>
          </CardTitle>
        </CardHeader>
        <CardFooter className="px-0">
          <Chart>
            <Area.Series data={connsSeries} />
          </Chart>
        </CardFooter>
      </Card>
    </div>
  )
}
