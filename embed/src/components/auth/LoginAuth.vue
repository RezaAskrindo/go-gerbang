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
          <Input v-model="form.identity" class="z-10" type="text" placeholder="Username/Email" required :aria-invalid="isIdentityError" />
        </div>
        <div class="grid gap-2">
          <Label>Password</Label>
          <Input v-model="form.password" placeholder="Password" required :use-password-show="true" :aria-invalid="isPasswordError" />
        </div>

        <div class="flex justify-center">
          <Button v-if="isLoading" class="w-full z-10" type="button" disabled>Loading...</Button>
          <Button v-else class="w-full z-10" type="submit">Login</Button>
        </div>

        <div class="text-center text-sm z-10">
          Forget Password?
          <RouterLink :to="`/auth/reset-password${pathQuery}`" class="underline underline-offset-4">
            Get Link Password
          </RouterLink>
        </div>

        <div class="text-center text-sm z-10">
          Don't have an account?
          <RouterLink :to="`/auth/sign-up${pathQuery}`" class="underline underline-offset-4">
            Sign Up
          </RouterLink>
        </div>
        <!-- <div class="text-center text-sm z-10">
          Your account is not active yet?
          <RouterLink to="sign-up" class="underline underline-offset-4">
            Request to active
          </RouterLink>
        </div> -->

        <!-- <div class="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-t after:border-border">
          <span class="relative z-10 bg-background px-2 text-muted-foreground">Or continue with</span>
        </div> -->

        <div class="flex flex-col gap-4">
          <!-- <GoogleLogin :callback="callbackGoogleAccount" prompt>
            <Button variant="outline" class="w-full">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path
                  d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
                  fill="currentColor"
                />
              </svg>
              Login with Google
            </Button>
          </GoogleLogin> -->
        </div>

      </form>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
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
import { pathQuery } from '@/stores/app.store';
import { sendNotification } from '@/lib/notification';

// import { GoogleLogin } from "vue3-google-login"

interface FormLogin {
  identity: string
  password: string
}

// interface FormGoogleLogin {
//   id_token: string
//   client_id: string
// }

const form: FormLogin = reactive({
  identity: '',
  password: ''
})

const route = useRoute();

const isLoading = ref(false);

const isIdentityError = ref(false);
const isPasswordError = ref(false);

const LoginExecution = async (formLogin: FormLogin, urlLogin: string) => {
  try {
    isIdentityError.value = false;
    isPasswordError.value = false;
    isLoading.value = true;
    const getCsrf = await getCSRFToken();

    const url = `${baseHost}/api/v1/auth/${urlLogin}${pathQuery.value}`;
    const response = await fetch(url, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'X-SGCsrf-Token': getCsrf
      },
      body: JSON.stringify(formLogin),
    });
    
    let result
    
    if (!response.ok) {
      result = await response.json();
      toast.error(result?.message);
      if (result?.message === 'record not found') {
        isIdentityError.value = true;
      } else if (result?.message === 'wrong password') {
        isPasswordError.value = true;
      }
    } else {
      result = await response.json();
      console.log(result);
      toast.error(result?.message);
      const getUrl = route.query?.url;
      if (getUrl) {
        window.location.href = getUrl.toString();
      }
    }
    isLoading.value = false;
  } catch (error) {
    isLoading.value = false;
    console.log(error)
    toast.error('Error occurred while logging in. Please try again later.');
    sendNotification(`Error Login Auth Go Gerbang!\n details: ${JSON.stringify(formLogin.identity)}`);
    return error
  }
}
 
// const callbackGoogleAccount = async (response: any) => {
//   if (response.credential) {
//     const formLogin: FormGoogleLogin = { 
//       id_token: response.credential, 
//       client_id: '114782695264-qbhlmm64mf883aetb4l07tf4m7jv4ek1.apps.googleusercontent.com'
//     }
//     await LoginExecution(formLogin, 'login-with-google');
//   } else if (response?.code) {
//     const getResponse = await fetch("https://oauth2.googleapis.com/token", {
//       method: "POST",
//       headers: { "Content-Type": "application/x-www-form-urlencoded" },
//       body: new URLSearchParams({
//         code: response.code,
//         client_id: '114782695264-qbhlmm64mf883aetb4l07tf4m7jv4ek1.apps.googleusercontent.com',
//         client_secret: "eDNrWTVP7l5uQNYlDeuj1hCx",
//         redirect_uri: "https://auth.siskor.web.id/callback", // must match your app
//         grant_type: "authorization_code"
//       })
//     });

//     const tokens = await getResponse.json();
//     const formLogin: FormGoogleLogin = { 
//       id_token: tokens, 
//       client_id: '114782695264-qbhlmm64mf883aetb4l07tf4m7jv4ek1.apps.googleusercontent.com'
//     }
//     await LoginExecution(formLogin, 'login-with-google');
//   } else {
//     console.log(response);
//     toast.error("failed to get JWT Token");
//   }
// }

async function submitLogin() {
  await LoginExecution(form, 'login');
}

</script>