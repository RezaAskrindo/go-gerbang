<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-center">SIGN UP</CardTitle>
      <CardDescription class="text-center">Please Fill Form</CardDescription>
    </CardHeader>
    <CardContent>
      <form @submit.prevent="submitLogin" class="grid gap-4">

        <!-- <div class="grid gap-2">
          <Label>Username</Label>
          <Input v-model="form.username" @input="removeSpaces" type="text" placeholder="Username" required />
        </div> -->
        <div class="grid gap-2">
          <Label>Full Name</Label>
          <Input v-model="form.fullName" type="text" placeholder="Full Name" required />
        </div>
        <div class="grid gap-2">
          <Label>Email</Label>
          <Input v-model="form.email" type="email" placeholder="Email" required />
        </div>
        <div class="grid gap-2">
          <Label>Phone Number</Label>
          <Input v-model="form.phoneNumber" type="text" placeholder="Phone Number" required />
        </div>
        <p></p>
        <div class="grid gap-2">
          <Label>Password</Label>
          <Input v-model="form.password" placeholder="Password" required :use-password-show="true" :aria-invalid="wrongAtPassword" />
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2.5 mb-2 dark:bg-gray-700">
          <div v-if="passwordStrength.level === 3" class="bg-green-600 h-2.5 rounded-full dark:bg-green-500" :style="{ width: passwordStrength.percentage+'%' }"></div>
          <div v-else-if="passwordStrength.level === 2" class="bg-teal-600 h-2.5 rounded-full dark:bg-teal-500" :style="{ width: passwordStrength.percentage+'%' }"></div>
          <div v-else-if="form.password.length > 4" class="bg-orange-600 h-2.5 rounded-full dark:bg-orange-500" style="width: 50%"></div>
          <div v-else class="bg-red-600 h-2.5 rounded-full dark:bg-red-500" :style="{ width: passwordStrength.percentage+'%' }"></div>
        </div>
        <div class="grid gap-2">
          <Label>Ulangi Password</Label>
          <Input v-model="form.password_repeat" placeholder="Ulangi Password" required :use-password-show="true" :aria-invalid="wrongAtPassword" />
        </div>

        <div style="display: flex; justify-content: center;">
          <Button v-if="isLoading" class="w-full" type="button" disabled>Sign Up</Button>
          <Button v-else class="w-full" type="submit">Sign Up</Button>
        </div>

        <div class="text-center text-sm">
          Alread have an account?
          <RouterLink :to="`/auth/login${pathQuery}`" class="underline underline-offset-4">
            Sign In
          </RouterLink>
        </div>

      </form>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { pathQuery } from '@/stores/app.store';
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
import { sendNotification } from '@/lib/notification';
import { CheckPasswordStrong } from '@/lib/checkPasswordStrong';

interface FormLogin {
  username: string
  fullName: string
  email: string
  phoneNumber: string
  password: string
  password_repeat: string
}

const route = useRoute();
const router = useRouter();

const form: FormLogin = reactive({
  username: '',
  fullName: '',
  email: '',
  phoneNumber: '08',
  password: '',
  password_repeat: ''
})

const isLoading = ref(false);

const wrongAtPassword = ref(false);

const passwordStrength = computed(() => CheckPasswordStrong(form.password));

// const removeSpaces = () => {
//   form.username = form.username.replace(/[^a-zA-Z0-9]/g, '');
// };

async function submitLogin() {
  if (form.password !== form.password_repeat) {
    wrongAtPassword.value = true;
    toast.error("You're Password Not Equal");
    return '';
  }

  try {
    isLoading.value = true;
    form.username = form.email;
    const getCsrf = await getCSRFToken();

    const sender = route?.query?.sender ? `&sender=${route.query.sender}` : '';

    const response = await fetch(`${baseHost}/api/v1/auth/sign-up?notif=true&active=true${sender}`, {
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
      if (result.status) {
        router.push(`/auth/login${pathQuery.value}`);
      }
    }
    isLoading.value = false;
  } catch(err) {
    isLoading.value = false;
    toast.error('Error occurred while Signing up. Please try again later.');
    sendNotification(`Error Sign Up Go Gerbang!\n details: ${JSON.stringify(form)}`);
    console.error('Error:', err);
  }
}

</script>