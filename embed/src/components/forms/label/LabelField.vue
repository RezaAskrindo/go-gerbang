<script setup lang="ts">
import { cn } from '@/lib/utils'
import { type HTMLAttributes } from 'vue'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()

</script>

<template>
  <label :class="cn('text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70',props.class)">
    <slot />
  </label>
</template>