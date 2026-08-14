import { lazy, useEffect, useState } from "react";
import { Ban, Check, CheckCircle, EllipsisVertical, Files, Loader, Plus, Save, X } from "lucide-react";
import { useDropzone, type FileWithPath } from 'react-dropzone';  
import { toast } from "sonner";

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { cn } from "@/lib/utils";
import type { ColumnDef } from "@tanstack/react-table";
import useSWR from "swr";

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
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { SheetClose } from "@/components/ui/sheet";

import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService";
import { fetchSWR, useConfiguration, useDeleteConfiguration } from "@/services/use-swr-service";

import CardInformation from "@/components/card-information";
import SheetForm from "@/components/sheet-form";
import { useHash } from "@/hooks/use-hash";
import { Switch } from "@/components/ui/switch";

const CodeHighlighter = lazy(() => import("@/components/CodeHighlighter"));

type TDetailModule = {
  module_name: string
  module_type: string
  location: string
  execution?: string
  desist?: string
  url?: string
  use_s3?: boolean
  s3_config?: string
}

const formSchema = z.object({
  module_type: z.string().min(2, {message: "Module Type is required"}),
  module_name: z.string().min(2, {message: "Module Name is required"}),
  location: z.string().min(2, {message: "Location is required"}),
  execution: z.string().optional(),
  desist: z.string().optional(),
  url: z.string().optional(),
  use_s3: z.boolean().optional(),
})

