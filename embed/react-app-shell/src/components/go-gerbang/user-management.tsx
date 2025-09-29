import useSWR from 'swr'

import { type ColumnDef } from "@tanstack/react-table";

import DataTable from '@/components/data-table';
import { Button } from '@/components/ui/button';

import { BackendUrlBase } from '@/services/baseService';

import TableSkeleton from './table-skeleton';
import { ArrowUpDown } from 'lucide-react';

type Account = {
  idAccount: string
  identityNumber: string
  username: string
  fullName: string
  email: string
  phoneNumber: string
  loginTime: number
  statusAccount: number
}

const columns: ColumnDef<Account>[] = [
  {
    accessorKey: "identityNumber",
    header: "ID Number",
  },
  {
    accessorKey: "username",
    header: ({ column }) => {
      return (
        <Button
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          variant="ghost"
          className="w-full"
        >
          Email
          <ArrowUpDown />
        </Button>
      )
    },
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
  },
  {
    accessorKey: "statusAccount",
    header: "Status",
  },
];

export default function UserManagement() {
  const { data, isLoading } = useSWR(
    `${BackendUrlBase}/users/all`, 
    (url: string) => fetch(url).then(r => r.json()),
    {
      revalidateOnFocus: false
    }
  );

  if (isLoading) return <TableSkeleton />;

  return (
    <div>
      {data?.items?.length ? <DataTable columns={columns} data={data?.items} /> : <TableSkeleton />}
    </div>
  )
}