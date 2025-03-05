<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">SIGN UP</CardTitle>
      <CardDescription class="text-center">Please Fill Form</CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="submitLogin" class="grid gap-4">

        <div class="grid gap-2">
          <Label>Username</Label>
          <Input v-model="form.username" type="text" placeholder="Username" required />
        </div>
        <div class="grid gap-2">
          <Label>Full Name</Label>
          <Input v-model="form.full_name" type="text" placeholder="Full Name" required />
        </div>
        <div class="grid gap-2">
          <Label>Email</Label>
          <Input v-model="form.email" type="email" placeholder="Email" required />
        </div>
        <div class="grid gap-2">
          <Label>Phone Number</Label>
          <Input v-model="form.phone_number" type="text" placeholder="Phone Number" required />
        </div>
        <div class="grid gap-2">
          <Label>Password</Label>
          <Input v-model="form.password" type="password" placeholder="Password" required />
        </div>

        <div style="display: flex; justify-content: center;">
          <Button class="w-full"  type="submit">Sign Up</Button>
        </div>

        <div class="text-center text-sm">
          Alread have an account?
          <RouterLink to="/auth/login" class="underline underline-offset-4">
            Sign In
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
  username: string
  full_name: string
  email: string
  phone_number: string
  password: string
}

const form: FormLogin = reactive({
  username: '',
  full_name: '',
  email: '',
  phone_number: '',
  password: '',
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