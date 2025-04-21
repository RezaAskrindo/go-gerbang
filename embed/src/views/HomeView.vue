<script setup lang="ts">
const listApp = [
  {
    name: "SISKOR",
    url: "https://app.siskor.web.id"
  }
]
</script>

<template>
  <div class="flex min-h-svh w-full items-center justify-center">
    <div class="w-full max-w-sm ">
      <div class="font-semibold text-xl text-center">Welcome to <a href="https://github.com/RezaAskrindo/go-gerbang" class="text-blue-700 dark:text-blue-500 underline">GO Gerbang</a></div>
      <div class="font-semibold text-xl text-center">Link Aplikasi</div>
      <ol class="text-center mt-2">
        <li v-for="list in listApp" :key="list.name">
          <a :href="list.url" class="text-blue-700 dark:text-blue-500 underline">{{ list.name }}</a>
        </li>
      </ol>
    </div>
  </div>
</template>
