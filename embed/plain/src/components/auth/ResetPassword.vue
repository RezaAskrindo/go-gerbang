<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">RESET</CardTitle>
      <CardDescription class="text-center">Get Link Reset Password</CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="submitRequestResetPassword" class="grid gap-4">

        <div class="grid gap-2">
          <Label>Username/Email</Label>
          <Input v-model="form.identity" class="z-10" type="text" placeholder="Username/Email" required />
        </div>

        <div class="flex justify-center">
          <Button class="w-full" type="submit">Send</Button>
        </div>

        <div class="text-center text-sm">
          Already had account?
          <RouterLink :to="`/auth/login${pathQuery}`" class="underline underline-offset-4">
            Login Here
          </RouterLink>
        </div>

      </form>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { pathQuery } from '@/stores/app.store';
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
import { sendNotification } from '@/lib/notification';

interface ResetPassword {
  identity: string
}

const form: ResetPassword = reactive({
  identity: '',
})

async function submitRequestResetPassword() {
  try {
    const getCsrf = await getCSRFToken();

    const response = await fetch(`${baseHost}/api/v1/auth/request-reset-password${pathQuery.value}&baseUrl=https://auth.siskor.web.id/auth`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'X-SGCsrf-Token': getCsrf
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
    toast.error('Error occurred while getting reset password. Please try again later.');
    sendNotification(`Error Get Reset Password Go Gerbang!\n details: ${JSON.stringify(form.identity)}`);
    console.error('Error:', err);
  }
}

</script>