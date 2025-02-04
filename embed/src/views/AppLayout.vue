<script setup lang="ts">
import AppTopbar from './AppTopbar.vue';

import type { AppLayoutProps } from '@/types/LayoutType';


defineProps<AppLayoutProps>()

</script>

<template>
  <main class="layout-wrapper mt-16">
    <AppTopbar :logo-image="logoImage" />
    <!-- <AppSidebar /> -->
    <div class="layout-main-container">
      <div class="layout-main">
        <slot></slot>
      </div>
    </div>
  </main>
</template>