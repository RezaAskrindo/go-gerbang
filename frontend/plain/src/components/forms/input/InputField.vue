<script setup lang="ts">
import { ref, type HTMLAttributes, watch } from 'vue'
import { cn } from '@/lib/utils'
import { useVModel } from '@vueuse/core'
import { Button } from '@/components/forms/button'
import MultipleIcon from '@/components/icons/MultipleIcon.vue'

const props = defineProps<{
  defaultValue?: string | number
  modelValue?: string | number
  class?: HTMLAttributes['class']
  placeholder?: string
  usePasswordShow?: boolean
  ariaInvalid?: boolean
  required?: boolean
}>()

const emits = defineEmits<{
  (e: 'update:modelValue', payload: string | number): void
}>()

// useVModel to bind the model value
const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue,
})

const showPassword = ref(false)

// Watch for changes in modelValue to reset aria-invalid if needed
watch(() => props.ariaInvalid, (newVal) => {
  // Handle aria-invalid based on your form validation logic
  if (newVal) {
    // Apply logic when invalid
  }
})
</script>

<template>
  <div v-if="usePasswordShow" class="relative w-full items-center">
    <input
      v-model="modelValue"
      :class="cn('flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 z-10 aria-invalid:ring-red-500/20 dark:aria-invalid:ring-red-500/40 aria-invalid:border-red-500', props.class)"
      :type="showPassword ? 'text' : 'password'"
      :placeholder="placeholder"
      :aria-invalid="props.ariaInvalid ? 'true' : 'false'"
      :required="props.required"
    />
    <Button
      type="button"
      variant="ghost"
      size="sm"
      class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent z-10"
      @click="showPassword = !showPassword"
    >
      <MultipleIcon :iconName="showPassword ? 'eye' : 'eye-closed'" class="h-4 w-4" />
      <span class="sr-only">{{ showPassword ? 'Hide password' : 'Show password' }}</span>
    </Button>
  </div>
  
  <input
    v-else
    v-model="modelValue"
    :class="cn('flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 z-10 aria-invalid:ring-red-500/20 dark:aria-invalid:ring-red-500/40 aria-invalid:border-red-500', props.class)"
    :placeholder="placeholder"
    :aria-invalid="props.ariaInvalid ? 'true' : 'false'"
    :required="props.required"
  />
</template>
