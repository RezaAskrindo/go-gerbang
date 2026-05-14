import { type ComponentType, useState } from "react";
import { useTheme } from "@/components/useTheme";

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { ChevronRight, Lock, LockOpen } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { cn } from "@/lib/utils";

const formSchema = z.object({
  identity: z.string().min(2, {
    message: "Identity Wajib Di Isi",
  }),
  password: z.string().min(2, {
    message: "Password Wajib Di Isi",
  })
})

type LoginFormProps = React.ComponentProps<"form"> & {
  loginSend?: (valuez:{ identity: string; password: string }) => void
  HeaderLogin?: ComponentType
  FooterLogin?: ComponentType
  ResetPasswordForm?: ComponentType
}

function LoginForm({
  loginSend,
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

  function onSubmit(values: z.infer<typeof formSchema>) {
    if (!loginSend) return true;
    
    setLoadingButton(true);

    loginSend(values);

    setTimeout(() => {
      setLoadingButton(false);
    }, 1500)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-6" {...props}>

        {HeaderLogin && <HeaderLogin />}
        
        <div className="grid gap-6">
          <FormField
            control={form.control}
            name="identity"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Identity</FormLabel>
                <FormControl>
                  <Input placeholder="name@example.com" className={cn(field.value && 'border-green-700 focus:ring-green-700/40!')} {...field} />
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
                      className={cn(`hide-password-toggle pr-10 ${field.value && 'border-green-700 focus:ring-green-700/40!'}`)}
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
                        <LockOpen className="h-4 w-4" aria-hidden="true" />
                      ) : (
                        <Lock className="h-4 w-4" aria-hidden="true" />
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

          <Button variant="outline" type={loadingButton ? 'button' : 'submit'} className="w-full" disabled={loadingButton ? true : false}>
            {loadingButton ? <Spinner /> : <ChevronRight />}
            {loadingButton ? 'Loading...' : 'Login' }
          </Button>

          {FooterLogin && <FooterLogin />}

        </div>
      </form>
    </Form>
  )
}


interface LoginProps {
  ImageLogo?: string
  ImageLogoWhite?: string
  ImageBanner?: string
  loginSend: (valuez:{ identity: string; password: string }) => void
  HeaderLogin?: ComponentType
  FooterLogin?: ComponentType
  ResetPasswordForm?: ComponentType
}

export default function AppLogin({
  ImageLogo,
  ImageLogoWhite,
  ImageBanner,
  loginSend,
  HeaderLogin,
  FooterLogin,
  ResetPasswordForm,
}: LoginProps) {
  const { theme } = useTheme();

  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-start">
          {ImageLogo && ImageLogoWhite ? <img
            src={theme === "dark" && ImageLogoWhite ? ImageLogoWhite : ImageLogo}
            alt="Image"
            className="h-12"
          /> : "LOGO"}
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            <LoginForm 
              loginSend={loginSend} 
              HeaderLogin={HeaderLogin} 
              FooterLogin={FooterLogin} 
              ResetPasswordForm={ResetPasswordForm} 
            />
          </div>
        </div>
      </div>
      {ImageBanner && <div className="bg-muted relative hidden lg:block">
        <img
          src={ImageBanner}
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover"
        />
      </div>}
    </div>
  )
}