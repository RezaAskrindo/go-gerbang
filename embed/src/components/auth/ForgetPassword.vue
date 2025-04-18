<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">New Password</CardTitle>
      <CardDescription class="text-center">Fill New Password</CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="submitLogin" class="grid gap-4">

        <div class="grid gap-2">
          <Label>New Password</Label>
          <Input v-model="form.password" class="z-10" type="password" placeholder="New Password" required />
        </div>

        <div class="grid gap-2">
          <Label>Repeat Password</Label>
          <Input v-model="form.passwordConfirm" class="z-10" type="password" placeholder="Repeat Password" required />
        </div>

        <div class="flex justify-center">
          <Button class="w-full" type="submit">Update Password</Button>
        </div>

        <div class="text-center text-sm">
          Already had account?
          <RouterLink to="/auth/login" class="underline underline-offset-4">
            Login Here
          </RouterLink>
        </div>

      </form>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { useRoute } from 'vue-router';
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

interface ResetPassword {
  password: string
  passwordConfirm: string
}

const route = useRoute();

const form: ResetPassword = reactive({
  password: '',
  passwordConfirm: '',
})

async function submitLogin() {
  try {
    const getCsrf = await getCSRFToken();

    const response = await fetch(`${baseHost}/api/v1/auth/reset-password?token=${route.query?.token}`, {
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
  } catch(err) {
    console.error('Error:', err);
  }
}

</script>