import { lazy, useMemo, type ReactNode } from "react";

import {
  ChevronRight,
  ChevronsUpDown,
  LogOut,
  Settings,
  ShieldPlus,
} from "lucide-react";

import { toast } from "sonner"

import { ThemeProvider } from "@/components/theme-provider";

import { useSidebar } from "@/components/ui/useHelper";
import { Toaster } from "@/components/ui/sonner";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";

import Placeholder from "@/assets/placeholder.svg";
import LogoReact from "@/assets/react.svg";

import NotFound from "./components/go-gerbang/not-found";
// import ResetPasswordForm from "./components/go-gerbang/reset-password-form";

import {
  LoginUser,
  LogoutUser,
} from "./services/baseService";
import { 
  CheckMigration, 
  GetAuthSession 
} from "./services/use-swr-service";

import AppLogin from "@/components/app-login";
import AppBeginning from "./components/go-gerbang/app-beginning";
import { useHash } from "./hooks/use-hash";

const AppShell = lazy(() => import("@/components/app-shell"));
const Dashboard = lazy(() => import("@/components/go-gerbang/dashboard"));

const CaddyManagement = lazy(() => import("@/components/go-gerbang/page/caddy-management"));
const LogManagement = lazy(() => import("@/components/go-gerbang/page/log-management"));
const ModuleManagement = lazy(() => import("@/components/go-gerbang/page/module-management"));
const NotificationManagement = lazy(() => import("@/components/go-gerbang/page/notification-management"));
const RbacManagement = lazy(() => import("@/components/go-gerbang/page/rbac-management"));
const UserManagement = lazy(() => import("@/components/go-gerbang/page/user-management"));

const TeamSwitcher = () => {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton className="data-[slot=sidebar-menu-button]:p-1.5!" asChild>
          <a href="#">
            <ShieldPlus className="size-5!" />
            <span className="text-base font-semibold">GO Gerbang</span>
          </a>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

const NavMain = () => {
  const hash = useHash();

  const items = [
    {
      title: "Configuration",
      url: "#",
      icon: Settings,
      isActive: true,
      items: [
        {title: "caddy", url: "#/caddy"},
        {title: "log", url: "#/log"},
        {title: "module", url: "#/module"},
        {title: "notification", url: "#/notification"},
        {title: "user", url: "#/user"},
        // {title: "rbac", url: "#/rbac"},
      ]
    }
  ]

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Go Gerbang Menu</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => (
          <Collapsible
            key={item.title}
            asChild
            defaultOpen={item.isActive}
            className="group/collapsible"
          >
            <SidebarMenuItem>
              <CollapsibleTrigger asChild>
                <SidebarMenuButton tooltip={item.title}>
                  {item.icon && <item.icon />}
                  <span>{item.title}</span>
                  <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                </SidebarMenuButton>
              </CollapsibleTrigger>
              <CollapsibleContent>
                <SidebarMenuSub>
                  {item.items?.map((subItem) => (
                    <SidebarMenuSubItem key={subItem.title}>
                      <SidebarMenuSubButton asChild isActive={subItem.url === hash}>
                        <a href={subItem.url}>
                          <span>{subItem.title}</span>
                        </a>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                  ))}
                </SidebarMenuSub>
              </CollapsibleContent>
            </SidebarMenuItem>
          </Collapsible>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}

const NavUser = () => {
  const { isMobile } = useSidebar();
  
  const { user } = useAuth();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage src={user?.avatar} alt={user?.fullName} />
                <AvatarFallback className="rounded-lg">CN</AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user?.fullName}</span>
                <span className="truncate text-xs">{user?.username}</span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarImage src={user?.avatar} alt={user?.fullName} />
                  <AvatarFallback className="rounded-lg">CN</AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user?.fullName}</span>
                  <span className="truncate text-xs">{user?.username}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            {/* <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <Bell />
                Notifications
              </DropdownMenuItem>
            </DropdownMenuGroup> */}
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={LogoutUser}>
              <LogOut />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

function useAuth() {
  const { data, isLoading } = GetAuthSession();

  return { auth: data?.status, user: data?.data, isLoading };
}

const routes: Record<string, React.ComponentType> = {
  "/": Dashboard,
  "/caddy": CaddyManagement,
  "/log": LogManagement,
  "/module": ModuleManagement,
  "/notification": NotificationManagement,
  "/rbac": RbacManagement,
  "/user": UserManagement,
};

function AppRouter(): ReactNode {
  const hash = useHash();

  const CurrentView = useMemo(() => {
    return routes[hash.slice(1) || "/"] || NotFound;
  }, [hash]);

  return <CurrentView />;
}

function BreadcrumbCom() {
  const hash = useHash();

  const hashMenu = hash.replace(/#\//g, '').replace(/^\w/, c => c.toUpperCase());

  return (
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem className="hidden md:block">
          <BreadcrumbLink href="#/">
            Dashboard
          </BreadcrumbLink>
        </BreadcrumbItem>
        {hash !== '#/' && <BreadcrumbSeparator className="hidden md:block" />}
        {hash !== '#/' && <BreadcrumbItem>
          <BreadcrumbPage>{hashMenu} Management</BreadcrumbPage>
        </BreadcrumbItem>}
      </BreadcrumbList>
    </Breadcrumb>
  )
}

function App() {
  const loginSend = async (valuez:{ identity: string; password: string }) => {
    toast.promise(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      LoginUser(valuez).catch(async (err: any) => {
        if (err.message === "CSRF validation failed") {
          const retryData = await LoginUser(valuez);
          return retryData;
        }
        throw err;
      }),
      {
        loading: "Waiting...",
        success: () => {
          window.location.href = "/";
          return "Success"
        },
        error: (err) => {
          return err.message || "Failed"
        },
      }
    )
  }

  const { auth, isLoading } = useAuth();

  const { data: checkMigration } = CheckMigration();

  const headerName = () => <div className="flex flex-col items-center gap-2 text-center">
    <h1 className="bg-linear-to-r from-blue-900 to-green-500 bg-clip-text text-4xl font-extrabold text-transparent">GO GERBANG</h1>
    <p className="text-muted-foreground text-sm text-balance">
      APPS Gateway Build by Muhammad Reza
    </p>
  </div>

  const loadingIndicator = <div className="fixed inset-0 flex items-center justify-center">
      <div className="relative w-10 h-10 border-2 border-black/70 border-b-transparent rounded-full animate-spin"></div>
  </div>

  if (isLoading) {
    return loadingIndicator
  }

  if (checkMigration?.data?.missing?.length && checkMigration?.data?.tables?.length) {
    return <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
        <AppBeginning  
          tables={checkMigration?.data?.tables}
          missing={checkMigration?.data?.missing}
        />
        <Toaster closeButton />
      </ThemeProvider>
  }

  if (!auth) {
    return <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
      <AppLogin 
        ImageLogo={LogoReact}
        ImageLogoWhite={LogoReact}
        ImageBanner={Placeholder}
        loginSend={loginSend}
        HeaderLogin={headerName}
      />
      <Toaster closeButton />
    </ThemeProvider>
  }

  return <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
    <AppShell 
      TeamSwitcher={TeamSwitcher}
      NavMain={NavMain}
      NavUser={NavUser}
      Breadcrumb={BreadcrumbCom}
      PageContent={AppRouter}
      variant="inset"
    />
    <Toaster closeButton />
  </ThemeProvider>
}

export default App
