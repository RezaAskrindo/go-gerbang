

import { zodResolver } from "@hookform/resolvers/zod"
import { Controller, useForm } from "react-hook-form"
import * as z from "zod"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"

import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService"
import { Lock, Unlock } from "lucide-react"
import { useState } from "react"

const formSchema = z.object({
  password: z
    .string()
    .min(4, "Password must be at least 4 characters."),
})

export default function AppBeginning({
  tables,
  missing,
}: {
  tables: string[]
  missing: string[]
}) {
  const [showPassword, setShowPassword] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      password: ""
    },
  })

  async function onSubmit(data: z.infer<typeof formSchema>) {
    const getCsrf = await FetchCsrfToken();

    const res = await (await fetch(`${BackendUrlBase}/migration-admin`, {
      method: "POST",
      credentials: "include",
      headers: { 
        "Content-Type": "application/json", 
        "X-SGCsrf-Token": getCsrf 
      },
      body: JSON.stringify(data),
    })).json();

    if (res.status) {
      toast.success(res.message || "Success to generate Admin");
    } else {
      toast.error(res.message || "Failed to generate Admin");
    }
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <Card>
          <CardHeader>
            <CardTitle>Check Migration</CardTitle>
            <CardDescription className="sr-only">Check Migration</CardDescription>
          </CardHeader>
          <CardContent>
            {tables[0] === "admin" ? <form id="form-rhf-demo" onSubmit={form.handleSubmit(onSubmit)}>
              <FieldGroup className="gap-y-3">
                <Field className="flex-row">
                  <FieldLabel>Username</FieldLabel>
                  <div>: admin</div>
                </Field>
                <Field className="flex-row">
                  <FieldLabel>Fullname</FieldLabel>
                  <div>: Admin</div>
                </Field>
                <Controller
                  name="password"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="form-password">
                        Password
                      </FieldLabel>
                      <InputGroup>
                        <InputGroupInput 
                          {...field}
                          type={showPassword ? "text" : "password"}
                          id="form-password"
                          aria-invalid={fieldState.invalid}
                          placeholder="Password Admin"
                          autoComplete="off"
                        />
                        <InputGroupAddon onClick={() => setShowPassword(!showPassword)} className="cursor-pointer" align="inline-end">
                          {showPassword ? <Unlock className="text-muted-foreground" /> : <Lock className="text-muted-foreground" />}
                        </InputGroupAddon>
                      </InputGroup>
                      <FieldDescription>Please Save Your Password</FieldDescription>
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </Field>
                  )}
                />
              </FieldGroup>
            
            </form> : 
            tables.map(el => missing.indexOf(el) > -1 && <div key={el}><strong>{el}</strong> is not exist</div>)}
          </CardContent>
          <CardFooter className="grid w-full">
            {tables[0] === "admin" ? <Button type="submit" form="form-rhf-demo">Submit</Button>:<Button variant="outline" onClick={() => {
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
            }}>Migration</Button>}
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}