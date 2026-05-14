import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { toast } from "sonner"
import { BackendUrlBase } from "@/services/baseService"


export default function AppBeginning({
  tables,
  missing,
}: {
  tables: string[]
  missing: string[]
}) {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <Card>
          <CardHeader>
            <CardTitle>Check Migration</CardTitle>
            <CardDescription>
              Enter your email below to login to your account
            </CardDescription>
          </CardHeader>
          <CardContent>
            {tables.map(el => missing.indexOf(el) > -1 && <div key={el}><strong>{el}</strong> is not exist</div>)}
          </CardContent>
          <CardFooter className="flex-col gap-2">
            <Button onClick={() => {
              toast.promise(
                fetch(`${BackendUrlBase}/migration`, {
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
                    window.location.reload();
                    return "Success"
                  },
                  error: (err) => {
                    return err.message || "Failed"
                  },
                }
              )
            }}>
              Migration
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}