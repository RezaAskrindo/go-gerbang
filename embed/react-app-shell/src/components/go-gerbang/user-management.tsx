import { useEffect, useState, type FC } from 'react';
import useSWR, { type KeyedMutator } from 'swr'
import { Ban, Ellipsis, EyeIcon, EyeOff, Plus, Save } from 'lucide-react';

import { type ColumnDef } from "@tanstack/react-table";

import DataTable from '@/components/data-table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';

import { BackendUrlBase, FetchCsrfToken } from '@/services/baseService';

import TableSkeleton from './table-skeleton';
import { toast } from 'sonner';

type FormProps<T> = {
  value?: T
  setOpen?: (val: boolean) => void
  mutate?: KeyedMutator<Account[]>;
}

type ModalFormProps = {
  typeForm: "user" | "password"
  value?: Account | AccountPassword
  open: boolean
  setOpen: (val: boolean) => void
  mutate?: KeyedMutator<Account[]>;
}

const ModalForm = ({typeForm, value, open, setOpen, mutate}: ModalFormProps) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const forms: Record<string, FC<FormProps<any>>> = {
    user: UserForm,
    password: UserPasswordForm
  }

  const FormComponent = forms[typeForm];

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-4xl">
        <DialogHeader>
          <DialogTitle>{typeForm.toUpperCase()} Form</DialogTitle>
          <DialogDescription className="sr-only">...</DialogDescription>
        </DialogHeader>
        {FormComponent && <FormComponent value={value} setOpen={setOpen} mutate={mutate} />}
      </DialogContent>
    </Dialog>
  )
}

const required = <span className="text-red-500 -ms-1.5">*</span>

type Account = {
  idAccount?: string
  identityNumber: string
  username: string
  fullName: string
  email: string
  phoneNumber: string
  statusAccount: number
  loginTime?: number
}

const UserForm = ({value, setOpen, mutate}: FormProps<Account>) => {
  const [form, setForm] = useState<Account>({
    idAccount: value?.idAccount || "",
    identityNumber: value?.identityNumber || "",
    username: value?.username || "",
    fullName: value?.fullName || "",
    email: value?.email || "",
    phoneNumber: value?.phoneNumber || "",
    statusAccount: value?.statusAccount || 0,
  });

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [error, setError] = useState<any>()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setLoading(true);
    
    const method = value?.idAccount ? "PUT" : "POST"
    const endpoint = value?.idAccount
      ? `${BackendUrlBase}/users/${value.idAccount}`
      : `${BackendUrlBase}/users`

    const { idAccount, ...rest } = form;
    const payload = idAccount === "" ? rest : form;

    const getCsrf = await FetchCsrfToken();

    const res = await fetch(endpoint, {
      method,
      credentials: "include",
      headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
      body: JSON.stringify(payload),
    })

    const response = await res.json();
    if (!res.ok) {
      toast.error(response?.message ?? "Failed to save")
    } else if (!response.status) {
      toast.info(response?.message ?? "Error to save");
      setError(response.data);
    } else {
      toast.success(response?.message ?? "Success to save");
    }

    setLoading(false);
    setOpen?.(false);
    mutate?.();
  }

  useEffect(() => {
    setForm({
      idAccount: value?.idAccount || "",
      identityNumber: value?.identityNumber || "",
      username: value?.username || "",
      fullName: value?.fullName || "",
      email: value?.email || "",
      phoneNumber: value?.phoneNumber || "",
      statusAccount: value?.statusAccount || 0,
    });

  }, [value])

  return (
    <form onSubmit={handleSubmit}>
      <div className="grid grid-cols-2 gap-4">
        <div className="grid gap-2">
          <Label htmlFor="username" className={`${error?.username ? "text-destructive" : ""}`}>Username{required}</Label>
          <Input id="username" name="username" value={form.username} onChange={handleChange} placeholder="Username" aria-invalid={error?.username ? "true" : "false"} required />
          <p className="text-destructive text-sm">{error?.username?.desc}</p>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="fullName" className={`${error?.fullName ? "text-destructive" : ""}`}>Full Name{required}</Label>
          <Input id="fullName" name="fullName" value={form.fullName} onChange={handleChange} placeholder="Nama Lengkap" aria-invalid={error?.fullName ? "true" : "false"} required />
          <p className="text-destructive text-sm">{error?.fullName?.desc}</p>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="identityNumber">No Identity</Label>
          <Input id="identityNumber" name="identityNumber" value={form.identityNumber} onChange={handleChange} placeholder="Nomor Identitas" />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input id="email" name="email" value={form.email} onChange={handleChange} placeholder="Alamat Email" type="email" />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input id="email" name="email" value={form.email} onChange={handleChange} placeholder="Alamat Email" type="email" />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">Status</Label>
          <Select value={form.statusAccount?.toString()} onValueChange={(e) => setForm({ ...form, statusAccount: parseInt(e) })}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Select a Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="9">Non Aktif</SelectItem>
                <SelectItem value="10">Aktif</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
      </div>
      <DialogFooter className="mt-6 flex !justify-between">
        <DialogClose asChild>
          <Button variant="outline">
            <Ban />
            Cancel
          </Button>
        </DialogClose>
        <Button type={loading? "button" : "submit"} disabled={loading}>
          <Save />
          {loading ? ("Loading...") : (value ? "Update" : "Create")}
        </Button>
      </DialogFooter>
    </form>
  )
}

type AccountPassword = {
  id?: string
  password: string
}

