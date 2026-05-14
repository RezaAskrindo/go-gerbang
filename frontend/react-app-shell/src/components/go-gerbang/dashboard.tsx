import { lazy, Suspense, useState } from "react";
import useSWR from "swr";
import { toast } from "sonner";

import { 
  Check, 
  RefreshCcw, 
  Trash2, 
  X
} from "lucide-react";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Switch } from "@/components/ui/switch";

import { BackendUrlBase } from "@/services/baseService";

import MetricsInfoSkeleton from "./metrics-info-skeleton";
import { fetchSWR, SWRDashboardConfig } from "@/services/use-swr-service";

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

  const [open, setOpen] = useState(false);

  const lastUpdate = circuitData?.lastStateChange ? new Date(circuitData?.lastStateChange).toLocaleString() : null;

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
          setTimeout(() => {
            window.location.reload();
          }, 5000);
          return "Success"
        },
        error: (err) => {
          return err.message || "Failed"
        },
      }
    )
  }

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
      <CardFooter className="flex flex-col gap-2">
        <div className="text-muted-foreground font-semibold">Last Update: {lastUpdate}</div>
        <AlertDialog open={open} onOpenChange={setOpen}>
          <AlertDialogTrigger asChild>
            <Button className="w-full" variant="destructive">
              <RefreshCcw />
              Restart GATEWAY
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action will restart the gateway services. All services will stop temporarily and may cause downtime.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogAction>Cancel</AlertDialogAction>
              <AlertDialogCancel onClick={handleRestart}>Agree</AlertDialogCancel>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </CardFooter>
    </Card>
  )
}

const ServicesInfo = () => {
  const [edit, setEdit] = useState(false);
  const [listInfoData, setListInfoData] = useState<InfoDataType[]>([]);

  const { data: infoData } = useSWR(`${BackendUrlBase}/info`, fetchSWR, SWRDashboardConfig);

  const handleProtectionChange = (field: string, val: boolean | string, index: number) => {
    const updatedData = [...listInfoData];
    updatedData[index] = { ...updatedData[index], [field]: val };
    setListInfoData(updatedData);
  };

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
              <TableHead rowSpan={2} className="w-12.5">Status</TableHead>
              <TableHead rowSpan={2}>Service</TableHead>
              <TableHead rowSpan={2}>Path</TableHead>
              <TableHead rowSpan={2}>URL</TableHead>
              <TableHead colSpan={4} className="text-center">Middleware</TableHead>
            </TableRow>
            <TableRow>
              <TableHead className="w-25 text-center">CSRF</TableHead>
              <TableHead className="w-25 text-center">Auth</TableHead>
              <TableHead className="w-25 text-center">Session</TableHead>
              <TableHead className="w-25 text-center">RBAC</TableHead>
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
            <Button variant="outline" onClick={handleUpdate}>Save</Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

export default function Dashboard() {

  return (
    <div className="grid grid-cols-4 gap-4">
      <div className="col-span-4">
        <Suspense fallback={<MetricsInfoSkeleton />}>
          <MetricsInfo />
        </Suspense>
      </div>
      <div className="col-span-4 sm:col-span-1">
        <CircuitBreaker />
      </div>
      <div className="col-span-4 sm:col-span-3">
        <ServicesInfo />
      </div>
    </div>
  )
}