import { useState, useEffect } from "react";
import { Plus, Send, Ban, Save, EllipsisVertical } from "lucide-react";
import useSWR from "swr";

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import type { ColumnDef } from "@tanstack/react-table";
import { toast } from "sonner";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button";
import { ButtonGroup, ButtonGroupSeparator } from "@/components/ui/button-group"
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { SheetClose } from "@/components/ui/sheet";
import { Spinner } from "@/components/ui/spinner";

import SheetForm from "@/components/sheet-form";
import CardInformation from "@/components/card-information";

import { fetchSWR } from "@/services/use-swr-service";
import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService";

import { cn } from "@/lib/utils";
import { ENV_STATUS } from "@/lib/constants";
import { 
  reverseTransformCaddyConfig, 
  transformCaddyConfig, 
  type FlatCaddyConfig 
} from "@/lib/caddy";

import { ComboboxMultipleCreatable } from "../combobox-multiple-creatable";


const formSchema = z.object({
  type: z.string(),
  // match_host: z.string().min(1, {message: "Match Host is required"}),
  match_host: z.array(z.string().min(1, {message: "Match Host is required"})),
  rewrite_uri: z.string().optional(),
  rewrite_strip_path_prefix: z.string().optional(),
  file_server: z.string().optional(),
  reverse_proxy: z.string().optional(),
  match_path: z.array(z.string()).optional(),
  host: z.string().optional(),
  i: z.number().optional(),
});

const defaultForm: Record<string, Partial<FlatCaddyConfig>> = {
  s_three_object: {
    type: "FE",
    rewrite_uri: "/mfe/react-app-shell/index.html",
    match_path: [
      "/mfe/react-app-shell/assets/*",
      "/mfe/react-app-shell/vite.svg"
    ],
    host: "s3.nevaobjects.id"
  },
  frontend: {
    type: "FE",
    rewrite_uri: "/index.html",
    file_server: "/Users/reza/Documents/web/golang/go-gerbang/frontend/react-app-shell/dist"
  },
  backend: {
    type: "BE",
    match_path: ["/api/*"],
    rewrite_strip_path_prefix: "/api",
    reverse_proxy: "localhost:8000"
  }
}

