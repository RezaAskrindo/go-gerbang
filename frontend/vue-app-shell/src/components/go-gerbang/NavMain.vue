<script setup lang="ts">
import { BookOpen, ChevronRight } from "lucide-vue-next"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"

const items = [
  {
    title: "All Menu",
    url: "#",
    icon: BookOpen,
    isActive: true,
    items: [
      {title: "user", url: "#/user"},
      // {title: "rbac", url: "#/rbac"},
    ]
  }
]
</script>

<template>
  <SidebarGroup>
      <SidebarGroupLabel>Go Gerbang Menu</SidebarGroupLabel>
      <SidebarMenu>
        <Collapsible
          v-for="item in items"
          :key="item.title"
          :default-open="item.isActive"
          class="group/collapsible"
          as-child
        >
          <SidebarMenuItem>
            <CollapsibleTrigger as-child>
              <SidebarMenuButton tooltip={item.title}>
                <component :is="item.icon" v-if="item.icon" />
                <span>{{ item.title }}</span>
                <ChevronRight class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
              </SidebarMenuButton>
            </CollapsibleTrigger>
            <CollapsibleContent>
              <SidebarMenuSub>
                <SidebarMenuSubItem v-for="subItem in item.items" :key="subItem.title">
                  <SidebarMenuSubButton as-child>
                    <a :href="subItem.url">
                      <span>{{ subItem.title }}</span>
                    </a>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>
      </SidebarMenu>
    </SidebarGroup>
</template>
