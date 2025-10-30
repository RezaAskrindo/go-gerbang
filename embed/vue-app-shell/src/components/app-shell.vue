<script setup lang="ts">
import { defineAsyncComponent, type Component, type VNodeChild } from 'vue';
import { useColorMode } from '@vueuse/core';
import { Moon, Sun } from 'lucide-vue-next';

import AppSidebar from "./app-sidebar.vue"
import {
  SidebarInset,
  SidebarProvider,
} from "@/components/ui/sidebar"
import { Button } from '@/components/ui/button';
import LoadingPage from './loading-page.vue';

const LoginView = defineAsyncComponent(() => import('./auth/LoginView.vue'))

interface AppShellProps {
  Breadcrumb?: Component
  TeamSwitcher?: Component
  NavMain?: Component
  CardInformation?: Component
  NavUser?: Component
  PageContent?: () => VNodeChild
  AuthChecking?: boolean
  AuthPassed?: boolean
  ImageLogo?: string
  ImageLogoWhite?: string
  ImageBanner?: string
  loginSend?: (valuez: { identity: string; password: string }) => void
  HeaderLogin?: Component
  FooterLogin?: Component
  ResetPasswordForm?: Component
}

const props = defineProps<AppShellProps>();
// console.log(props.AuthPassed)

const mode = useColorMode();
</script>

<template>
  <LoadingPage v-if="props.AuthChecking" />

  <Suspense v-else-if="typeof props.AuthPassed !== undefined && props.AuthPassed === false">
    <template #default>
      <LoginView 
        :ImageLogo="props.ImageLogo"
        :ImageLogoWhite="props.ImageLogoWhite"
        :ImageBanner="props.ImageBanner"
        :loginSend="props.loginSend"
        :HeaderLogin="props.HeaderLogin"
        :FooterLogin="props.FooterLogin"
        :ResetPasswordForm="props.ResetPasswordForm"
      />
    </template>
    <template #fallback>
      <div class="flex items-center justify-center h-screen">
        <div class="w-10 h-10 border-2 border-black border-b-transparent rounded-full animate-spin"></div>
      </div>
    </template>
  </Suspense>

  <SidebarProvider v-else>
    <AppSidebar />
    <SidebarInset>
      <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 border-b">
        <div class="flex items-center gap-2 px-4">
          
        </div>
        <div className="ms-auto">
          <Button
            variant="outline"
            size="icon"
            @click="mode = mode === 'dark' ? 'light' : 'dark'"
          >
            <Sun v-if="mode === 'dark'" className="hidden [html.dark_&]:block" />
            <Moon v-else className="hidden [html.light_&]:block" />
            <span className="sr-only">Toggle theme</span>
          </Button>
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4">
        <div class="grid auto-rows-min gap-4 md:grid-cols-3">
          <div class="aspect-video rounded-xl bg-muted/50" />
          <div class="aspect-video rounded-xl bg-muted/50" />
          <div class="aspect-video rounded-xl bg-muted/50" />
        </div>
        <div class="min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min" />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
