import z from "zod";
import useSWR from "swr";

import { useEffect, useState } from "react";
import { Ban, EllipsisVertical, Plus, Save } from "lucide-react";
import { toast } from "sonner";
import { cn } from "cn";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
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
  Dialog,
  DialogContent,
  DialogClose,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
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
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";

import CardInformation from "@/components/card-information";
import TableSkeleton from "@/components/go-gerbang/table-skeleton";

import { fetchSWR } from "@/services/use-swr-service";
import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService";

const baseUrlMailNotif = `${BackendUrlBase}/api/notification/mail/accounts`

type TFormNotificationMail = {
  id: string
  tenant_id: string
  name: string
  from_name: string
  from_email: string
  smtp_host: string
  smtp_port: number
  smtp_user: string
  use_ssl?: boolean
  body_html?: string
  body_text?: string
  status: string
  created_at: string
  updated_at: string
}

const formSchema = z.object({
  name: z.string().min(2, {message: "Name is required"}),
  from_name: z.string().min(2, {message: "From Name is required"}),
  from_email: z.string().email({message: "Valid email required"}),
  smtp_host: z.string().min(2, {message: "SMTP Host is required"}),
  smtp_port: z.string().regex(/^\d+$/, {message: "Port must be a number"}),
  smtp_user: z.string().min(2, {message: "SMTP User is required"}),
  use_ssl: z.boolean().default(true).optional(),
  status: z.enum(["active", "inactive"]),
});

type TFormNotificationMailSchema = z.infer<typeof formSchema>;

export function SheetNotificationMailForm({
  openSheet,
  setOpenSheet,
  dataForm,
  setDataForm,
  mutateData,
}: {
  openSheet: boolean;
  setOpenSheet: (open: boolean) => void;
  dataForm?: TFormNotificationMail;
  setDataForm?: (data: TFormNotificationMail | undefined) => void;
  mutateData: () => void
}) {
  const form = useForm<TFormNotificationMailSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      from_name: "",
      from_email: "",
      smtp_host: "",
      smtp_port: "",
      smtp_user: "",
      use_ssl: true,
      status: "active",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (dataForm?.id) {
      toast.success("Update Not Working Yet!");
      return;
    }
    // console.log(dataForm)
    // console.log(values);
    const urlSubmit = baseUrlMailNotif;

    const getCsrf = await FetchCsrfToken();
            
    const response = await (await fetch(urlSubmit, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
      body: JSON.stringify(values),
    })).json();

    if (response.status) {
      form.reset();
      setOpenSheet(false);
      mutateData();
      setDataForm?.(undefined);
      toast.success("Success Save Notification");
    } else {
      toast.error("Failed Save Notification");
    }
  }

  useEffect(() => {
    if (dataForm) {
      form.setValue("name", dataForm.name);
      form.setValue("from_name", dataForm.from_name);
      form.setValue("from_email", dataForm.from_email);
      form.setValue("smtp_host", dataForm.smtp_host);
      form.setValue("smtp_port", String(dataForm.smtp_port));
      form.setValue("smtp_user", dataForm.smtp_user);
      form.setValue("status", dataForm.status as "active" | "inactive");
      if (dataForm.use_ssl) form.setValue("use_ssl", dataForm.use_ssl);
    } else {
      form.reset();
    }
  }, [dataForm, form]);

  return (
    <Sheet open={openSheet} onOpenChange={setOpenSheet}>
      <SheetContent>
        <SheetHeader className="pb-0">
          <SheetTitle>
            {dataForm ? "Edit Notification Mail" : "Add Notification Mail"}
          </SheetTitle>
          <SheetDescription>
            Configure SMTP settings and email templates. Click save when you&apos;re done.
          </SheetDescription>
        </SheetHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col h-full overflow-hidden">
            <div className="grid flex-1 auto-rows-min gap-6 px-4 overflow-auto pb-3">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Configuration Name</FormLabel>
                    <FormControl>
                      <Input 
                        type="text" 
                        placeholder="e.g., Production SMTP" 
                        className={cn(field.value && 'border-green-700')}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="from_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>From Name</FormLabel>
                    <FormControl>
                      <Input 
                        type="text" 
                        placeholder="e.g., Support Team" 
                        className={cn(field.value && 'border-green-700')}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="from_email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>From Email</FormLabel>
                    <FormControl>
                      <Input 
                        type="email" 
                        placeholder="noreply@example.com" 
                        className={cn(field.value && 'border-green-700')}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="smtp_host"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SMTP Host</FormLabel>
                    <FormControl>
                      <Input 
                        type="text" 
                        placeholder="smtp.gmail.com" 
                        className={cn(field.value && 'border-green-700')}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="smtp_port"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SMTP Port</FormLabel>
                    <FormControl>
                      <Input 
                        type="text" 
                        placeholder="587" 
                        className={cn(field.value && 'border-green-700')}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="smtp_user"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SMTP User</FormLabel>
                    <FormControl>
                      <Input 
                        type="text" 
                        placeholder="your-email@gmail.com" 
                        className={cn(field.value && 'border-green-700')}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="use_ssl"
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
                    <div className="space-y-0.5">
                      <FormLabel>Use SSL/TLS</FormLabel>
                      <p className="text-sm text-muted-foreground">
                        Enable secure SMTP connection
                      </p>
                    </div>
                    <FormControl>
                      <Switch 
                        checked={field.value} 
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="status"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Status</FormLabel>
                    <Select name={field.name} value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger className={cn(field.value && 'border-green-700')}>
                          <SelectValue placeholder="Select status" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="active">Active</SelectItem>
                        <SelectItem value="inactive">Inactive</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

            </div>

            <div className="mt-auto flex flex-row w-full gap-2 p-4 border-t">
              <SheetClose asChild>
                <Button className="flex-1" variant="outline">
                  <Ban className="w-4 h-4" />
                  Close
                </Button>
              </SheetClose>
              <Button className="flex-1" type="submit">
                <Save className="w-4 h-4" />
                Save
              </Button>
            </div>
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  );
}

