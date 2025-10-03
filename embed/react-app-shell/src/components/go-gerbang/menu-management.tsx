import useSWR from "swr";
import { Edit, Plus } from "lucide-react"
import { 
  BackendUrlBase, 
  fetchSWR
} from "@/services/baseService";


import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"

import { useState } from "react";
import { toast } from "sonner";

type TAuthRole = {
  id_auth_role: number
  name_auth_role: string
  desc_auth_role: string | null
}

type AuthRoleProps = {
  value?: TAuthRole
}

const AuthRoleForm = ({ value }: AuthRoleProps) => {
  const [form, setForm] = useState({
    name_auth_role: value?.name_auth_role || "",
    desc_auth_role: value?.desc_auth_role || "",
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const method = value?.id_auth_role ? "PUT" : "POST"
    const endpoint = value?.id_auth_role
      ? `${BackendUrlBase}/auth/role/${value.id_auth_role}`
      : `${BackendUrlBase}/auth/role`

    const res = await fetch(endpoint, {
      method,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(form),
    })

    if (!res.ok) {
      toast.error((await res.json())?.message ?? "Failed to save")
    } else {
      toast.success((await res.json())?.message ?? "Failed to save");
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="icon">
          {value ? <Edit /> : <Plus />}
        </Button>
      </DialogTrigger>

      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Auth Role</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4">
            <div className="grid gap-3">
              <Label htmlFor="name_auth_role">Name</Label>
              <Input 
                id="name_auth_role"
                name="name_auth_role"
                value={form.name_auth_role}
                onChange={handleChange}
                required
              />
            </div>
            <div className="grid gap-3">
              <Label htmlFor="desc_auth_role">Description</Label>
              <Textarea
                id="desc_auth_role"
                name="desc_auth_role"
                value={form.desc_auth_role}
                onChange={handleChange}
              />
            </div>
          </div>
          <DialogFooter className="mt-3">
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button type="submit" variant="outline">
              {value ? "Update" : "Create"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

const AuthRoleList = () => {
  const { data } = useSWR(`${BackendUrlBase}/auth/role/all`, fetchSWR)

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="mt-3">Auth Role</CardTitle>
        <CardAction>
          <AuthRoleForm />
        </CardAction>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">Rule Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead className="text-right">#</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data?.data?.length ? data.data?.map((el: TAuthRole) => <TableRow key={el.id_auth_role}>
              <TableCell>{el.name_auth_role}</TableCell>
              <TableCell>{el.desc_auth_role}</TableCell>
              <TableCell className="text-right"><AuthRoleForm value={el} /></TableCell>
            </TableRow>) : <TableRow><TableCell colSpan={3} className="h-24 text-center">No Data</TableCell></TableRow>}
          </TableBody>
        </Table>

      </CardContent>
    </Card>
  )
}

const PolicyList = () => {

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Policy List</CardTitle>
        {/* <CardAction>
          <Button variant="link" size="icon">
            <RefreshCcw />
          </Button>
        </CardAction> */}
      </CardHeader>
      <CardContent>
        <Table>
          <TableCaption>A list of your recent invoices.</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">Rule Type</TableHead>
              <TableHead>Tenant</TableHead>
              <TableHead>Method</TableHead>
              <TableHead className="text-right">Amount</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>

          </TableBody>
        </Table>

      </CardContent>
    </Card>
  )
}


export default function MenuManagement() {

  return (
    <div className="grid grid-cols-3 gap-4">
      <div>
        <AuthRoleList />
      </div>
      <div className="col-span-2">
        <PolicyList />
      </div>
    </div>
  )
}