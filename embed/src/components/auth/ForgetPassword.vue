<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">New Password {{ isExpired ? 'Is Expired' : '' }}</CardTitle>
      <CardDescription class="text-center">Fill New Password</CardDescription>
    </CardHeader>
    <CardContent>
      <form v-if="!isExpired" @submit.prevent="submitNewPassword" class="grid gap-4">

        <div class="grid gap-2">
          <div class="flex items-center">
            <Label>New Password</Label>
            <Button @click="toggleShowPassword" class="ml-auto text-sm !h-0" variant="link">
              <span v-if="showPassword">Show</span>
              <span v-else>Hide</span>
              Password
            </Button>
          </div>
          <Input v-model="form.password" class="z-10" :type="showPassword ? 'password' : 'text'" placeholder="New Password" required />
        </div>

        <div class="grid gap-2">
          <Label>Repeat Password</Label>
          <Input v-model="form.passwordConfirm" class="z-10" :type="showPassword ? 'password' : 'text'" placeholder="Repeat Password" required />
        </div>

        <div class="flex justify-center">
          <Button class="w-full" type="submit">Update Password</Button>
        </div>

        <div class="text-center text-sm">
          Already had account?
          <RouterLink :to="`/auth/login${pathQuery}`" class="underline underline-offset-4">
            Login Here
          </RouterLink>
        </div>

      </form>
      <Button v-else @click="backToLogin" class="w-full">
        Ajukan Ulang Reset Password
      </Button>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
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
import { pathQuery } from '@/stores/app.store';
import { sendNotification } from '@/lib/notification';

interface ResetPassword {
  password: string
  passwordConfirm: string
}

const route = useRoute();
const router = useRouter();

const isExpired = ref(false);

const showPassword = ref(true);
const toggleShowPassword = () => showPassword.value = !showPassword.value

const form: ResetPassword = reactive({
  password: '',
  passwordConfirm: '',
})

async function submitNewPassword() {
  try {
    const getCsrf = await getCSRFToken();

    const response = await fetch(`${baseHost}/api/v1/auth/reset-password?token=${route.query?.token}`, {
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
      window.history.back();
    }
  } catch(err) {
    toast.error('Error occurred while filling new password. Please try again later.');
    sendNotification(`Error Fill New Password Go Gerbang!\n details: on token ${route.query?.token}`);
    console.error('Error:', err);
  }
}

const backToLogin = () => {
  router.push(`/auth/reset-password${pathQuery.value}`)
}

onMounted(() => {
  const token = route.query?.token;
  if (token) {
    const parts = token.toString().split("_");
    const expire = parseInt(parts[1]);
    const now = Date.now();
    const diffMs = now - expire;

    console.log(diffMs)
    const diffHours = diffMs / (1000 * 60 * 60);
    
    if (diffHours > 0) {
      isExpired.value = true;
    }
  }
})

</script>