import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Skeleton } from "@/components/ui/skeleton"

export default function TableSkeleton() {
  const totalTable = 10;

  return (
    <div className="overflow-hidden rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            {Array.from({length: totalTable}).map((_, i) => <TableHead key={i} className="bg-muted sticky top-0 z-10">
              <Skeleton className="h-[20px] w-full bg-muted-foreground" />
            </TableHead>)}
          </TableRow>
        </TableHeader>
        <TableBody>
          {Array.from({length: totalTable*2}).map((_, i) => (
            <TableRow key={i}>
              {Array.from({length: totalTable}).map((_, i) => <TableCell key={i}>
                <Skeleton className="h-[20px] w-full" />
              </TableCell>)}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}