<template>
  <div class="card">
    <form @submit.prevent="submitLogin" class="card-body">
      <h3 style="text-align: center; margin-bottom: 2rem;">Go Gerbang</h3>

      <label>Username</label>
      <TextInput v-model="form.identity" style="margin-bottom: 1rem; margin-top: .5rem;" required />
      <label>Password</label>
      <TextInput v-model="form.password" typeInput="password" style="margin-bottom: 1rem; margin-top: .5rem;" required />

      <div style="display: flex; justify-content: center;">
        <ButtonForm typeButton="submit" label="Login" style="min-width: 7rem;" />
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import ButtonForm from '@/components/forms/ButtonForm.vue';
import TextInput from '@/components/forms/TextInput.vue';

interface FormLogin {
  identity: string
  password: string
}

const form: FormLogin = reactive({
  identity: '',
  password: ''
})

function submitLogin() {
  console.log(form)
}

</script>