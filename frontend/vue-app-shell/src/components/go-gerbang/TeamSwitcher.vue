<script setup lang="ts">
import { ShieldPlus } from "lucide-vue-next";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

</script>

<template>
  <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton class="data-[slot=sidebar-menu-button]:p-1.5!" asChild>
          <a href="#">
            <ShieldPlus class="size-5!" />
            <span class="text-base font-semibold">GO Gerbang</span>
          </a>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
</template>
