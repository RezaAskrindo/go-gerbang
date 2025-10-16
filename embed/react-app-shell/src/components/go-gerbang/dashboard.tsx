import { lazy, Suspense, useState, type FC } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import useSWR from "swr";
import useSWRMutation from "swr/mutation";
import { toast } from "sonner";

import { 
  CalendarDays, 
  Check, 
  ChevronRight, 
  Clock2, 
  RefreshCcw, 
  Search, 
  Timer, 
  Trash2, 
  User, 
  X 
} from "lucide-react";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs"
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Input } from "../ui/input";

import DataTable from "@/components/data-table";

import { 
  BackendUrlBase, 
  fetchSWR,
  SWRDashboardConfig
} from "@/services/baseService";

import TableSkeleton from "./table-skeleton";

const CodeHighlighter = lazy(() => import("@/components/CodeHighlighter"));
const MetricsInfo = lazy(() => import("./metrics-info"));

type InfoDataType = {
  auth_protection: boolean
  csrf_protection: boolean
  session_protection: boolean
  path: string
  rbac_protection: boolean
  service: string
  status?: boolean
  url: string
}

const CircuitBreaker = () => {
  const { data: circuitData } = useSWR(`${BackendUrlBase}/metrics/circuit`, fetchSWR, SWRDashboardConfig);

  const lastUpdate = circuitData?.lastStateChange ? new Date(circuitData?.lastStateChange).toLocaleString() : null;

  return (
    <Card className="gap-3">
      <CardHeader>
        <CardDescription>TOTAL Request</CardDescription>
        <CardTitle className="text-2xl font-bold tabular-nums @[250px]/card:text-3xl">
          {circuitData?.totalRequests?.toLocaleString()}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-2">
        <div className="flex items-center justify-between">
          <div>Success:</div>
          <div className="font-bold">{circuitData?.successThreshold}</div>
        </div>
        <div className="flex items-center justify-between">
          <div>Success On Failed:</div>
          <div>{circuitData?.successes}</div>
        </div>
        <div className="flex items-center justify-between">
          <div>Failures:</div>
          <div className="font-bold">{circuitData?.failures}</div>
        </div>
        <div className="flex items-center justify-between">
          <div>Failures Threshold:</div>
          <div>{circuitData?.failureThreshold}</div>
        </div>
        <div className="flex items-center justify-between">
          <div>Rejected:</div>
          <div className="font-bold">{circuitData?.rejectedRequests}</div>
        </div>
      </CardContent>
      <CardFooter className="mx-auto mt-2">
        <div className="text-muted-foreground font-semibold">Last Update: {lastUpdate}</div>
      </CardFooter>
    </Card>
  )
}

