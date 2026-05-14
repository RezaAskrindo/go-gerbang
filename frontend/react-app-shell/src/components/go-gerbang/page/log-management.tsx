import { lazy, Suspense, useState, type FC } from "react";
import useSWR from "swr";
import useSWRMutation from "swr/mutation";
import { CalendarDays, ChevronRight, Clock2, RefreshCcw, Search, Timer, User } from "lucide-react";
import type { ColumnDef } from "@tanstack/react-table";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs"

import { BackendUrlBase } from "@/services/baseService";
import { fetchSWR } from "@/services/use-swr-service";

import TableSkeleton from "../table-skeleton";
import DataTable from "../../data-table";

const CodeHighlighter = lazy(() => import("@/components/CodeHighlighter"));

const methodListRequest = {
  "GET": "bg-green-600",
  "POST": "bg-blue-500",
  "PUT": "bg-orange-500",
  "DELETE": "bg-red-500",
}

const methodListStatus = {
  200: "bg-green-600",
  404: "bg-orange-500",
  502: "bg-red-500",
}

type TLogProxy = {
  avg_duration: number
  method: string
  path: string
  request_count: number
  service: string
  status: number
} 

type TLogProxyDetail = {
  level: string
  service: string
  method: string
  path: string
  status: number
  duration: number
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  fields: Record<string, any>;
  timestamp: string
  user_auth: string | null
}

type DetailLogProxyProps = {
  row: TLogProxy
  dateBegin?: Date
  dateEnd?: Date
}

const DetailLogProxy: FC<DetailLogProxyProps> = ({ row, dateBegin, dateEnd }) => {
  const { isLoading, data: logProxyDataDetail, error } = useSWR(dateBegin && dateEnd ? `${BackendUrlBase}/log-stats-proxy?detail=true&service=${row.service}&method=${row.method}&path=${row.path}&status=${row.status}&from=${dateBegin.toISOString()}&to=${dateEnd.toISOString()}` : `${BackendUrlBase}/log-stats-proxy`, fetchSWR);

  const columnsDetail: ColumnDef<TLogProxyDetail>[] = [
    {
      accessorKey: "timestamp",
      header: "Time",
      cell: ({row}) => {
        const raw = row.original;
        const date = new Date(raw.timestamp);
        const now = new Date();

        const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const target = new Date(date.getFullYear(), date.getMonth(), date.getDate());

        const diffMs = target.getTime() - today.getTime();
        const diffDays = Math.round(diffMs / (1000 * 60 * 60 * 24));

        const rtf = new Intl.RelativeTimeFormat("en", { style: "short" });

        return <div className="flex flex-col justify-center">
          <div className="font-bold">{date.toLocaleTimeString()}</div>
          <div className="text-muted-foreground">{diffDays === 0 ? "today" : rtf.format(diffDays, "day")}</div>
        </div>
      }
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const raw = row.original;

        return <Button key={raw.path} className={`${methodListStatus[raw.status as keyof typeof methodListStatus] ?? ''} hover:${methodListStatus[raw.status as keyof typeof methodListStatus] ?? ''} h-5`} size="sm">{raw.status}</Button>
      }
    },
    {
      accessorKey: "method",
      header: "Endpoint",
      cell: ({ row }) => {
        const raw = row.original;

        return <div className="flex flex-row items-center gap-2">
          <Button className={`${methodListRequest[raw.method as keyof typeof methodListRequest] ?? ''} hover:${methodListRequest[raw.method as keyof typeof methodListRequest] ?? ''} h-5`} size="sm">{raw.method}</Button>
          <div className="flex flex-col gap-1">
            <div className="font-semibold">{raw.path}</div>
            <div className="flex flex-row gap-2">
              {raw.user_auth && <div className="flex flex-row items-center gap-1 text-muted-foreground"><User className="size-4" />{raw.user_auth}</div>}
              <div className="flex flex-row items-center gap-1 text-muted-foreground"><Timer className="size-4" />{raw.duration} ms</div>
            </div>
          </div>
        </div>
      }
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const raw = row.original;

        return <Dialog>
          <DialogTrigger asChild>
            <Button variant="outline" size="icon">
              <Search />
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-5xl">
            <DialogHeader>
              <DialogTitle>Detail Logs</DialogTitle>
              <DialogDescription className="sr-only">detail of logs</DialogDescription>
            </DialogHeader>
            <Tabs defaultValue="request">
              <TabsList>
                <TabsTrigger value="request">Request Log</TabsTrigger>
                <TabsTrigger value="response">Response Log</TabsTrigger>
              </TabsList>
              <TabsContent value="request">
                {raw.fields?.request ? <pre
                  className="no-scrollbar min-w-0 overflow-x-auto px-4 py-3.5 outline-none has-data-highlighted-line:px-0 has-data-line-numbers:px-0 has-data-[slot=tabs]:p-0"
                >
                  <Suspense fallback={<p>Loading...</p>}>
                    <CodeHighlighter code={raw.fields?.request} />
                  </Suspense>
                </pre> : "no request body"}
              </TabsContent>
              <TabsContent value="response">
                {raw.fields?.response ? <pre
                  className="no-scrollbar min-w-0 overflow-x-auto px-4 py-3.5 outline-none has-data-highlighted-line:px-0 has-data-line-numbers:px-0 has-data-[slot=tabs]:p-0"
                >
                  <Suspense fallback={<p>Loading...</p>}>
                    <CodeHighlighter code={raw.fields?.response} />
                  </Suspense>
                </pre> : "no response body"}
              </TabsContent>
            </Tabs>
          </DialogContent>
        </Dialog>
      }
    },
  ];

  if (isLoading) return <TableSkeleton />;
  if (error) return <div>Error loading details</div>;

  return (
    <div>
      <DataTable rowIdKey="timestamp" columns={columnsDetail} data={logProxyDataDetail?.data ?? []} />
    </div>
  )
}

