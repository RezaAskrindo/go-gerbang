
import { 
  Suspense, 
  useCallback, 
  type ComponentType, 
  type ComponentProps, 
  type ReactNode 
} from "react"
import { Moon, Sun } from "lucide-react"

import { Separator } from "@/components/ui/separator"
import {
  Sidebar,
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"

import { useTheme } from "@/components/useTheme";

import { AppSidebar } from "./app-sidebar"

type AppShellProps = ComponentProps<typeof Sidebar> & {
  Breadcrumb?: ComponentType
  TeamSwitcher?: ComponentType
  NavMain?: ComponentType
  CardInformation?: ComponentType
  NavUser?: ComponentType
  PageContent?: () => ReactNode
}

export default function AppShell({
  Breadcrumb,
  TeamSwitcher,
  NavMain,
  CardInformation,
  NavUser,
  PageContent,
  ...props
}: AppShellProps) {
  const { setTheme, theme } = useTheme();
  
  const toggleTheme = useCallback(() => {
    setTheme(theme === "dark" ? "light" : "dark");
  }, [theme, setTheme]);

  return (
    <SidebarProvider>
      <AppSidebar 
        TeamSwitcher={TeamSwitcher}
        NavMain={NavMain}
        CardInformation={CardInformation}
        NavUser={NavUser}
        {...props}
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
        <div className="flex flex-1 flex-col">
          <div className="@container/main flex flex-1 flex-col gap-2">
            <div className="p-4">
              {PageContent && <PageContent />}
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