const ServicesInfo = () => {
  const [edit, setEdit] = useState(false);
  const [listInfoData, setListInfoData] = useState<InfoDataType[]>([]);

  const { data: infoData } = useSWR(`${BackendUrlBase}/info`, fetchSWR, SWRDashboardConfig);

  // useEffect(() => {
  //   if (infoData) {
  //     setListInfoData(infoData);
  //   }
  // }, [infoData]);

  // if (typeof window === "undefined") {
  //   return null;
  // }

  const handleProtectionChange = (field: string, val: boolean | string, index: number) => {
    const updatedData = [...listInfoData];
    updatedData[index] = { ...updatedData[index], [field]: val };
    setListInfoData(updatedData);
  };

  const handleRestart = async () => {
    toast.promise(
      fetch(`${BackendUrlBase}/restart`, {
        method: "POST",
        credentials: "include",
      }).then(async (res) => {
        if (!res.ok) throw new Error("Request failed")
        const data = await res.json()
        if (!data.status) throw new Error(data.message || "Failed to login")
        return data
      }),
      {
        loading: "Waiting...",
        success: () => {
          return "Success"
        },
        error: (err) => {
          return err.message || "Failed"
        },
      }
    )
  }

  const handleUpdate = () => {
    // listInfoData.forEach(obj => {
    //   delete obj.status;
    // });
    toast.promise(
      fetch(`${BackendUrlBase}/config-file`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({services: listInfoData}),
      }).then(async (res) => {
        if (!res.ok) throw new Error("Request failed")
        const data = await res.json()
        if (!data.status) throw new Error(data.message || "Failed to login")
        return data
      }),
      {
        loading: "Waiting...",
        success: () => {
          // mutate(`${BackendUrlBase}/info`);
          return "Success"
        },
        error: (err) => {
          return err.message || "Failed"
        },
      }
    )
  }

  return (
    <Card className="gap-2">
      <CardHeader>
        <CardTitle>Service Info</CardTitle>
        <CardDescription>List of service gateway</CardDescription>
        <CardAction>
          <div className="flex items-center space-x-2">
            <Switch id="edit-mode" onCheckedChange={setEdit} checked={edit} />
            <Label htmlFor="edit-mode">Edit Mode</Label>
          </div>
        </CardAction>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead rowSpan={2} className="w-[50px]">Status</TableHead>
              <TableHead rowSpan={2}>Service</TableHead>
              <TableHead rowSpan={2}>Path</TableHead>
              <TableHead rowSpan={2}>URL</TableHead>
              <TableHead colSpan={4} className="text-center">Middleware</TableHead>
            </TableRow>
            <TableRow>
              <TableHead className="w-[100px] text-center">CSRF</TableHead>
              <TableHead className="w-[100px] text-center">Auth</TableHead>
              <TableHead className="w-[100px] text-center">Session</TableHead>
              <TableHead className="w-[100px] text-center">RBAC</TableHead>
            </TableRow>
          </TableHeader>

          <TableBody>
            {infoData?.length > 0 ? (
              infoData?.map((el: InfoDataType, index: number) => (
                <TableRow key={`${el.url}-${index}`}>
                  <TableCell>
                    {edit ? (
                      <Button variant="destructive" size="icon" onClick={() => setListInfoData((prevItems) => prevItems.filter((_, i) => i !== index))}>
                        <Trash2 />
                      </Button>
                    ) : (
                      el.status ? 
                        <Button className="bg-green-600 h-5" size="sm">live</Button> : 
                        <Button className="bg-red-500 h-5" size="sm">off</Button>
                      )
                    }
                  </TableCell>
                  <TableCell>{
                    edit ? 
                    (<Input value={el.service} onChange={(e) => handleProtectionChange('service', e.target.value, index)} />):
                    el.service
                  }</TableCell>
                  <TableCell>{
                    edit ? 
                    (<Input value={el.path} onChange={(e) => handleProtectionChange('path', e.target.value, index)} />):
                    <a href={`${BackendUrlBase}${el.path}`} className="underline text-blue-700 dark:text-blue-400" target="_blank"> {el.path}</a>
                  }</TableCell>
                  <TableCell>{
                    edit ? 
                    (<Input value={el.url} onChange={(e) => handleProtectionChange('url', e.target.value, index)} />):
                    el.url
                  }</TableCell>
                  <TableCell>
                    {edit ? (
                      <div className="flex justify-center"><Switch
                        id={`csrf-toggle-${index}`}
                        checked={el.csrf_protection}
                        onCheckedChange={(checked: boolean) => handleProtectionChange('csrf_protection', checked, index)}
                      /></div>
                    ) : (
                      el.csrf_protection ? <Check className="text-green-500 mx-auto size-5" /> : <X className="text-red-500 mx-auto size-5" />
                    )}
                  </TableCell>
                  <TableCell>
                    {edit ? (
                      <div className="flex justify-center"><Switch
                        id={`auth-toggle-${index}`}
                        checked={el.auth_protection}
                        onCheckedChange={(checked: boolean) => handleProtectionChange('auth_protection', checked, index)}
                      /></div>
                    ) : (
                      el.auth_protection ? <Check className="text-green-500 mx-auto size-5" /> : <X className="text-red-500 mx-auto size-5" />
                    )}
                  </TableCell>
                  <TableCell>
                    {edit ? (
                      <div className="flex justify-center"><Switch
                        id={`session-toggle-${index}`}
                        checked={el.session_protection}
                        onCheckedChange={(checked: boolean) => handleProtectionChange('session_protection', checked, index)}
                      /></div>
                    ) : (
                      el.session_protection ? <Check className="text-green-500 mx-auto size-5" /> : <X className="text-red-500 mx-auto size-5" />
                    )}
                  </TableCell>
                  <TableCell>
                    {edit ? (
                      <div className="flex justify-center"><Switch
                        id={`rbac-toggle-${index}`}
                        checked={el.rbac_protection}
                        onCheckedChange={(checked: boolean) => handleProtectionChange('rbac_protection', checked, index)}
                      /></div>
                    ) : (
                      el.rbac_protection ? <Check className="text-green-500 mx-auto size-5" /> : <X className="text-red-500 mx-auto size-5" />
                    )}
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={8} className="text-center">No Data Available</TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
        {edit && (
          <div className="flex flex-row gap-2 mt-2">
            <Button variant="outline" onClick={() => setListInfoData((prevItems) => [...prevItems, {
              auth_protection: false,
              csrf_protection: false,
              session_protection: false,
              path: "",
              rbac_protection: false,
              service: "",
              status: false,
              url: "",
            }])}>Add</Button>
            <Button variant="outline" onClick={handleRestart}>Restart</Button>
            <Button variant="outline" onClick={handleUpdate}>Save</Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

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
            </DialogHeader>
            <Tabs defaultValue="request">
              <TabsList>
                <TabsTrigger value="request">Request Log</TabsTrigger>
                <TabsTrigger value="response">Response Log</TabsTrigger>
              </TabsList>
              <TabsContent value="request">
                {raw.fields?.request ? <pre
                  className="no-scrollbar min-w-0 overflow-x-auto px-4 py-3.5 outline-none has-[[data-highlighted-line]]:px-0 has-[[data-line-numbers]]:px-0 has-[[data-slot=tabs]]:p-0"
                >
                  <Suspense fallback={<p>Loading...</p>}>
                    <CodeHighlighter code={raw.fields?.request} />
                  </Suspense>
                </pre> : "no request body"}
              </TabsContent>
              <TabsContent value="response">
                {raw.fields?.response ? <pre
                  className="no-scrollbar min-w-0 overflow-x-auto px-4 py-3.5 outline-none has-[[data-highlighted-line]]:px-0 has-[[data-line-numbers]]:px-0 has-[[data-slot=tabs]]:p-0"
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

const LogProxyData = () => {
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
              <div className="border-b-1">
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
              <div className="border-b-1">
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
  )
}

export default function Dashboard() {

  return (
    <div className="grid grid-cols-4 gap-4">
      <div className="col-span-4">
        <Suspense fallback={<p>Loading...</p>}>
          <MetricsInfo />
        </Suspense>
      </div>
      <div className="col-span-4 sm:col-span-1">
        <CircuitBreaker />
      </div>
      <div className="col-span-4 sm:col-span-3">
        <ServicesInfo />
      </div>
      <div className="col-span-4">
        <LogProxyData />
      </div>
    </div>
  )
}