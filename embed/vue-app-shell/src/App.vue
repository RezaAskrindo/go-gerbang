<script setup lang="ts">
import { computed } from 'vue';
import { toast } from 'vue-sonner';
import { Toaster } from './components/ui/sonner';
import 'vue-sonner/style.css'

import { GetAuthSession, LoginUser } from './services/baseService';


import LogoVue from './assets/vue.svg'
import Placeholder from './assets/placeholder.svg'

import LoadingPage from './components/loading-page.vue';
import AppLogin from './components/app-login.vue';
import AppShell from './components/app-shell.vue';
import NavUser from './components/go-gerbang/NavUser.vue';
import TeamSwitcher from './components/go-gerbang/TeamSwitcher.vue';
import NavMain from './components/go-gerbang/NavMain.vue';

const { data, isLoading } = GetAuthSession();

const user = computed(() => ({
  name: data.value?.data.fullName,
  email: data.value?.data.username,
  avatar: "/avatars/shadcn.jpg",
}))

const loginSend = async (valuez:{ identity: string; password: string }) => {
  toast.promise(
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
      error: (err: any) => {
        return err.message || "Failed"
      },
    }
  )
}

</script>

<template>
  <LoadingPage v-if="isLoading" />
  <AppLogin v-else-if="!isLoading && !data?.status" 
    :-image-logo="LogoVue"
    :-image-logo-white="LogoVue"
    :-image-banner="Placeholder"
    :login-send="loginSend"
  />
  <AppShell v-else variant="inset">
    <template #sidebar-header>
      <TeamSwitcher />
    </template>
    <template #sidebar-content>
      <NavMain />
    </template>
    <template #sidebar-footer>
      <NavUser :user="user" />
    </template>
  </AppShell>
  <Toaster />
</template>