const columns: ColumnDef<TLogProxy>[] = [
  {
    id: "service",
    header: "Service",
    cell: ({ row }) => <div className="flex flex-row items-center">
      <Button onClick={row.getToggleExpandedHandler()} variant="ghost" size="icon" className="data-[state=open]:bg-muted text-muted-foreground flex size-7"><ChevronRight /></Button>
      <div className="font-medium">{row.original.service}</div>
    </div>,
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const raw = row.original;
      return <Button key={raw.path} className={`${methodListStatus[raw.status as keyof typeof methodListStatus] ?? ''} hover:${methodListStatus[raw.status as keyof typeof methodListStatus] ?? ''} h-5`} size="sm">{raw.status}</Button>
    }
  },
  {
    accessorKey: "method",
    header: "Endpoint",
    cell: ({ row }) => {
      const raw = row.original;

      return <div className="flex flex-row gap-1">
        <Button className={`${methodListRequest[raw.method as keyof typeof methodListRequest] ?? ''} hover:${methodListRequest[raw.method as keyof typeof methodListRequest] ?? ''} h-5`} size="sm">{raw.method}</Button>
        <div>{raw.path}</div>
      </div>
    }
  },
  {
    accessorKey: "request_count",
    header: "Request",
  },
  {
    accessorKey: "avg_duration",
    header: "Duration",
    cell: ({ row }) => row.original.avg_duration.toFixed(2)
  },
];


