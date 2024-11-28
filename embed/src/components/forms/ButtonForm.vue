<style scoped>
.button-control {
  background-color: var(--color-heading);
  color: var(--color-background-mute);
  font-weight: 700;
  font-size: .875rem;
  line-height: 1.25rem;
  padding: .5rem 1rem;
  border-radius: calc(.5rem - 1px);
  gap: .5rem;
  justify-content: center;
  align-items: center;
  display: inline-flex;
  cursor: pointer;
}
</style>

<template>
  <button :type="typeButton" class="button-control">{{ label }}</button>
</template>

<script setup lang="ts">
interface Props {
  typeButton?: 'button' | 'submit' | 'reset'
  label?: string
  labels?: string[]
}

// const props = 
withDefaults(defineProps<Props>(), {
  typeButton: 'button',
  label: 'Test',
  labels: () => ['one', 'two']
})

const model = defineModel();

</script>