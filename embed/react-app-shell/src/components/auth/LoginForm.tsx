import { useState, type ComponentType, type JSX } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { EyeIcon, EyeOffIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { toast } from "sonner"

const formSchema = z.object({
  identity: z.string().min(2, {
    message: "Identity Wajib Di Isi",
  }),
  password: z.string().min(2, {
    message: "Password Wajib Di Isi",
  })
})

type LoginFormProps = React.ComponentProps<"form"> & {
  // loginSend?: (valuez:{ identity: string; password: string }) => void
  BackendUrlBase: string
  FrontendUrl: string
  HeaderLogin?: JSX.Element
  FooterLogin?: ComponentType
  ResetPasswordForm?: ComponentType
}

export default function LoginForm({
  // loginSend,
  BackendUrlBase,
  FrontendUrl,
  HeaderLogin,
  FooterLogin,
  ResetPasswordForm,
  ...props
}: LoginFormProps) {
  
  const [showPassword, setShowPassword] = useState(false);
  const [loadingButton, setLoadingButton] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      identity: "",
      password: "",
    },
  })


  async function FetchCsrfToken(): Promise<string> {
    const urlCsrf = `${BackendUrlBase}/secure-gateway-c`;
    const res = await fetch(urlCsrf, { credentials: 'include' });
    if (!res.ok) return "";
    const data = await res.json();
    return data?.data;
  }

  function onSubmit(values: z.infer<typeof formSchema>) {
    // if (!loginSend) return true;
    
    setLoadingButton(true);

    const doLogin = async () => {
      const getCsrf = await FetchCsrfToken();
      const url = `${BackendUrlBase}/api/v1/auth/login?httponly=true&session=true&domain=localhost&url=${FrontendUrl}`;

      const res = await fetch(url, {
        method: "POST",
        credentials: "include",
        headers: { 
          "Content-Type": "application/json", 
          "X-SGCsrf-Token": getCsrf 
        },
        body: JSON.stringify(values),
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

    toast.promise(
      doLogin().catch(async (err) => {
        if (err.message === "CSRF validation failed") {
          const retryData = await doLogin();
          return retryData;
        }
        throw err;
      }),
      {
        loading: "Waiting...",
        success: () => {
          window.location.href = FrontendUrl;
          return "Success"
        },
        error: (err) => {
          return err.message || "Failed"
        },
      }
    )

    setTimeout(() => {
      setLoadingButton(false);
    }, 2000)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-6" {...props}>

        {HeaderLogin ?? null}
        
        <div className="grid gap-6">
          <FormField
            control={form.control}
            name="identity"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Identity</FormLabel>
                <FormControl>
                  <Input placeholder="name@example.com" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  <span>Password</span>
                  {ResetPasswordForm && <ResetPasswordForm />}
                </FormLabel>
                <FormControl>
                  <div className="relative">
                    <Input
                      type={showPassword ? 'text' : 'password'}
                      className="hide-password-toggle pr-10"
                      {...field}
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                      onClick={() => setShowPassword((prev) => !prev)}
                    >
                      {showPassword ? (
                        <EyeIcon className="h-4 w-4" aria-hidden="true" />
                      ) : (
                        <EyeOffIcon className="h-4 w-4" aria-hidden="true" />
                      )}
                      <span className="sr-only">{showPassword ? 'Hide password' : 'Show password'}</span>
                    </Button>

                    <style>{`
                        .hide-password-toggle::-ms-reveal,
                        .hide-password-toggle::-ms-clear {
                          visibility: hidden;
                          pointer-events: none;
                          display: none;
                        }
                      `}</style>
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          {loadingButton ? <Button type="button" variant="ghost" className="w-full" disabled={true}>
            Loading...
          </Button> : <Button type="submit" className="w-full">
            Login
          </Button>}

          {FooterLogin && <FooterLogin />}

        </div>
      </form>
    </Form>
  )
}
