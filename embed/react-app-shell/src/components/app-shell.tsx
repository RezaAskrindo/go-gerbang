
import { 
  lazy,
  Suspense, 
  useCallback, 
  type ComponentType, 
  type JSX, 
  type ReactNode 
} from "react"
import { Moon, Sun } from "lucide-react"

import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"

import { useTheme } from "@/components/useTheme";

import { AppSidebar } from "./app-sidebar"

const LoginView = lazy(() => import("./auth/LoginView"));

type AppShellProps = {
  Breadcrumb?: ComponentType
  TeamSwitcher?: ComponentType
  NavMain?: ComponentType
  CardInformation?: ComponentType
  NavUser?: ComponentType
  PageContent?: () => ReactNode
  AuthChecking?: boolean
  AuthPassed?: boolean
  ImageLogo?: string
  ImageLogoWhite?: string
  ImageBanner?: string
  loginSend?: (valuez:{ identity: string; password: string }) => void
  HeaderLogin?: JSX.Element
  FooterLogin?: ComponentType
  ResetPasswordForm?: ComponentType
}

export default function AppShell({
  Breadcrumb,
  TeamSwitcher,
  NavMain,
  CardInformation,
  NavUser,
  PageContent,
  AuthChecking,
  AuthPassed=false,
  ImageLogo,
  ImageLogoWhite,
  ImageBanner,
  loginSend,
  HeaderLogin,
  FooterLogin,
  ResetPasswordForm,
}: AppShellProps) {
  const { setTheme, theme } = useTheme();
  
  const toggleTheme = useCallback(() => {
    setTheme(theme === "dark" ? "light" : "dark");
  }, [theme, setTheme]);

  const loadingIndicator = <div className="fixed inset-0 flex items-center justify-center">
    <div className="relative w-10 h-10 border-2 border-black/70 border-b-transparent rounded-full animate-spin"></div>
  </div>

  if (AuthChecking) {
    return loadingIndicator;
  }

  if (!AuthPassed) {
    return <Suspense fallback={loadingIndicator}>
      <LoginView 
        ImageLogo={ImageLogo}
        ImageLogoWhite={ImageLogoWhite}
        ImageBanner={ImageBanner}
        loginSend={loginSend}
        HeaderLogin={HeaderLogin}
        FooterLogin={FooterLogin}
        ResetPasswordForm={ResetPasswordForm}
      />
    </Suspense>
  }

  return (
    <SidebarProvider>
      <AppSidebar 
        TeamSwitcher={TeamSwitcher}
        NavMain={NavMain}
        CardInformation={CardInformation}
        NavUser={NavUser}
      />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 border-b">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mr-2 data-[orientation=vertical]:h-4"
            />
            {Breadcrumb && (
              <Suspense fallback={<Skeleton className="h-[20px] w-[200px]" />}>
                <Breadcrumb />
              </Suspense>
            )}
          </div>
          <div className="ms-auto">
            <Button
              variant="outline"
              size="icon"
              className="group/toggle size-8 me-4"
              onClick={toggleTheme}
            >
              <Sun className="hidden [html.dark_&]:block" />
              <Moon className="hidden [html.light_&]:block" />
              <span className="sr-only">Toggle theme</span>
            </Button>
          </div>
        </header>
        <div className="flex-1 p-4">
          {PageContent && <PageContent />}
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
