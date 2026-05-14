<script setup lang="ts">
import { ref, type Component } from 'vue';
import { useColorMode } from '@vueuse/core';
import { ChevronRight, Lock, LockOpen } from 'lucide-vue-next';

import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { Spinner } from "@/components/ui/spinner"

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

const showPassword = ref(false);
const loadingButton = ref(false);

const form = ref({
  identity: "",
  password: "",
})

const onFormSubmit = async() => {
  if (!props.loginSend) return true;

  loadingButton.value = true;

  props.loginSend(form.value);

  setTimeout(() => {
    loadingButton.value = false;
  }, 1500)
}

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
        <div v-else>LOGO</div>
      </div>
      <div class="flex flex-1 items-center justify-center">
        <div class="w-full max-w-xs">

          <form @submit.prevent="onFormSubmit" class="flex flex-col gap-6">
    
            <component v-if="props.HeaderLogin" :is="props.HeaderLogin" />

            <div class="grid gap-6">
              <div class="grid w-full max-w-sm items-center gap-1.5">
                <Label for="identity">Identity</Label>
                <Input v-model="form.identity" id="identity" type="text" placeholder="name@example.com" />
              </div>
              <div class="grid w-full max-w-sm items-center gap-1.5">
                <Label for="password">Password</Label>
                <div class="relative">
                  <Input v-model="form.password" id="password" :type="showPassword ? 'text' : 'password'" placeholder="name@example.com" />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    @click="showPassword = !showPassword"
                  >
                    <LockOpen v-if="showPassword" class="size-4" aria-hidden="true" />
                    <Lock v-else class="size-4" aria-hidden="true" />
                  </Button>
                </div>
              </div>

              <Button :type="loadingButton ? 'button' : 'submit'" variant="outline" :disabled="loadingButton">
                <Spinner v-if="loadingButton" />
                <ChevronRight v-else />
                {{ loadingButton ? 'Loading...' : 'Login' }}
              </Button>

              <component v-if="props.FooterLogin" :is="props.FooterLogin" />

            </div>
          </form>

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