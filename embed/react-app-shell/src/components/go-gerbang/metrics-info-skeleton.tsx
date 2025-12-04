
import { Skeleton } from "@/components/ui/skeleton";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function MetricsInfoSkeleton() {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      { [0,1,2,3].map(el => <Card key={el} className="@container/card">
        <CardHeader>
          <CardDescription>
            <Skeleton className="h-6 w-32" />
          </CardDescription>
          <CardTitle className="text-2xl font-bold tabular-nums @[250px]/card:text-3xl">
            <Skeleton className="h-10 w-full" />
          </CardTitle>
        </CardHeader>
        <CardFooter>
          <Skeleton className="h-48 w-full" />
        </CardFooter>
      </Card>) }
    </div>
  )
}