function FormModule({
  openDialog,
  setOpenDialog,
  indexData=0,
  data,
  moduleType
}: {
  openDialog: boolean
  setOpenDialog: (v: boolean) => void
  indexData?: number
  data?: TDetailModule
  moduleType?: string
}) {
  console.log(indexData)
  const [loading, setLoading] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      module_type: moduleType ?? "",
      module_name: "",
      location: "",
      execution: "",
      use_s3: false,
    },
  });

  const [module_type] = form.watch(["module_type"]);

  // const messageErrors = useMemo(
  //   () => extractMessages(form.formState.errors),
  //   [form.formState.errors]
  // )

  const [droppedFiles, setDroppedFiles] = useState<FileWithPath[]>([]);
  const [isDragging, setIsDragging] = useState(false); 

  useEffect(() => {
    if (data && openDialog) {
      form.setValue("module_type", data.module_type);
      form.setValue("module_name", data.module_name);
      form.setValue("location", data.location);
      if (data.execution) form.setValue("execution", data.execution);
      if (data.desist) form.setValue("desist", data.desist);
      if (data.url) form.setValue("url", data.url);
      if (data.s3_config) form.setValue("use_s3", true);
    } else {
      form.reset();
      setDroppedFiles([]);
    }
  }, [data, openDialog])

  const onDrop = (droppedFiles: File[]) => {  
    // Handle dropped files (we’ll expand this later)  
    // console.log('Folder contents:', droppedFiles);
    setDroppedFiles(droppedFiles);
  };

  const { getRootProps, getInputProps } = useDropzone({  
    onDrop,
    onDragEnter: () => setIsDragging(true),
    onDragLeave: () => setIsDragging(false),
    onDropAccepted: () => setIsDragging(false),
    onDropRejected: () => setIsDragging(false),
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (droppedFiles.length === 0 && !data) {
      alert("Please select a file or folder");
      return;
    }

    if (data && values.module_name) {
      await useDeleteConfiguration("MODULE_CONFIG", values.module_name);
    }

    if (loading) {
      return;
    }

    setLoading(true);
    
    let payload = [
      {
        idConfiguration: undefined,
        configurationGroup: "MODULE_CONFIG",
        configurationName: values.module_name,
        configurationKey: "module_name",
        configurationValue: values.module_name,
        configurationIndex: indexData,
      },
      {
        idConfiguration: undefined,
        configurationGroup: "MODULE_CONFIG",
        configurationName: values.module_name,
        configurationKey: "module_type",
        configurationValue: values.module_type,
        configurationIndex: indexData,
      },
      {
        idConfiguration: undefined,
        configurationGroup: "MODULE_CONFIG",
        configurationName: values.module_name,
        configurationKey: "location",
        configurationValue: values.location,
        configurationIndex: indexData,
      },
    ];

    let encodedConfig: string | null = null;
    if (values.use_s3) {
      const config = {
        access_key: "K5C9MOVDVOJMJ1ZMCIKH",
        secret_key: "S02M0rtfIyFnJESZrTOoaGw4gmwRxDOIboQg4N6Z",
        region: "us-east-1",
        endpoint: "s3.nevaobjects.id",
        bucket: "mfe",
      };

      const bytes = new TextEncoder().encode(JSON.stringify(config));

      encodedConfig = btoa(String.fromCharCode(...bytes));

      payload = payload.concat([
        {
          idConfiguration: undefined,
          configurationGroup: "MODULE_CONFIG",
          configurationName: values.module_name,
          configurationKey: "s3_config",
          configurationValue: encodedConfig,
          configurationIndex: indexData,
        }
      ])
    }

    if (values.execution && values.desist && values.url) {
      payload = payload.concat([
        {
          idConfiguration: undefined,
          configurationGroup: "MODULE_CONFIG",
          configurationName: values.module_name,
          configurationKey: "execution",
          configurationValue: values.execution,
          configurationIndex: indexData,
        },
        {
          idConfiguration: undefined,
          configurationGroup: "MODULE_CONFIG",
          configurationName: values.module_name,
          configurationKey: "desist",
          configurationValue: values.desist,
          configurationIndex: indexData,
        },
        {
          idConfiguration: undefined,
          configurationGroup: "MODULE_CONFIG",
          configurationName: values.module_name,
          configurationKey: "url",
          configurationValue: values.url,
          configurationIndex: indexData,
        },
      ]);
    }

    const getCsrf = await FetchCsrfToken();
    
    const res = await fetch(`${BackendUrlBase}/Configuration`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
      body: JSON.stringify(payload),
    })

    const response = await res.json();

    if (response.status && droppedFiles.length) {
      const uploadData = new FormData();
      if (values.use_s3 && encodedConfig) {
        uploadData.append("s3_config", encodedConfig);
      }
      uploadData.append("file-location", values.location); // don't change file-location
      droppedFiles.forEach(file => {
        uploadData.append("files", file, file.relativePath || file.name);
      });

      try {
        const response = await fetch(`${BackendUrlBase}/upload-file`, {
          method: "POST",
          body: uploadData,
        });

        const result = await response.json();
        if (!response.ok) throw new Error(result.message || "Upload failed");
        setDroppedFiles([]);
        // toast.success(result.message || "Upload successful");
        toast.success("Succes Save Module");
        // TODO: reset form or show success
      } catch (err) {
        console.error("Upload error:", err);
        const errorMessage = err instanceof Error ? err.message : "Unknown error";
        alert("Upload failed: " + errorMessage);
      }
    }
    
    setLoading(false);
    setOpenDialog(false);
  }

  return (
    <SheetForm name="Edit Module" openSheet={openDialog} setOpenSheet={setOpenDialog} >
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col h-full overflow-hidden">
          <div className="grid flex-1 auto-rows-min gap-6 px-4 overflow-auto pb-3">

            <FormField
              control={form.control}
              name="module_type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Module Type</FormLabel>
                  <FormControl>
                    <Select name={field.name} value={field.value} onValueChange={field.onChange}>
                      <SelectTrigger id="module_type" className={cn(`w-full ${field.value && 'border-green-700'}`)}>
                        <SelectValue placeholder="Module Type" />
                      </SelectTrigger>
                      <SelectContent position="item-aligned">
                        <SelectItem value="Frontend">Frontend</SelectItem>
                        <SelectItem value="Backend">Backend</SelectItem>
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {module_type === "Frontend" && <FormField
              control={form.control}
              name="use_s3"
              render={({ field }) => (
                <FormItem className="flex flex-row justify-between">
                  <FormLabel>S3 Object?</FormLabel>
                  <FormControl>
                    <Switch name={field.name} checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />}
            <FormField
              control={form.control}
              name="module_name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Module Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Module Name" className={cn(field.value && 'border-green-700 focus:ring-green-700/40!')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="location"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Location</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="Location" className={cn(field.value && 'border-green-700 focus:ring-green-700/40!')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {module_type === "Backend" && <FormField
              control={form.control}
              name="execution"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>File Exec Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="File Exec Name" className={cn(field.value && 'border-green-700 focus:ring-green-700/40!')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />}
            {module_type === "Backend" && <FormField
              control={form.control}
              name="desist"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>File Desist Name</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="File Desist Name" className={cn(field.value && 'border-green-700 focus:ring-green-700/40!')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />}
            {module_type === "Backend" && <FormField
              control={form.control}
              name="url"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>URL for status</FormLabel>
                  <FormControl>
                    <Input type="text" placeholder="URL for status" className={cn(field.value && 'border-green-700 focus:ring-green-700/40!')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />}
            
            <div>
              <Label>Module Files</Label>
              <div className={`my-2 flex justify-center rounded-lg border-dashed px-6 py-10 ${isDragging ? 'border-3' : 'border-2'}`} {...getRootProps()}>
                <div className="text-center">
                  <Files className="mx-auto size-10" />
                  <div className="mt-4 flex justify-center text-sm/6">
                    <Label htmlFor="file-upload">
                      <span>{droppedFiles.length > 0 ? `${droppedFiles.length} file(s) selected` : "Click here to Upload"}</span>
                      <input id="file-upload" name="files" {...getInputProps()} />
                    </Label>
                  </div>
                  {droppedFiles.length === 0 && <p className="text-sm mt-2 italic">or drag and drop here</p>}
                </div>
              </div>
            </div>
          </div>
          
          <div className="mt-auto flex flex-row w-full gap-2 p-4">
            <SheetClose asChild>
              <Button className="flex-1" variant="outline">
                <Ban />
                Close
              </Button>
            </SheetClose>
            <Button className="flex-1" type={loading ? "button" : "submit"} disabled={loading}>
              <Save />
              {loading ? "Loading..." : "Save" }
            </Button>
          </div>

        </form>
      </Form>
    </SheetForm>
  )
}

function ContentBackendGuide() {
  return (
    <div>
      <h2 className="font-bold mb-2">Petunjuk Untuk Backend Module:</h2>
      <ol className="text-sm list-decimal flex flex-col gap-2 ms-5">
        <li>
          <span className="font-semibold underline">Module Name</span> berisikan nama Module yang akan di buat.
        </li>
        <li>
          <span className="font-semibold underline">Location</span> berisikan full-path folder dimana Module akan di upload.
        </li>
        <li>
          <span className="font-semibold underline">File Exec Name</span> berisikan Module akan di jalankan / execute (windows: .exe, linux: no ext).
        </li>
        <li>
          <span className="font-semibold underline">File Desist Name</span> berisikan file running script (.sh) untuk menjalankan File Exec di atas. Seperti running_script.sh, berikut contoh script yang digunakan.
          <pre className="no-scrollbar min-w-0 overflow-x-auto outline-none has-data-highlighted-line:px-0 has-data-line-numbers:px-0 has-data-[slot=tabs]:p-0">
            <CodeHighlighter lang="sh" code={`
#!/bin/bash

pkill -f go-service
sleep 2

BIN="/full/path/go-service"

if [ ! -x "$BIN" ]; then
  echo "Error: binary '$BIN' not found or not executable."
  exit 1
fi

export PORT=8080
export DATABASE_URL="postgres://dbUser:dbPass@localhost:5432/dbName?sslmode=disable"
export ORIGINS="http://localhost:3000"

"$BIN" & disown

echo "go-service restarted successfully."
              `} />
          </pre>
        </li>
        <li>
          <span className="font-semibold underline">URL for status</span> berisikan url path dengan port untuk melakukan test kepada service tersebut. Seperti http://localhost:8080.
        </li>
        <li>
          <span className="font-semibold underline">Module File</span> berisikan file-file yang di upload ke folder, minimal file exec dan file desist.
        </li>
      </ol>
    </div>
  )
}


function ContentFrontendGuide() {
  return (
    <div>
      <h2 className="font-bold mb-2">Petunjuk Untuk Frontend Module:</h2>
      <ol className="text-sm list-decimal flex flex-col gap-2 ms-5">
        <li>
          <span className="font-semibold underline">S3 Object</span> jika iya, maka akan di upload ke S3 Server. Jika tidak maka akan di upload ke local.
        </li>
        <li>
          <span className="font-semibold underline">Module Name</span> berisikan nama Module yang akan di buat.
        </li>
        <li>
          <span className="font-semibold underline">Location</span> tergantung apakah menggunakan S3 object atau tidak? jika iya, maka cukup folder nya saja, apabila tidak maka harus full path.
        </li>
      </ol>
    </div>
  )
}

export default function ModuleManagement() {
  const [openDialog, setOpenDialog] = useState(false);
  const [openAlert, setOpenAlert] = useState(false);
  const [dataForm, setDataFrom] = useState<TDetailModule>();

  const hash = useHash();
  const hashMenu = hash.replace(/#\//g, '').replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase()).split(" ");
  const getModule = hashMenu?.[1] ?? "backend";

  const { data: moduleConfig, mutate } = useConfiguration("MODULE_CONFIG");

  const filterModule = moduleConfig?.data?.filter((el: TDetailModule) => el.module_type === getModule);

  useEffect(() => {
    if (!openDialog) {
      mutate();
      setDataFrom(undefined);
    }
  }, [openDialog, mutate]);

  
  const RunScript = (work_dir: string, file: string) => {
    toast.promise(
      fetch(`${BackendUrlBase}/Configuration/execute?work_dir=${work_dir}&file=${file}`).then(async (res) => {
        if (!res.ok) throw new Error("Request failed")
        const data = await res.json()
        if (!data.status) throw new Error(data.message || "Failed to Execute")
        return data
      }),
      {
        loading: "Waiting...",
        success: () => {
          return "Success Execute"
        },
        error: (err) => {
          return err.message || "Failed"
        },
      }
    )
  }

  const columnsDetailBackend: ColumnDef<TDetailModule>[] = [
    {
      accessorKey: "module_name",
      header: "Module Name",
    },
    {
      accessorKey: "location",
      header: "Location",
    },
    {
      accessorKey: "url",
      header: "Status",
      cell: ({ row }) => {
        const module = row.original;
        const { data, isLoading } = useSWR(`${BackendUrlBase}/check-local-service?url=${module.url}`, fetchSWR, {
          revalidateOnFocus: true
        });
        if (isLoading) return <div><Loader /></div>
        else if (data.data) return <div><Check className="text-green-500 size-5" /></div>
        else return <div><X className="text-red-600" /></div>
      }
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const module = row.original;

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="data-[state=open]:bg-muted text-muted-foreground flex size-8"
                size="icon"
              >
                <EllipsisVertical />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-32">
              <DropdownMenuItem onClick={() => {setOpenDialog(true);setDataFrom(module)}}>Edit</DropdownMenuItem>
              {/* {module.execution && <DropdownMenuItem onClick={() => RunScript(module.location, module.execution as string)}>Start</DropdownMenuItem>} */}
              {module.desist && <DropdownMenuItem onClick={() => RunScript(module.location, module.desist as string)}>Restart</DropdownMenuItem>}
              <DropdownMenuSeparator />          
              <DropdownMenuItem onClick={() => {setOpenAlert(true);setDataFrom(module)}} variant="destructive">
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      }
    }
  ];

  const columnsDetailFrontend: ColumnDef<TDetailModule>[] = [
    {
      accessorKey: "module_name",
      header: "Module Name",
    },
    {
      accessorKey: "location",
      header: "Location",
    },
    {
      accessorKey: "s3_config",
      header: "S3 Object",
      cell: ({ row }) => {
        return row.original.s3_config ? <CheckCircle className="text-green-500 size-5" /> : null;
      }
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const module = row.original;

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="data-[state=open]:bg-muted text-muted-foreground flex size-8"
                size="icon"
              >
                <EllipsisVertical />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-32">
              <DropdownMenuItem onClick={() => {setOpenDialog(true);setDataFrom(module)}}>Edit</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => {setOpenAlert(true);setDataFrom(module)}} variant="destructive">
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      }
    }
  ];


  return (
    <div className="flex-1 flex-col gap-8 md:flex">
      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-col gap-1">
          <h2 className="text-2xl font-semibold tracking-tight">
            Module {getModule} Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s Module Management list. Only work for two type module: Backend (execute file) and Front end (static file)
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => {setOpenDialog(true);setDataFrom(undefined)}} variant="outline">
            <Plus />
            MODULE
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2">
          <CardInformation name={`${getModule.toUpperCase()} MODULE`} rowIdKey="module_name" columnsDetail={getModule === "Frontend" ? columnsDetailFrontend : columnsDetailBackend} data={filterModule} />
        </div>
        <div>
          <CardInformation name={`${getModule.toUpperCase()} MODULE GUIDE`} 
            Content={
              getModule === "Frontend" ?
              ContentFrontendGuide :
              ContentBackendGuide
            }
          />
        </div>
      </div>
      <FormModule 
        key={getModule}
        openDialog={openDialog} 
        setOpenDialog={setOpenDialog} 
        data={dataForm} 
        // indexData={(moduleConfig?.data?.length ?? 1) - 1} 
        indexData={moduleConfig?.data?.length} 
        moduleType={getModule}
      />
      <AlertDialog open={openAlert} onOpenChange={setOpenAlert}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {dataForm?.module_name} Module?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will delete the data. Data deleted cannot be restored.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction variant="outline">Cancel</AlertDialogAction>
            {dataForm?.module_name && <AlertDialogCancel variant="destructive" onClick={() => useDeleteConfiguration("MODULE_CONFIG", dataForm?.module_name)}>Delete</AlertDialogCancel>}
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}