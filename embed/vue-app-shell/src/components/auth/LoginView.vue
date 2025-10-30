<script setup lang="ts">
import { useColorMode } from '@vueuse/core';
import type { Component } from 'vue';
import LoginForm from './LoginForm.vue';

interface LoginProps {
  ImageLogo?: string
  ImageLogoWhite?: string
  ImageBanner?: string
  loginSend?: (valuez:{ identity: string; password: string }) => void
  HeaderLogin?: Component
  FooterLogin?: Component
  ResetPasswordForm?: Component
}

const props = defineProps<LoginProps>();

const mode = useColorMode();
</script>

<template>
  <div class="grid min-h-svh lg:grid-cols-2">
    <div class="flex flex-col gap-4 p-6 md:p-10">
      <div class="flex justify-center gap-2 md:justify-start">
        <img
          v-if="props.ImageLogo && props.ImageLogoWhite"
          :src="mode === 'dark' && ImageLogoWhite ? ImageLogoWhite : ImageLogo"
          alt="Image"
          class="h-12"
        />
      </div>
      <div class="flex flex-1 items-center justify-center">
        <div class="w-full max-w-xs">
          <LoginForm 
            :loginSend="props.loginSend"
            :HeaderLogin="props.HeaderLogin"
            :FooterLogin="props.FooterLogin"
            :ResetPasswordForm="props.ResetPasswordForm"
          />
        </div>
      </div>
    </div>
    <div v-if="props.ImageBanner" class="bg-muted relative hidden lg:block">
      <img
        :src="props.ImageBanner"
        alt="Image"
        class="absolute inset-0 h-full w-full object-cover"
      />
    </div>
  </div>
</template>