const templateFormSchema = z.object({
  body_html: z.string().optional(),
  body_text: z.string().optional(),
});

type TFormTemplateSchema = z.infer<typeof templateFormSchema>;

export function NotificationMailTemplateDialog({
  openTemplateDialog,
  setOpenTemplateDialog,
  dataForm,
  setDataForm,
  mutateData,
}: {
  openTemplateDialog: boolean;
  setOpenTemplateDialog: (open: boolean) => void;
  dataForm?: TFormNotificationMail;
  setDataForm?: (data: TFormNotificationMail | undefined) => void;
  mutateData: () => void
}) {
  const form = useForm<TFormTemplateSchema>({
    resolver: zodResolver(templateFormSchema),
    defaultValues: {
      body_html: "",
      body_text: "",
    },
  });

  const bodyHtml = form.watch("body_html");
  const bodyText = form.watch("body_text");
  const [previewTab, setPreviewTab] = useState<"html" | "text">("html");

  async function onSubmit(values: z.infer<typeof templateFormSchema>) {
    if (dataForm?.id) {
      const urlSubmit = `${baseUrlMailNotif}/${dataForm.id}/layout`;
  
      const getCsrf = await FetchCsrfToken();
              
      const response = await (await fetch(urlSubmit, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
        body: JSON.stringify(values),
      })).json();
  
      if (response.status) {
        form.reset();
        setOpenTemplateDialog(false);
        mutateData();
        setDataForm?.(undefined);
        toast.success("Success Save Layout");
      } else {
        toast.error("Failed Save Layout");
      }
    } else {
      toast.error("Need ID Params");
    }
  }

  useEffect(() => {
    if (dataForm) {
      if (dataForm.body_html) form.setValue("body_html", dataForm.body_html);
      if (dataForm.body_text) form.setValue("body_text", dataForm.body_text);
    } else {
      form.reset();
    }
  }, [dataForm, form]);

  return (
    <Dialog open={openTemplateDialog} onOpenChange={setOpenTemplateDialog}>
      <DialogContent className="min-w-3/4 max-h-[90vh] overflow-hidden flex flex-col">
        <DialogHeader>
          <DialogTitle>Edit Email Templates</DialogTitle>
          <DialogDescription>
            Configure HTML and plain text email templates with live preview.
          </DialogDescription>
        </DialogHeader>
          <Form {...form}>
            <form 
              onSubmit={form.handleSubmit(onSubmit)}
              className="flex flex-col h-full overflow-hidden flex-1"
            >
              <div className="grid grid-cols-2 gap-4 flex-1 overflow-hidden px-6 py-4">
                
                {/* Left side - Editors */}
                <div className="flex flex-col gap-4 overflow-auto">
                  <FormField
                    control={form.control}
                    name="body_html"
                    render={({ field }) => (
                      <FormItem className="flex-1 flex flex-col">
                        <FormLabel>HTML Body</FormLabel>
                        <FormControl className="flex-1">
                          <Textarea 
                            placeholder="<html><body>Your email template here...</body></html>" 
                            className={cn(field.value && 'border-green-700', "flex-1 resize-none font-mono text-xs")}
                            {...field} 
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="body_text"
                    render={({ field }) => (
                      <FormItem className="flex-1 flex flex-col">
                        <FormLabel>Plain Text Body</FormLabel>
                        <FormControl className="flex-1">
                          <Textarea 
                            placeholder="Plain text version of your email template..." 
                            className={cn(field.value && 'border-green-700', "flex-1 resize-none font-mono text-xs")}
                            {...field} 
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                {/* Right side - Preview */}
                <div className="flex flex-col gap-4 overflow-hidden border rounded-lg p-4 bg-muted/50">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-semibold">Preview</span>
                    <div className="ml-auto flex gap-2">
                      <button
                        type="button"
                        onClick={() => setPreviewTab("html")}
                        className={cn(
                          "px-3 py-1 text-xs rounded-md transition-colors",
                          previewTab === "html"
                            ? "bg-primary text-primary-foreground"
                            : "bg-muted text-muted-foreground hover:bg-muted-foreground/20"
                        )}
                      >
                        HTML
                      </button>
                      <button
                        type="button"
                        onClick={() => setPreviewTab("text")}
                        className={cn(
                          "px-3 py-1 text-xs rounded-md transition-colors",
                          previewTab === "text"
                            ? "bg-primary text-primary-foreground"
                            : "bg-muted text-muted-foreground hover:bg-muted-foreground/20"
                        )}
                      >
                        Text
                      </button>
                    </div>
                  </div>

                  <div className="flex-1 border rounded p-4 bg-white overflow-auto">
                    {previewTab === "html" ? (
                      bodyHtml ? (
                        <iframe
                          srcDoc={bodyHtml}
                          title="HTML Preview"
                          className="w-full h-full border-none"
                          sandbox="allow-same-origin"
                        />
                      ) : (
                        <div className="flex items-center justify-center h-full">
                          <p className="text-xs text-muted-foreground">HTML preview will appear here</p>
                        </div>
                      )
                    ) : bodyText ? (
                      <div className="p-4">
                        <p className="text-sm whitespace-pre-wrap">{bodyText}</p>
                      </div>
                    ) : (
                      <div className="flex items-center justify-center h-full">
                        <p className="text-xs text-muted-foreground">Text preview will appear here</p>
                      </div>
                    )}
                  </div>
                </div>

              </div>

              <div className="mt-auto flex flex-row gap-2 p-4 border-t">
                <DialogClose asChild>
                  <Button className="flex-1" variant="outline">
                    <Ban className="w-4 h-4" />
                    Close
                  </Button>
                </DialogClose>
                <Button className="flex-1" type="submit">
                  <Save className="w-4 h-4" />
                  Save Templates
                </Button>
              </div>
            </form>
          </Form>
        
      </DialogContent>
    </Dialog>
  );
}

export default function NotificationManagement() {
  const tenant_id = "7ebf9c33-fd7f-41eb-93e9-f05e546d9ba6";

  const [openSheet, setOpenSheet] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [openTemplateDialog, setOpenTemplateDialog] = useState(false);
  const [dataForm, setDataForm] = useState<TFormNotificationMail | undefined>(() => ({
    id: "",
    tenant_id: tenant_id,
    name: "",
    from_name: "",
    from_email: "",
    smtp_host: "",
    smtp_port: 587,
    smtp_user: "",
    use_ssl: true,
    body_html: undefined,
    body_text: undefined,
    status: "inactive",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }));

  const { data, isLoading, mutate } = useSWR(`${baseUrlMailNotif}?tenant_id=${tenant_id}`, fetchSWR)

  const columnsNotificationMail: ColumnDef<TFormNotificationMail>[] = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "from_name",
      header: "From Name",
    },
    {
      accessorKey: "from_email",
      header: "From Email",
    },
    {
      accessorKey: "smtp_host",
      header: "SMTP Host",
    },
    {
      accessorKey: "smtp_port",
      header: "SMTP Port",
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
          <span className={cn(
            "px-2 py-1 rounded text-xs font-medium",
            status === "active" ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"
          )}>
            {status}
          </span>
        );
      }
    },
    {
      accessorKey: "id",
      header: () => null,
      cell: ({ row }) => {
        const raw = row.original;

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
              <DropdownMenuItem onClick={() => {setOpenSheet(true); setDataForm(raw)}}>
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {setOpenTemplateDialog(true); setDataForm(raw)}}>
                Edit Layout
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {
                toast.promise(
                  fetch(`${BackendUrlBase}/publish?p=mail`, {
                    method: "POST",
                    credentials: "include",
                    signal: AbortSignal.timeout(10000)
                  }).then(async (res) => {
                    if (!res.ok) throw new Error("Request failed");
                    const data = await res.json();
                    if (!data.status) throw new Error(data.message || "Failed to test");
                    return data;
                  }),
                  {
                    loading: "Sending test email...",
                    success: () => "Test email sent successfully",
                    error: (err) => err.message || "Failed to send test email",
                  }
                );
              }}>
                Test
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem variant="destructive" onClick={() => {setOpenDialog(true); setDataForm(raw)}}>
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
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
          <Button onClick={() => {setOpenSheet(true);setDataForm(undefined)}} variant="outline">
            <Plus />
            NOTIFICATION
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4">
        {isLoading ? <TableSkeleton /> : <CardInformation name="SMTP Setup" description="There is list of SMTP Email for Notification" rowIdKey="email" columnsDetail={columnsNotificationMail} data={data?.data ?? []} />}
      </div>

      <SheetNotificationMailForm 
        openSheet={openSheet} 
        setOpenSheet={setOpenSheet} 
        dataForm={dataForm} 
        mutateData={mutate}
        />

      <NotificationMailTemplateDialog
        openTemplateDialog={openTemplateDialog}
        setOpenTemplateDialog={setOpenTemplateDialog}
        dataForm={dataForm}
        setDataForm={setDataForm}
        mutateData={mutate}
      />

      <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {dataForm?.from_email} Notification?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will delete the data. Data deleted cannot be restored.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction>Cancel</AlertDialogAction>
            <AlertDialogCancel variant="destructive" onClick={() => {
              if (dataForm?.id) {
                toast.promise(
                  fetch(`${baseUrlMailNotif}/${dataForm.id}`, {
                    method: "DELETE",
                    credentials: "include",
                    signal: AbortSignal.timeout(10000)
                  }).then(async (res) => {
                    if (!res.ok) throw new Error("Request failed");
                    const data = await res.json();
                    if (!data.status) throw new Error(data.message || "Failed to test");
                    return data;
                  }),
                  {
                    loading: "Waiting...",
                    success: () => {
                      mutate();
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