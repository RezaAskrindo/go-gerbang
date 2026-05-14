<script setup lang="ts">
import { useRoute } from 'vue-router'
import LogoSvg from '@/components/LogoSvg.vue';
import { Toaster } from '@/components/sonner'

const route = useRoute();
</script>

<template>
  <Toaster />
  <div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
    <div class="w-full max-w-sm">
      <div class="flex flex-col gap-6">
        <div v-if="route.query?.sender" class="flex items-center gap-2 self-center font-medium">
          <div class="flex items-center justify-center rounded-md bg-primary text-primary-foreground">
            <LogoSvg :logo-name="route.query?.sender.toString()" />
          </div>
          {{ route.query?.sender }}
        </div>
        <RouterView />
        <div class="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 hover:[&_a]:text-primary">
          Please 
          Login or 
          Sign Up 
          If You Don't Have an Account
        </div>
      </div>
    </div>
  </div>
</template>