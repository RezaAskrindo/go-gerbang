<script setup lang="ts">
import UserDropdown from '@/components/topbar/UserDropdown.vue';
import type { AppTobarProps } from '@/types/LayoutType';
import { UseImage } from '@vueuse/components';

const { logoName = 'LOGO APP' } = defineProps<AppTobarProps>()

const classButton = "flex justify-center items-center rounded-full cursor-pointer w-10 h-10 hover:bg-gray-100"

</script>

<template>
  <div class="flex items-center fixed top-0 left-0 w-full h-16 px-4 z-[997] shadow-sm bg-white">
    <div class="min-w-80 flex flex-row items-center">
      <button type="button" :class="classButton">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 6l16 0" /><path d="M4 12l16 0" /><path d="M4 18l16 0" /></svg>
      </button>
      <a href="/" class="flex w-auto h-10 items-center">
        <UseImage v-if="logoImage" :src="logoImage">
          <template #loading>
            Loading..
          </template>
          <template #error>
            Failed
          </template>
        </UseImage>
        <span v-else>
          {{ logoName }}
        </span>
      </a>
    </div>

    <div class="flex gap-x-1 ms-auto">
      <div class="flex">
        
        <UserDropdown />
      </div>
    </div>
  </div>  
</template>