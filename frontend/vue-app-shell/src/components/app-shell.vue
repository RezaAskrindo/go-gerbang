<script setup lang="ts">
import { useColorMode } from '@vueuse/core';
import { Moon, Sun } from 'lucide-vue-next';

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
  type SidebarProps,
} from "@/components/ui/sidebar"
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';

const props = defineProps<SidebarProps>();

const mode = useColorMode();
</script>

<template>
  <SidebarProvider>
    <Sidebar v-bind="props">
      <SidebarHeader>
        <slot name="sidebar-header"></slot>
      </SidebarHeader>
      <SidebarContent>
        <slot name="sidebar-content"></slot>
      </SidebarContent>
      <SidebarFooter>
        <slot name="sidebar-footer"></slot>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
    <SidebarInset>
      <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 border-b">
        <div class="flex items-center gap-2 px-4">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
          <slot name="header-breadcrumb"></slot>
        </div>
        <div className="ms-auto">
          <Button
            variant="outline"
            size="icon"
            @click="mode = mode === 'dark' ? 'light' : 'dark'"
            class="me-4"
          >
            <Sun v-if="mode === 'dark'" className="hidden [html.dark_&]:block" />
            <Moon v-else className="hidden [html.light_&]:block" />
            <span className="sr-only">Toggle theme</span>
          </Button>
        </div>
      </header>

      <div className="flex flex-1 flex-col">
        <div className="@container/main flex flex-1 flex-col gap-2">
          <slot />
        </div>
      </div>

    </SidebarInset>
  </SidebarProvider>
</template>
