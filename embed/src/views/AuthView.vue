<script setup lang="ts">
import { Toaster } from '@/components/sonner'
</script>

<template>
  <Toaster />
  <div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
    <div class="w-full max-w-sm">
      <div class="flex flex-col gap-6">
        <RouterView />
        <div class="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 hover:[&_a]:text-primary">
          Please Login or Sign Up If You Don't Have an Account
        </div>
      </div>
    </div>
  </div>
</template>