<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">RESET</CardTitle>
      <CardDescription class="text-center">Get Link Reset Password</CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="submitLogin" class="grid gap-4">

        <div class="grid gap-2">
          <Label>Username/Email</Label>
          <Input v-model="form.identity" class="z-10" type="email" placeholder="Username/Email" required />
        </div>

        <div class="flex justify-center">
          <Button class="w-full" type="submit">Send</Button>
        </div>

        <div class="text-center text-sm">
          Sudah Punya Akun?
          <RouterLink to="/auth/login" class="underline underline-offset-4">
            Silahkan Login
          </RouterLink>
        </div>

      </form>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
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

interface FormLogin {
  identity: string
  password: string
}

const form: FormLogin = reactive({
  identity: '',
  password: ''
})

async function submitLogin() {
  // console.log(form)
  // try {
  //   const response = await fetch(url, {
  //     method: 'POST',
  //     headers: {
  //       'Content-Type': 'application/json',
  //     },
  //     body: JSON.stringify(form),
  //   });


  // } catch (error) {
    
  // }
}

</script>