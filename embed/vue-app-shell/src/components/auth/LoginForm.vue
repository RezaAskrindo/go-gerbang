<script setup lang="ts">
import { ref, type Component } from 'vue';
import { Lock, LockOpen } from 'lucide-vue-next';

import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"

interface LoginFormProps {
  loginSend?: (valuez:{ identity: string; password: string }) => void
  HeaderLogin?: Component
  FooterLogin?: Component
  ResetPasswordForm?: Component
}

const props = defineProps<LoginFormProps>();

const showPassword = ref(false);
const loadingButton = ref(false);

const form = ref({
  identity: "",
  password: "",
})

const onFormSubmit = async() => {
  loadingButton.value = true;
  console.log(form.value);

  setTimeout(() => {
    loadingButton.value = false;
  }, 1500)
}

</script>

<template>
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

      <Button :type="loadingButton ? 'button' : 'submit'" variant="outline" className="w-full">
        {{ loadingButton ? 'Loading...' : 'Login' }}
      </Button>

      <component v-if="props.FooterLogin" :is="props.FooterLogin" />

    </div>
  </form>
</template>