const UserPasswordForm = ({value, setOpen, mutate}: FormProps<Account>) => {
  const [form, setForm] = useState<AccountPassword>({
    password: "",
  })
 
  const [showPassword, setShowPassword] = useState(false);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [error, setError] = useState<any>()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    setLoading(true);
    
    const getCsrf = await FetchCsrfToken();

    const res = await fetch(`${BackendUrlBase}/api/v1/auth/change-password`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
      body: JSON.stringify(form),
    })

    const response = await res.json();
    if (!res.ok) {
      toast.error(response?.message ?? "Failed to save")
    } else if (!response.status) {
      toast.info(response?.message ?? "Error to save");
      setError(response.data);
    } else {
      toast.success(response?.message ?? "Success to save");
    }
    
    mutate?.();
    setOpen?.(false);
    setLoading(false);
  }

  useEffect(() => {
    setForm({
      id: value?.idAccount || "",
      password: ""
    });
  }, [value])

  return (
    <form onSubmit={handleSubmit}>
      <div className="grid gap-4">
        <div className="flex flex-row gap-2 my-4">
          <Label htmlFor="password" className={`${error?.password ? "text-destructive" : ""}`}>Password{required}</Label>
          <div className="relative w-full">
            <Input id="password" name="password" value={form.password} onChange={handleChange} placeholder="Password" aria-invalid={error?.password ? "true" : "false"} type={showPassword ? "text" : "password"} required />
            <Button type="button" variant="ghost" size="sm" className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent" onClick={() => setShowPassword((prev) => !prev)}
            >
              {showPassword ? (
                <EyeIcon className="h-4 w-4" aria-hidden="true" />
              ) : (
                <EyeOff className="h-4 w-4" aria-hidden="true" />
              )}
              <span className="sr-only">{showPassword ? 'Hide password' : 'Show password'}</span>
            </Button>
          </div>
          <p className="text-destructive text-sm">{error?.password?.desc}</p>
        </div>
      </div>
      <DialogFooter className="mt-6 flex !justify-between">
        <DialogClose asChild>
          <Button variant="outline">
            <Ban />
            Cancel
          </Button>
        </DialogClose>
        <Button type={loading? "button" : "submit"} disabled={loading}>
          <Save />
          {loading ? ("Loading...") : (value ? "Update" : "Create")}
        </Button>
      </DialogFooter>
    </form>
  )
}

export default function UserManagement() {
  const [formUser, setFormUser] = useState<Account | undefined>();
  const [open, setOpen] = useState(false);
  const [modalType, setModalType] = useState<"user" | "password">("user");

  const { data, isLoading, mutate } = useSWR(
    `${BackendUrlBase}/users/all`, 
    (url: string) => fetch(url).then(r => r.json()),
  );

  const doDelete = async (id: string) => {
    const getCsrf = await FetchCsrfToken();
    const url = `${BackendUrlBase}/users/${id}`;

    const res = await fetch(url, {
      method: "DELETE",
      credentials: "include",
      headers: { 
        "Content-Type": "application/json", 
        "X-SGCsrf-Token": getCsrf 
      },
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.message || "Request failed");
    }
    if (!data.status) {
      throw new Error(data.message || "Failed to login");
    }

    return data;
  };

  const columns: ColumnDef<Account>[] = [
    {
      accessorKey: "identityNumber",
      header: "ID Number",
    },
    {
      accessorKey: "username",
      header: "Username",
    },
    {
      accessorKey: "fullName",
      header: "Full Name",
    },
    {
      accessorKey: "email",
      header: "Email",
    },
    {
      accessorKey: "phoneNumber",
      header: "Phone Number",
    },
    {
      accessorKey: "loginTime",
      header: "Login Time",
      cell: ({ row }) => row.original.loginTime && new Date(row.original.loginTime).toLocaleTimeString()
    },
    {
      accessorKey: "statusAccount",
      header: "Status",
      cell: ({ row }) => row.original.statusAccount === 10 ? "Aktif" : "Non Aktif"
    },
    {
      id: "actions",
      cell: ({ row }) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              className="data-[state=open]:bg-muted text-muted-foreground flex size-8"
              size="icon"
            >
              <Ellipsis />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-32">
            <DropdownMenuItem onClick={() => {setFormUser(row.original);setOpen(true);setModalType("user")}}>Edit</DropdownMenuItem>
            <DropdownMenuItem onClick={() => {setFormUser(row.original);setOpen(true);setModalType("password")}}>Set Password</DropdownMenuItem>
            <DropdownMenuItem >Add Role</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem variant="destructive" onClick={() => toast.promise(
              doDelete(row.original.idAccount ?? "").catch(async (err) => {
                if (err.message === "CSRF validation failed") {
                  const retryData = await doDelete(row.original.idAccount ?? "");
                  return retryData;
                }
                throw err;
              }),
              {
                loading: "Waiting...",
                success: () => {
                  mutate();
                  return "Success"
                },
                error: (err) => {
                  return err.message || "Failed"
                },
              } 
            )}>Delete</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ];

  const rows = data?.items ? data.items?.map((item: Account) => ({
    ...item,
    id: item.idAccount
  })) : []
  
  return (
    <div className="flex-1 flex-col gap-8 md:flex">
      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-col gap-1">
          <h2 className="text-2xl font-semibold tracking-tight">
            User Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s a list of users.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => {setFormUser(undefined);setOpen(true);setModalType("user")}} variant="outline">
            <Plus />
            USER
          </Button>
        </div>
      </div>
      
      <ModalForm typeForm={modalType} open={open} setOpen={setOpen} value={formUser} mutate={mutate} />

      {isLoading ? <TableSkeleton /> : <DataTable rowIdKey="idAccount" columns={columns} data={rows} useFilter={true} />}
    </div>
  )
}