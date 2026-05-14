import { useState } from "react"
import { Filter, FilterX } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import DataTable from "@/components/data-table";

type TCardInformationProps = {
  name: string
  description?: string
  rowIdKey?: string
  /* eslint-disable @typescript-eslint/no-explicit-any */
  columnsDetail?: any[]
  /* eslint-disable @typescript-eslint/no-explicit-any */
  data?: any[]
}

export default function CardInformation({
  name,
  description,
  rowIdKey,
  columnsDetail,
  data,
}: TCardInformationProps) {
  const [openFilter, setOpenFilter] = useState(false)

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{name}</CardTitle>
        <CardDescription className={description || "sr-only"}>{description ? description : "no description"}</CardDescription>
        {rowIdKey && columnsDetail && data && <CardAction>
          <Button onClick={() => setOpenFilter(!openFilter)} variant={openFilter ? "outline" : "default"} size="icon">
            {openFilter ? <FilterX /> : <Filter />}
          </Button>
        </CardAction>}
      </CardHeader>
      <CardContent>
        {rowIdKey && columnsDetail && data && <DataTable rowIdKey={rowIdKey} columns={columnsDetail} data={data} useFilter={openFilter} />}
      </CardContent>
    </Card>
  )
}