export default function LogManagement() {
  const[openBegin, setOpenBegin] = useState(false); 
  const[dateBegin, setDateBegin] = useState<Date | undefined>(undefined);
  const[openEnd, setOpenEnd] = useState(false); 
  const[dateEnd, setDateEnd] = useState<Date | undefined>(undefined);

  const handleTimeBeginChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const [h, m, s] = value.split(":").map(Number);
    if (!dateBegin) return;
    const newDate = new Date(dateBegin);
    newDate.setHours(h || 0, m || 0, s || 0);
    setDateBegin(newDate);
  };

  const handleTimeEndChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const [h, m, s] = value.split(":").map(Number);
    if (!dateEnd) return;
    const newDate = new Date(dateEnd);
    newDate.setHours(h || 0, m || 0, s || 0);
    setDateEnd(newDate);
  };

  const { trigger, data: logProxyData, isMutating } = useSWRMutation(dateBegin && dateEnd ? `${BackendUrlBase}/log-stats-proxy?from=${dateBegin.toISOString()}&to=${dateEnd.toISOString()}` : `${BackendUrlBase}/log-stats-proxy`, fetchSWR);

  return (
    <div className="flex-1 flex-col gap-8 md:flex">
      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-col gap-1">
          <h2 className="text-2xl font-semibold tracking-tight">
            Log Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s a list of log of proxy requests that have been made.
          </p>
        </div>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex flex-row gap-2">
            <Label htmlFor="time-from">Date Start</Label>
            <Popover open={openBegin} onOpenChange={setOpenBegin}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  id="date"
                  className="w-48 justify-between font-normal"
                >
                  {dateBegin ? dateBegin.toLocaleTimeString('ID', {year: "numeric", month: "2-digit", day: "2-digit"}) : "Select date"}
                  <CalendarDays />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto overflow-hidden p-0" align="start">
                <div className="border-b">
                  <Calendar
                    mode="single"
                    selected={dateBegin}
                    captionLayout="dropdown"
                    disabled={{ after: dateEnd ?? new Date() }}
                    onSelect={(date) => {
                      setDateBegin(date)
                      setOpenBegin(false)
                    }}
                  />
                </div>
                <div className="px-5 py-3">
                  <div className="flex w-full flex-col gap-3">
                    <Label htmlFor="time-from">Time (HH.MM.SS)</Label>
                    <div className="relative flex w-full items-center gap-2">
                      <Clock2 className="text-muted-foreground pointer-events-none absolute left-2.5 size-4 select-none" />
                      <Input
                        id="time-from"
                        type="time"
                        step="1"
                        value={dateBegin ? dateBegin.toTimeString().slice(0, 8) : ""}
                        onChange={handleTimeBeginChange}
                        className="appearance-none pl-8 [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                      />
                    </div>
                  </div>
                </div>
              </PopoverContent>
            </Popover>
          </div>
          <div className="flex flex-row gap-2">
            <Label htmlFor="time-from">Date End</Label>
            <Popover open={openEnd} onOpenChange={setOpenEnd}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  id="date"
                  className="w-48 justify-between font-normal"
                >
                  {dateEnd ? dateEnd.toLocaleTimeString('ID', {year: "numeric", month: "2-digit", day: "2-digit"}) : "Select date"}
                  <CalendarDays />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto overflow-hidden p-0" align="start">
                <div className="border-b">
                  <Calendar
                    mode="single"
                    selected={dateEnd}
                    captionLayout="dropdown"
                    disabled={{ before: dateBegin ?? new Date() }}
                    onSelect={(date) => {
                      setDateEnd(date)
                      setOpenEnd(false)
                    }}
                  />
                </div>
                <div className="px-5 py-3">
                  <div className="flex w-full flex-col gap-3">
                    <Label htmlFor="time-from">Time (HH.MM.SS)</Label>
                    <div className="relative flex w-full items-center gap-2">
                      <Clock2 className="text-muted-foreground pointer-events-none absolute left-2.5 size-4 select-none" />
                      <Input
                        id="time-from"
                        type="time"
                        step="1"
                        value={dateEnd ? dateEnd.toTimeString().slice(0, 8) : ""}
                        onChange={handleTimeEndChange}
                        className="appearance-none pl-8 [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                      />
                    </div>
                  </div>
                </div>
              </PopoverContent>
            </Popover>
          </div>
          <Button type="button" variant="outline" onClick={() => trigger()}>
            <RefreshCcw />
            Reload
          </Button>
        </div>
        {isMutating ? <TableSkeleton /> : <DataTable rowIdKey="path" columns={columns} data={logProxyData?.data ?? []} canExpand={true} DetailRow={DetailLogProxy} detailRowProps={{ dateBegin, dateEnd }} />}
      </div>
    </div>
  )
}