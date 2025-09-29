import { type ComponentProps, type ComponentType } from "react"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

type AppSidebarProps = ComponentProps<typeof Sidebar> & {
  TeamSwitcher?: ComponentType
  NavMain?: ComponentType
  CardInformation?: ComponentType
  NavUser?: ComponentType
};

export function AppSidebar({ 
  TeamSwitcher,
  NavMain,
  CardInformation,
  NavUser,
  ...props 
}: AppSidebarProps) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        {TeamSwitcher ? <TeamSwitcher /> : null}
      </SidebarHeader>
      <SidebarContent>
        {NavMain ? <NavMain /> : null}
      </SidebarContent>
      <SidebarFooter>
        {CardInformation ? <CardInformation /> : null}
        {NavUser ? <NavUser /> : null}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}