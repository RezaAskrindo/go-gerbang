import { 
  lazy,
  useEffect, 
  useMemo, 
  useState, 
  type ReactNode 
} from "react";

import {
  BookOpen, 
  ChevronRight, 
  ChevronsUpDown, 
  LogOut, 
  ShieldPlus 
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
  GetAuthSession,
  LoginUser,
  LogoutUser,
} from "./services/baseService";

import AppLogin from "@/components/app-login";

const AppShell = lazy(() => import("@/components/app-shell"));
const Dashboard = lazy(() => import("@/components/go-gerbang/dashboard"));
const UserManagement = lazy(() => import("@/components/go-gerbang/user-management"));
const RbacManagement = lazy(() => import("@/components/go-gerbang/rbac-management"));
const ModuleuManagement = lazy(() => import("@/components/go-gerbang/module-management"));

const TeamSwitcher = () => {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton className="data-[slot=sidebar-menu-button]:!p-1.5" asChild>
          <a href="#">
            <ShieldPlus className="!size-5" />
            <span className="text-base font-semibold">GO Gerbang</span>
          </a>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}

const NavMain = () => {
  const items = [
    {
      title: "All Menu",
      url: "#",
      icon: BookOpen,
      isActive: true,
      items: [
        {title: "user", url: "#/user"},
        {title: "module", url: "#/module"},
        {title: "rbac", url: "#/rbac"},
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
                      <SidebarMenuSubButton asChild>
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
  "/user": UserManagement,
  "/module": ModuleuManagement,
  "/rbac": RbacManagement,
};

function AppRouter(): ReactNode {
  const [currentPath, setCurrentPath] = useState(window.location.hash);

  useEffect(() => {
    const onHashChange = () => setCurrentPath(window.location.hash);
    window.addEventListener("hashchange", onHashChange);

    return () => window.removeEventListener("hashchange", onHashChange);
  }, []);

  const CurrentView = useMemo(() => {
    return routes[currentPath.slice(1) || "/"] || NotFound;
  }, [currentPath]);

  return <CurrentView />;
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

  const headerName = () => <div className="flex flex-col items-center gap-2 text-center">
    <h1 className="bg-linear-to-r from-blue-900 to-orange-500 bg-clip-text text-4xl font-extrabold text-transparent">GO GERBANG</h1>
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
      PageContent={AppRouter}
      variant="inset"
    />
    <Toaster closeButton />
  </ThemeProvider>
}

export default App
