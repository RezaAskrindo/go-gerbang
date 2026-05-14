import { useEffect, useState } from "react";
import { Ban, EllipsisVertical, Plus, Save } from "lucide-react";
import { toast } from "sonner";

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { cn } from "@/lib/utils";
import type { ColumnDef } from "@tanstack/react-table";


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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService";
import { useConfigurations, useDeleteConfigurations } from "@/services/use-swr-service";

import CardInformation from "@/components/card-information";

type TFormNotification = {
  notif_type: string
  sender: string
  email: string
  password?: string
  key?: string
  host?: string
  port?: string
  image_element?: string
}

type TNotifType = {
  EMAIL_SMTP_CONFIG: number
  EMAIL_RESEND_CONFIG: number
}

const formSchema = z.object({
  notif_type: z.string().min(2, {message: "Notification Type is required"}),
  sender: z.string().min(2, {message: "Sender is required"}),
  email: z.email(),
  password: z.string().optional(),
  key: z.string().optional(),
  host: z.string().optional(),
  port: z.string().optional(),
  image_element: z.string().optional(),
})

function SheetForm({
  openSheet,
  setOpenSheet,
  mutateData,
  indexData={
    EMAIL_SMTP_CONFIG: 0,
    EMAIL_RESEND_CONFIG: 0
  },
  dataForm,
}: {
  openSheet: boolean
  setOpenSheet: (v: boolean) => void
  mutateData: () => void
  indexData?: TNotifType
  dataForm?: TFormNotification
}) {

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      notif_type: "",
      sender: "",
      email: "",
    },
  });

  const [notify_type, image_element] = form.watch(["notif_type", "image_element"]);

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let payload = [
      {
        idConfiguration: undefined,
        configurationGroup: values.notif_type,
        configurationName: values.sender,
        configurationKey: "sender",
        configurationValue: values.sender,
        configurationIndex: indexData[values.notif_type as keyof TNotifType],
      },
      {
        idConfiguration: undefined,
        configurationGroup: values.notif_type,
        configurationName: values.sender,
        configurationKey: "email",
        configurationValue: values.email,
        configurationIndex: indexData[values.notif_type as keyof TNotifType],
      },
    ];
    if (values.notif_type === "EMAIL_SMTP_CONFIG" && values.password && values.host && values.port) {
      payload = payload.concat([
        {
          idConfiguration: undefined,
          configurationGroup: values.notif_type,
          configurationName: values.sender,
          configurationKey: "password",
          configurationValue: values.password,
          configurationIndex: indexData[values.notif_type as keyof TNotifType],
        },
        {
          idConfiguration: undefined,
          configurationGroup: values.notif_type,
          configurationName: values.sender,
          configurationKey: "host",
          configurationValue: values.host,
          configurationIndex: indexData[values.notif_type as keyof TNotifType],
        },
        {
          idConfiguration: undefined,
          configurationGroup: values.notif_type,
          configurationName: values.sender,
          configurationKey: "port",
          configurationValue: values.port,
          configurationIndex: indexData[values.notif_type as keyof TNotifType],
        }
      ])
    } else if (values.notif_type === "EMAIL_RESEND_CONFIG" && values.key) {
      payload = payload.concat([
        {
          idConfiguration: undefined,
          configurationGroup: values.notif_type,
          configurationName: values.sender,
          configurationKey: "key",
          configurationValue: values.key,
          configurationIndex: indexData[values.notif_type as keyof TNotifType],
        },
      ])
    }

    if (values.image_element) {
      payload.push({
        idConfiguration: undefined,
        configurationGroup: values.notif_type,
        configurationName: values.sender,
        configurationKey: "image_element",
        configurationValue: values.image_element,
        configurationIndex: indexData[values.notif_type as keyof TNotifType],
      })
    }
    
    if (dataForm && values.notif_type && values.sender) {
      await useDeleteConfigurations(values.notif_type, values.sender);
    }

    const getCsrf = await FetchCsrfToken();
        
    const response = await (await fetch(`${BackendUrlBase}/configurations`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
      body: JSON.stringify(payload),
    })).json();

    if (response.status) {
      form.reset();
      setOpenSheet(false);
      mutateData();
      toast.success("Success Save Notification");
    } else {
      toast.error("Failed Save Notification");
    }
  }

  useEffect(() => {
    if (dataForm) {
      form.setValue("notif_type", dataForm.notif_type);
      form.setValue("sender", dataForm.sender);
      form.setValue("email", dataForm.email);
      if (dataForm.image_element) form.setValue("image_element", dataForm.image_element);
      if (dataForm.password) form.setValue("password", dataForm.password);
      if (dataForm.host) form.setValue("host", dataForm.host);
      if (dataForm.port) form.setValue("port", dataForm.port);
      if (dataForm.key) form.setValue("key", dataForm.key);
    } else {
      form.reset();
    }
  }, [dataForm])

  return (
    <Sheet open={openSheet} onOpenChange={setOpenSheet}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>Edit Notification</SheetTitle>
          <SheetDescription>
            Make changes Notification here. Click save when you&apos;re done.
          </SheetDescription>
        </SheetHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col h-full overflow-hidden">
            <div className="grid flex-1 auto-rows-min gap-6 px-4 overflow-auto pb-3">
              <FormField
                control={form.control}
                name="notif_type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Module Type</FormLabel>
                    <FormControl>
                      <Select name={field.name} value={field.value} onValueChange={field.onChange}>
                        <SelectTrigger id="module_type" className={cn(`w-full ${field.value && 'border-green-700'}`)}>
                          <SelectValue placeholder="Module Type" />
                        </SelectTrigger>
                        <SelectContent position="item-aligned">
                          <SelectItem value="EMAIL_SMTP_CONFIG">SMTP</SelectItem>
                          <SelectItem value="EMAIL_RESEND_CONFIG">Resend</SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="sender"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Sender Name</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Sender Name" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Email" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              { notify_type === "EMAIL_SMTP_CONFIG" && <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Password" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              /> }
              { notify_type === "EMAIL_SMTP_CONFIG" && <FormField
                control={form.control}
                name="host"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Host</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Host" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              /> }
              { notify_type === "EMAIL_SMTP_CONFIG" && <FormField
                control={form.control}
                name="port"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Port</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Port" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              /> }
              { notify_type === "EMAIL_RESEND_CONFIG" && <FormField
                control={form.control}
                name="key"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Key</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Key" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              /> }
              <FormField
                control={form.control}
                name="image_element"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Image Element</FormLabel>
                    <FormControl>
                      <Input type="text" placeholder="Image Element" className={cn(field.value && 'border-green-700 focus:border-green-700! focus:ring-green-700/40!')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              
              {image_element && <div className="p-2 border rounded-lg"><div dangerouslySetInnerHTML={{ __html: image_element }} /></div>}
              
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
      </SheetContent>
    </Sheet>
  )
}

export default function NotificationManagement() {
  const [openSheet, setOpenSheet] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [dataForm, setDataFrom] = useState<TFormNotification>();

  const { data: dataResend, mutate: mutateResend } = useConfigurations("EMAIL_RESEND_CONFIG");
  const { data: dataSMTP, mutate: mutateSMTP } = useConfigurations("EMAIL_SMTP_CONFIG");

  const mutateData = () => {
    mutateResend();
    mutateSMTP();
  }

  useEffect(() => {
    if (!openSheet) {
      console.log("Here")
      mutateData();
      setDataFrom(undefined);
    }
  }, [openSheet]);

  const columnsDetailResend: ColumnDef<TFormNotification>[] = [
    {
      accessorKey: "sender",
      header: "Sender",
    },
    {
      accessorKey: "email",
      header: "User Email",
    },
    {
      accessorKey: "key",
      header: "Key Resend",
    },
    {
      accessorKey: "image_element",
      header: "Logo Email",
      cell: ({ row }) => {
        const raw = row.original;
        if (raw.image_element) return <div dangerouslySetInnerHTML={{__html: raw.image_element}}></div> 
        else return null;
      }
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const raw = row.original;
        raw.notif_type = "EMAIL_RESEND_CONFIG";

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
              <DropdownMenuItem onClick={() => {setOpenSheet(true);setDataFrom(raw)}}>
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {
                toast.promise(
                  fetch(`${BackendUrlBase}/check-mail?to=rezaoda@gmail.com&appName=${raw.sender}&provider=${raw.sender}`, {
                    method: "GET",
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
              }}>Test</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem variant="destructive" onClick={() => {setOpenDialog(true);setDataFrom(raw)}}>
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      }
    }
  ];

  const columnsDetailSMTP: ColumnDef<TFormNotification>[] = [
    {
      accessorKey: "sender",
      header: "Sender",
    },
    {
      accessorKey: "email",
      header: "User Email",
    },
    {
      accessorKey: "password",
      header: "Password Email",
    },
    {
      accessorKey: "host",
      header: "Host Email",
    },
    {
      accessorKey: "port",
      header: "Port Email",
    },
    {
      accessorKey: "image_element",
      header: "Logo Email",
      cell: ({ row }) => {
        const raw = row.original;
        if (raw.image_element) return <div dangerouslySetInnerHTML={{__html: raw.image_element}}></div> 
        else return null;
      }
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const raw = row.original;
        raw.notif_type = "EMAIL_SMTP_CONFIG";

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
              <DropdownMenuItem onClick={() => {setOpenSheet(true);setDataFrom(raw)}}>
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {
                toast.promise(
                  fetch(`${BackendUrlBase}/check-mail?to=rezaoda@gmail.com&appName=${raw.sender}&provider=${raw.sender}`, {
                    method: "GET",
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
              }}>Test</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem variant="destructive" onClick={() => {setOpenDialog(true);setDataFrom(raw)}}>
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
            Notification Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s list Notification Setup
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => setOpenSheet(true)} variant="outline">
            <Plus />
            NOTIFICATION
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 gap-4">
        <div>
          <CardInformation name="Resend Setup" description="There is list of Resend Email for Notification" rowIdKey="email" columnsDetail={columnsDetailResend} data={dataResend?.data ?? []} />
        </div>
        <div>
          <CardInformation name="SMTP Setup" description="There is list of SMTP Email for Notification" rowIdKey="email" columnsDetail={columnsDetailSMTP} data={dataSMTP?.data ?? []} />
        </div>
      </div>
      <SheetForm 
        openSheet={openSheet} 
        setOpenSheet={setOpenSheet} 
        dataForm={dataForm} 
        mutateData={mutateData}
        indexData={{
          EMAIL_SMTP_CONFIG: (dataSMTP?.data?.length ?? 1) - 1,
          EMAIL_RESEND_CONFIG: (dataResend?.data?.length ?? 1) - 1
        }} 
      />
      <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {dataForm?.sender} Notification?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will delete the data. Data deleted cannot be restored.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction>Cancel</AlertDialogAction>
            <AlertDialogCancel variant="destructive" onClick={() => {
              if (dataForm?.notif_type) {
                toast.promise(
                  useDeleteConfigurations(dataForm?.notif_type, dataForm?.sender),
                  {
                    loading: "Waiting...",
                    success: () => {
                      mutateData();
                      setOpenDialog(false);
                      return "Success"
                    },
                    error: (err) => {
                      setOpenDialog(false);
                      return err.message || "Failed"
                    },
                  } 
                )
              }
            }}>Delete</AlertDialogCancel>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}