function SheetFormChild({
  openSheet,
  dataForm,
  setOpenSheet,
  setCaddyConfig,
}: {
  openSheet: boolean
  dataForm?: FlatCaddyConfig
  setOpenSheet?: (data: boolean) => void
  setCaddyConfig?: (data: FlatCaddyConfig) => void
}) {
  // console.log(dataForm?.type)
  
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: "",
      match_host: [],
      rewrite_uri: "",
      rewrite_strip_path_prefix: "",
      file_server: "",
      reverse_proxy: "",
      match_path: [],
      host: "",
    },
  });

  const [type] = form.watch(["type"]);
  // const type: any = "";

  const getDefaultForm = (key: keyof typeof defaultForm) => {
    const values = defaultForm[key];
    if (!values) return;

    form.reset({
      ...values,
    });
  };

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let cleanValue = Object.fromEntries(
      Object.entries(values).filter(([_, value]) => {
        if (Array.isArray(value)) {
          return value.length > 0;
        }
        if (typeof value === "string") {
          return value.trim() !== "";
        }
        return true;
      })
    );

    if (dataForm?.i !== undefined) {
      cleanValue = {
        ...cleanValue,
        i: dataForm.i
      }
    }

    setCaddyConfig?.(cleanValue as FlatCaddyConfig);
    setOpenSheet?.(false);
  }

  useEffect(() => {
    if (dataForm && openSheet) {
      if (dataForm.type) form.setValue("type", dataForm.type); // NOTE: NOT WORKING
      // form.setValue("type", "BE");

      if (dataForm.match_host) form.setValue("match_host", dataForm.match_host);
      if (dataForm.match_path) form.setValue("match_path", dataForm.match_path);
      if (dataForm.rewrite_uri) form.setValue("rewrite_uri", dataForm.rewrite_uri); 
      if (dataForm.rewrite_strip_path_prefix) form.setValue("rewrite_strip_path_prefix", dataForm.rewrite_strip_path_prefix); 
      if (dataForm.file_server) form.setValue("file_server", dataForm.file_server); 
      if (dataForm.reverse_proxy) form.setValue("reverse_proxy", dataForm.reverse_proxy); 
      if (dataForm.host) form.setValue("host", dataForm.host); 
    }
  }, [dataForm])

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col h-full overflow-hidden">
        <div className="grid flex-1 auto-rows-min gap-6 px-4 overflow-auto pb-3">
          <div className="flex flex-row items-center justify-between">
            <div>Example:</div>
            <ButtonGroup>
              <Button variant="secondary" size="sm" type="button" onClick={() => getDefaultForm("s_three_object")}>S3 (FE)</Button>
              <ButtonGroupSeparator />
              <Button variant="secondary" size="sm" type="button" onClick={() => getDefaultForm("frontend")}>FE</Button>
              <ButtonGroupSeparator />
              <Button variant="secondary" size="sm" type="button" onClick={() => getDefaultForm("backend")}>BE</Button>
              <ButtonGroupSeparator />
              {/* <Button variant="secondary" size="sm" type="button" onClick={() => getDefaultForm("s_three_object")}>FE & BE</Button> */}
            </ButtonGroup>
          </div>
          <FormField
            control={form.control}
            name="type"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Type Config{type}</FormLabel>
                <FormControl>
                    <Select name={field.name} value={field.value ?? ""} onValueChange={field.onChange}>
                      <SelectTrigger id="module_type" className={cn(`w-full ${field.value && 'border-green-700'}`)}>
                        <SelectValue placeholder="Module Type" />
                      </SelectTrigger>
                      <SelectContent position="item-aligned">
                        <SelectItem value="FE">Frontend</SelectItem>
                        <SelectItem value="BE">Backend</SelectItem>
                        {/* <SelectItem value="FE-BE">Frontend - Backend</SelectItem> */}
                      </SelectContent>
                    </Select>
                  </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="match_host"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Match Host</FormLabel>
                <FormControl>
                  <ComboboxMultipleCreatable value={field.value} onValueChange={field.onChange} className={cn(field.value?.length && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} placeholder="Match Host" />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh: example.com, subdomain.example.com</p>
              </FormItem>
            )}
          />
          {type !== "BE" && <FormField
            control={form.control}
            name="rewrite_uri"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Rewrite FE</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Rewrite FE" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh SPA: /index.html</p>
                <p className="text-xs">Contoh S3: /path/index.html</p>
              </FormItem>
            )}
          />}
          {type !== "BE" && <FormField
            control={form.control}
            name="file_server"
            render={({ field }) => (
              <FormItem>
                <FormLabel>File Server</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="File Server" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh SPA: /home/web/spa</p>
              </FormItem>
            )}
          />}
          <FormField
            control={form.control}
            name="match_path"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Match Path</FormLabel>
                <FormControl>
                  <ComboboxMultipleCreatable value={field.value} onValueChange={field.onChange} className={cn(field.value?.length && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} placeholder="Match Path" />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh Prefix: <strong>/api/*</strong> di proxy ke <strong>example.com/api/*</strong></p>
                <p className="text-xs">Contoh S3: <strong>/path/assets/*</strong> atau <strong>/path/favicon.ico</strong></p>
              </FormItem>
            )}
          />
          {type !== "FE" && <FormField
            control={form.control}
            name="rewrite_strip_path_prefix"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Rewrite Path Prefix</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Rewrite Path Prefix" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh: <strong>/api</strong>, (tergantung pada match path di atas) maka proxy <strong>/api</strong> menjadi path <strong>localhost:9000/</strong></p>
              </FormItem>
            )}
          />}
          {type !== "FE" && <FormField
            control={form.control}
            name="reverse_proxy"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Reverse Proxy</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Reverse Proxy" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh: localhost:9000 menjadi sumber proxy</p>
              </FormItem>
            )}
          />}
          {type !== "BE" && <FormField
            control={form.control}
            name="host"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Host</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Host" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                </FormControl>  
                <FormMessage />
                <p className="text-xs">Contoh: s3.example.com</p>
              </FormItem>
            )}
          />}
        </div>
        <div className="mt-auto flex flex-row w-full gap-2 p-4">
          <SheetClose asChild>
            <Button className="flex-1" variant="outline">
              <Ban />
              Close
            </Button>
          </SheetClose>
          <Button className="flex-1" type="submit">
            <Save />
            Save
          </Button>
        </div>
      </form>
    </Form>
  )
}

export default function CaddyManagement() {
  const [openDialog, setOpenDialog] = useState(false);
  const [openSheet, setOpenSheet] = useState(false);

  const { data, isLoading, mutate } = useSWR(`${BackendUrlBase}/check-local-service?url=http://localhost:2019/config&getRes=true`, fetchSWR);

  const currectDataCaddy = transformCaddyConfig(data?.apps?.http?.servers?.srv0);
  const [caddyConfig, setCaddyConfig] = useState<FlatCaddyConfig>();

  const columnsDetailCaddy: ColumnDef<FlatCaddyConfig>[] = [
    {
      accessorKey: "i",
      header: "Type",
      cell: ({ row }) => {
        const raw = row.original;
        let type = "FE";
        if (raw.rewrite_uri && raw.rewrite_strip_path_prefix) {
          type = "FE-BE";
        } else if (raw.rewrite_strip_path_prefix) {
          type = "BE";
        }

        return type;
      }
    },
    {
      accessorKey: "match_host",
      header: "Match Host",
    },
    {
      accessorKey: "rewrite_uri",
      header: "Rewrite (FE)",
    },
    {
      accessorKey: "file_server",
      header: "File Server (FE)",
    },
    {
      accessorKey: "match_path",
      header: "Match Path (BE)",
    },
    {
      accessorKey: "rewrite_strip_path_prefix",
      header: "Rewrite (BE)",
    },
    {
      accessorKey: "reverse_proxy",
      header: "Reverse Proxy (BE)",
    },
    {
      accessorKey: "host",
      header: "Host (S3)",
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const raw = row.original;
        raw.type = "FE";
        if (raw.rewrite_uri && raw.rewrite_strip_path_prefix) {
          raw.type = "FE-BE";
        } else if (raw.rewrite_strip_path_prefix) {
          raw.type = "BE";
        }

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="data-[state=open]:bg-muted text-muted-foreground flex size-6"
                size="icon"
              >
                <EllipsisVertical />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-32">
              <DropdownMenuItem onClick={() => {setOpenSheet(true);setCaddyConfig(raw);}}>Edit</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem variant="destructive" onClick={() => {setCaddyConfig(raw);setOpenDialog(true)}}>Delete</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      }
    }
  ];

  const [applyState, setApplyState] = useState(false);

  const applyConfig = async () => {
    setApplyState(true);

    const PORT_CADDY = ENV_STATUS.MODE === "development" ? [":2001"] : [":443", ":80"];
    // FOR PROXY TO GATEWAY
    // const PORT_CADDY = [":443", ":80"];

    let getAllCaddyConfig = currectDataCaddy;
    if (caddyConfig) {
      if (caddyConfig?.i !== undefined) {
        if (Object.keys(caddyConfig).length > 1) {
          getAllCaddyConfig[caddyConfig?.i] = caddyConfig;
        } else {
          getAllCaddyConfig = getAllCaddyConfig.slice(0, caddyConfig?.i);
        }
      } else {
        getAllCaddyConfig.push(caddyConfig);
      }
    }
    
    const parseConfig = reverseTransformCaddyConfig(getAllCaddyConfig, PORT_CADDY);
    const listTlsHosts = parseConfig?.routes?.flatMap(el =>
      el.match?.flatMap(elem => elem?.host ?? [] ) ?? []
    );
    const uniqueHosts = [...new Set(listTlsHosts)];
    
    const appendixForm = ENV_STATUS.MODE === "development" ? {
      pki: { // FOR LOCALHOST DEV
        certificate_authorities: {
          local: {
            install_trust: false
          }
        }
      }
    } : {
      tls: { // FOR TLS PRODUCTION
        certificates: {
          automate: uniqueHosts
        }
      }
    };
    // FOR PROXY TO GATEWAY
    // const appendixForm ={
    //   tls: {
    //     certificates: {
    //       automate: uniqueHosts
    //     }
    //   }
    // };

    const oldDataCaddy = data?.apps?.http?.servers?.srv0 ?? {}

    const initialForm = {
      http: {
        servers: {
          srv0: {
            ...oldDataCaddy,
            ...parseConfig
          }
        }
      }
    }
    const payload = {
      apps: {
        ...initialForm,
        ...appendixForm
      }
    };

    console.log(payload)

    try {
      const getCsrf = await FetchCsrfToken();
  
      const response = await fetch(`${BackendUrlBase}/proxy-local-service?url=http://localhost:2019/load`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error(response.statusText || "Apply Config Failed");

      const data = await response.json();
      console.log(data);
      
      mutate();
      setApplyState(false);
      setCaddyConfig(undefined);
      toast.success("Apply Config Success");
    } catch (err) {
      setApplyState(false);
      // setCaddyConfig(undefined);
      console.error("Apply Config error:", err);
      const errorMessage = err instanceof Error ? err.message : "Unknown error";
      toast.error("Apply Config failed: " + errorMessage);
    }
  }

  return (
    <div className="flex-1 flex-col gap-8 md:flex">
      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-col gap-1">
          <h2 className="text-2xl font-semibold tracking-tight">
            Caddy Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s list Caddy Setup
          </p>
        </div>
        <div className="flex items-center gap-2">
          {caddyConfig ? <Button onClick={applyState ? () => console.log("") : applyConfig} variant="destructive" disabled={applyState}>
            {applyState ? <Spinner /> : <Send /> }
            {applyState ? "Loading..." : "Apply" }
            <Badge variant="secondary">1</Badge>
          </Button> : null}
          <Button onClick={() => {setOpenSheet(true);setCaddyConfig(undefined)}} variant="outline">
            <Plus />
            Setup
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 gap-4">
        <div>
          {isLoading ? <p>Loading...</p> : <CardInformation 
            name="Caddy Setup" 
            description="There is list of Caddy for Domain" 
            rowIdKey="match_host"
            data={currectDataCaddy}
            columnsDetail={columnsDetailCaddy}
          />}
        </div>
      </div>
      <SheetForm 
        name="Caddy Form"
        openSheet={openSheet} 
        setOpenSheet={setOpenSheet}
      >
        <SheetFormChild 
          openSheet={openSheet}
          dataForm={caddyConfig}
          setOpenSheet={setOpenSheet}
          setCaddyConfig={setCaddyConfig}
        />
      </SheetForm>
      <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Setting Caddy?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will delete the data. Data deleted cannot be restored.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="overflow-auto">
            <pre>{JSON.stringify(caddyConfig, null, 2)}</pre>
          </div>
          <AlertDialogFooter>
            <AlertDialogAction>Cancel</AlertDialogAction>
            <AlertDialogCancel variant="destructive" onClick={() => {
              if (caddyConfig?.i !== undefined) {
                const { i, ...rest } = caddyConfig;
                console.log(rest);
                setCaddyConfig({ i });
              }
            }}>Delete</AlertDialogCancel>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}