// import { type ComponentType, type JSX } from "react";

import { useTheme } from "@/components/useTheme";
// import LoginForm from "./LoginForm";

interface LoginProps {
  ImageLogo?: string
  ImageLogoWhite?: string
  ImageBanner?: string
  // loginSend: (valuez:{ identity: string; password: string }) => void
  // HeaderLogin?: JSX.Element
  // FooterLogin?: ComponentType
  // ResetPasswordForm?: ComponentType
}

export default function LoginView({
  ImageLogo,
  ImageLogoWhite,
  ImageBanner,
  // loginSend,
  // HeaderLogin,
  // FooterLogin,
  // ResetPasswordForm,
}: LoginProps) {
  const { theme } = useTheme();

  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-start">
          {ImageLogo && <img
            src={theme === "dark" && ImageLogoWhite ? ImageLogoWhite : ImageLogo}
            alt="Image"
            className="h-12"
          />}
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            {/* <LoginForm 
              loginSend={loginSend} 
              HeaderLogin={HeaderLogin} 
              FooterLogin={FooterLogin} 
              ResetPasswordForm={ResetPasswordForm} 
            /> */}
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