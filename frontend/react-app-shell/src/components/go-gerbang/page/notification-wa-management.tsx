import z from "zod";
import useSWR from "swr";

import { useEffect, useRef, useState } from "react";
import { Ban, EllipsisVertical, Plus, Save, Loader } from "lucide-react";
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
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox"
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

import CardInformation from "@/components/card-information";

import { fetchSWR } from "@/services/use-swr-service";
import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService";
import TableSkeleton from "../table-skeleton";
import type { Account } from "./user-management";

const baseUrlWaNotif = `${BackendUrlBase}/api/notification/whatsapp/accounts`

type TFormNotificationWa = {
  id: string
  tenant_id: string
  phone: string
  status: string
  created_at: string
  updated_at: string
}

const formSchema = z.object({
  tenant_id: z.string().min(1, {message: "Tenant is required"}),
  phone: z.string().min(10, {message: "Valid phone number required"}),
  status: z.enum(["active", "inactive"]),
});

type TFormNotificationWaSchema = z.infer<typeof formSchema>;

export function SheetNotificationWaForm({
  openSheet,
  setOpenSheet,
  dataForm,
  setDataForm,
  mutateData,
}: {
  openSheet: boolean;
  setOpenSheet: (open: boolean) => void;
  dataForm?: TFormNotificationWa;
  setDataForm?: (data: TFormNotificationWa | undefined) => void;
  mutateData: () => void
}) {
  const sheetContentRef = useRef<HTMLDivElement>(null);

  const form = useForm<TFormNotificationWaSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      tenant_id: "",
      phone: "",
      status: "active",
    },
  });

  const { data: tenantsData, isLoading: tenantsLoading } = useSWR(
    `${BackendUrlBase}/users/all`,
    fetchSWR
  );
  const tenants: Account[] = tenantsData?.items || [];

  async function onSubmit(values: z.infer<typeof formSchema>) {
    const urlSubmit = dataForm ? `${baseUrlWaNotif}/${dataForm.id}` : baseUrlWaNotif;
    const method = dataForm ? "PUT" : "POST";

    const getCsrf = await FetchCsrfToken();
            
    try {
      const response = await (await fetch(urlSubmit, {
        method: method,
        credentials: "include",
        headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
        body: JSON.stringify(values),
      })).json();

      if (response.success || response.status) {
        form.reset();
        setOpenSheet(false);
        mutateData();
        setDataForm?.(undefined);
        toast.success(`Success ${dataForm ? "Update" : "Create"} WhatsApp Account`);
      } else {
        toast.error(response.message || `Failed ${dataForm ? "Update" : "Create"} WhatsApp Account`);
      }
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "An error occurred");
    }
  }

  useEffect(() => {
    if (dataForm) {
      form.setValue("phone", dataForm.phone);
      form.setValue("status", dataForm.status as "active" | "inactive");
    } else {
      form.reset();
    }
  }, [dataForm, form]);

  return (
    <Sheet open={openSheet} onOpenChange={setOpenSheet}>
      <SheetContent ref={sheetContentRef}>
        <SheetHeader className="pb-0">
          <SheetTitle>
            {dataForm ? "Edit WhatsApp Account" : "Add WhatsApp Account"}
          </SheetTitle>
          <SheetDescription>
            Configure your WhatsApp account settings. Click save when you&apos;re done.
          </SheetDescription>
        </SheetHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col h-full overflow-hidden">
            <div className="grid flex-1 auto-rows-min gap-6 px-4 overflow-auto pb-3">
              <FormField
                control={form.control}
                name="tenant_id"
                render={({ field }) => (
                  <FormItem className="flex flex-col">
                    <FormLabel>Tenant</FormLabel>
                    <Combobox
                      key={tenantsLoading ? 'loading' : 'loaded'}
                      items={tenants}
                      itemToStringLabel={(tenant: Account) => tenant.username}
                      itemToStringValue={(tenant: Account) => String(tenant.idAccount)}
                      value={tenants.find(
                        (tenant) => tenant.idAccount && tenant.idAccount === field.value
                      ) ?? null}
                      onValueChange={(tenant: Account | null) => {
                        field.onChange(tenant?.idAccount ?? "");
                      }}
                    >
                      <ComboboxInput
                        placeholder="Search or select tenant..."
                        disabled={tenantsLoading}
                        className={cn(field.value && "border-green-700")}
                      />
                      <ComboboxContent>
                        <ComboboxEmpty>No items found.</ComboboxEmpty>
                        <ComboboxList>
                          {(tenant: Account) => (
                            <ComboboxItem key={tenant.idAccount} value={tenant}>
                              {tenant.username}
                            </ComboboxItem>
                          )}
                        </ComboboxList>
                      </ComboboxContent>
                    </Combobox>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Phone Number</FormLabel>
                    <FormControl>
                      <Input 
                        type="text" 
                        placeholder="e.g., +1234567890" 
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

interface WhatsAppQRResponse {
  success: boolean;
  message: string;
  data: string | null; // base64 PNG image
}

export function WhatsAppConnectDialog({
  openQRDialog,
  setOpenQRDialog,
  dataForm,
}: {
  openQRDialog: boolean;
  setOpenQRDialog: (open: boolean) => void;
  dataForm?: TFormNotificationWa;
}) {
  const [qrCode, setQrCode] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);

  const handleConnect = async () => {
    if (!dataForm?.id) return;

    setIsConnecting(true);
    try {
      // Step 1: Initiate connection
      const connectResponse = await fetch(`${baseUrlWaNotif}/${dataForm.id}/connect`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        signal: AbortSignal.timeout(10000),
      });

      if (!connectResponse.ok) {
        throw new Error("Failed to initiate connection");
      }

      const connectData = await connectResponse.json();
      if (!connectData.success) {
        throw new Error(connectData.message || "Connection failed");
      }

      setLoading(true);

      // Step 2: Get QR Code
      const qrResponse = await fetch(`${baseUrlWaNotif}/${dataForm.id}/qr`, {
        method: "GET",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        signal: AbortSignal.timeout(10000),
      });

      if (!qrResponse.ok) {
        throw new Error("Failed to fetch QR code");
      }

      const qrData: WhatsAppQRResponse = await qrResponse.json();
      if (!qrData.success || !qrData.data) {
        throw new Error(qrData.message || "Failed to generate QR code");
      }

      setQrCode(qrData.data);
      toast.success("QR Code generated successfully");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Connection failed");
    } finally {
      setIsConnecting(false);
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!openQRDialog) {
      setQrCode(null);
      setLoading(false);
      setIsConnecting(false);
    }
  }, [openQRDialog]);

  return (
    <Dialog open={openQRDialog} onOpenChange={setOpenQRDialog}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Connect WhatsApp Account</DialogTitle>
          <DialogDescription>
            Scan the QR code with your WhatsApp device to connect.
          </DialogDescription>
        </DialogHeader>

        <div className="flex flex-col gap-4 py-6">
          {qrCode ? (
            <div className="flex justify-center">
              <img 
                src={`data:image/png;base64,${qrCode}`} 
                alt="WhatsApp QR Code"
                className="w-64 h-64 border rounded-lg"
              />
            </div>
          ) : loading ? (
            <div className="flex flex-col items-center justify-center gap-2 py-12">
              <Loader className="w-8 h-8 animate-spin text-muted-foreground" />
              <p className="text-sm text-muted-foreground">Generating QR code...</p>
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center gap-2 py-12">
              <p className="text-sm text-muted-foreground">
                Click "Generate QR Code" to start the connection process.
              </p>
            </div>
          )}
        </div>

        <div className="flex flex-row gap-2">
          <DialogClose asChild>
            <Button className="flex-1" variant="outline">
              <Ban className="w-4 h-4" />
              Close
            </Button>
          </DialogClose>
          <Button 
            className="flex-1" 
            onClick={handleConnect}
            disabled={isConnecting || loading}
          >
            {isConnecting || loading ? (
              <>
                <Loader className="w-4 h-4 animate-spin" />
                Connecting...
              </>
            ) : (
              "Generate QR Code"
            )}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default function WhatsAppNotificationManagement() {
  const [openSheet, setOpenSheet] = useState(false);
  const [openDialog, setOpenDialog] = useState(false);
  const [openQRDialog, setOpenQRDialog] = useState(false);
  const [dataForm, setDataForm] = useState<TFormNotificationWa | undefined>();

  const { data, isLoading, mutate } = useSWR(baseUrlWaNotif, fetchSWR)

  const columnsNotificationWa: ColumnDef<TFormNotificationWa>[] = [
    // {
    //   accessorKey: "id",
    //   header: "ID",
    // },
    {
      accessorKey: "phone",
      header: "Phone Number",
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
          <span className={cn(
            "px-2 py-1 rounded text-xs font-medium",
            status === "connected" ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"
          )}>
            {status}
          </span>
        );
      }
    },
    {
      accessorKey: "created_at",
      header: "Created",
      cell: ({ row }) => {
        const date = new Date(row.getValue("created_at") as string);
        return date.toLocaleDateString();
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
            <DropdownMenuContent align="end" className="w-40">
              <DropdownMenuItem onClick={() => {setOpenSheet(true); setDataForm(raw)}}>
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {setOpenQRDialog(true); setDataForm(raw)}}>
                Connect WhatsApp
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {
                toast.promise(
                  fetch(`${BackendUrlBase}/publish?p=wa`, {
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
                    loading: "Sending test wa...",
                    success: () => "Test wa sent successfully",
                    error: (err) => err.message || "Failed to send test wa",
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
            WhatsApp Notification Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s list of WhatsApp accounts for notification
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => {setOpenSheet(true); setDataForm(undefined)}} variant="outline">
            <Plus />
            ADD ACCOUNT
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4">
        {isLoading ? <TableSkeleton /> : <CardInformation name="WhatsApp Accounts" description="List of WhatsApp accounts for sending notifications" rowIdKey="phone" columnsDetail={columnsNotificationWa} data={data?.data ?? []} />}
      </div>

      <SheetNotificationWaForm 
        openSheet={openSheet} 
        setOpenSheet={setOpenSheet} 
        dataForm={dataForm}
        setDataForm={setDataForm}
        mutateData={mutate}
      />

      <WhatsAppConnectDialog
        openQRDialog={openQRDialog}
        setOpenQRDialog={setOpenQRDialog}
        dataForm={dataForm}
      />

      <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {dataForm?.phone} Account?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will delete the WhatsApp account. This cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction 
              variant="destructive" 
              onClick={() => {
                if (dataForm?.id) {
                  toast.promise(
                    fetch(`${baseUrlWaNotif}/${dataForm.id}`, {
                      method: "DELETE",
                      credentials: "include",
                      signal: AbortSignal.timeout(10000)
                    }).then(async (res) => {
                      if (!res.ok) throw new Error("Request failed");
                      const data = await res.json();
                      if (!data.success && !data.status) {
                        throw new Error(data.message || "Failed to delete");
                      }
                      return data;
                    }),
                    {
                      loading: "Deleting...",
                      success: () => {
                        mutate();
                        setOpenDialog(false);
                        return "Account deleted successfully";
                      },
                      error: (err) => {
                        setOpenDialog(false);
                        return err.message || "Failed to delete";
                      },
                    } 
                  )
                }
              }}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}