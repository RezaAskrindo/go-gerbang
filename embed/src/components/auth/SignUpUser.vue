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
        <div class="grid gap-2">
          <div class="flex items-center">
            <Label>Password</Label>
            <Button @click="toggleShowPassword" class="ml-auto text-sm !h-0" variant="link">
              <span v-if="showPassword">Show</span>
              <span v-else>Hide</span>
              Password
            </Button>
          </div>
          <Input v-model="form.password" :type="showPassword ? 'password' : 'text'" placeholder="Password" required />
        </div>
        <div class="grid gap-2">
          <Label>Ulangi Password</Label>
          <Input v-model="form.password_repeat" :type="showPassword ? 'password' : 'text'" placeholder="Password" required />
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
import { reactive, ref } from 'vue';
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

const showPassword = ref(true);
const toggleShowPassword = () => showPassword.value = !showPassword.value

const form: FormLogin = reactive({
  username: '',
  fullName: '',
  email: '',
  phoneNumber: '',
  password: '',
  password_repeat: ''
})

const isLoading = ref(false);

// const removeSpaces = () => {
//   form.username = form.username.replace(/[^a-zA-Z0-9]/g, '');
// };

async function submitLogin() {
  if (form.password !== form.password_repeat) {
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
    console.error('Error:', err);
  }
}

</script>