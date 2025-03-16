<template>
  <Card class="relative line-border">
    <CardHeader>
      <CardTitle class="text-center">SIGN IN</CardTitle>
      <CardDescription class="text-center">Single Sign In Page</CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="submitLogin" class="grid gap-4">

        <div class="grid gap-2">
          <Label>Username/Email</Label>
          <Input v-model="form.identity" class="z-10" type="text" placeholder="Username/Email" required />
        </div>
        <div class="grid gap-2">
          <div class="flex items-center">
            <Label>Password</Label>
            <RouterLink to="/auth/reset-password" class="ml-auto text-sm underline-offset-4 hover:underline z-10">
              Forgot your password?
            </RouterLink>
          </div>
          <Input v-model="form.password" class="z-10" type="password" placeholder="Password" required />
        </div>

        <div class="flex justify-center">
          <Button class="w-full z-10" type="submit">Login</Button>
        </div>

        <div class="text-center text-sm z-10">
          Don't have an account?
          <RouterLink to="sign-up" class="underline underline-offset-4">
            Sign Up
          </RouterLink>
        </div>
        <div class="text-center text-sm z-10">
          Your account is not active yet?
          <RouterLink to="sign-up" class="underline underline-offset-4">
            Request to active
          </RouterLink>
        </div>

        <div class="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
          <span class="relative z-10 bg-background px-2 text-muted-foreground">Or continue with</span>
        </div>

        <div class="flex flex-col gap-4">
          <Button variant="outline" class="w-full">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path
                d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
                fill="currentColor"
              />
            </svg>
            Login with Google
          </Button>
        </div>

      </form>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { getCSRFToken, baseHost } from '@/stores/worker.service';

import { 
  Card, 
  CardDescription, 
  CardHeader, 
  CardTitle,
  CardContent
} from '@/components/card/index';
import { Input } from '@/components/forms/input';
import { Label } from '@/components/forms/label';
import { Button } from '@/components/forms/button';

import { toast } from 'vue-sonner'

interface FormLogin {
  identity: string
  password: string
}

const form: FormLogin = reactive({
  identity: 'reza',
  password: '9192'
})

async function submitLogin() {
  try {
    const getCsrf = await getCSRFToken();

    const url = `${baseHost}/api/v1/auth/login?session=true&block=true&httponly=true&domain=localhost`;
    const response = await fetch(url, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'X-SGCsrf-Token': getCsrf?.data
      },
      body: JSON.stringify(form),
    });
    
    let result
    
    if (!response.ok) {
      result = await response.json();
      toast.error(result?.message);
    } else {
      result = await response.json();
      toast.error(result?.message);
    }
  } catch (error) {
    console.log(error)
    return error
  }
}

</script>