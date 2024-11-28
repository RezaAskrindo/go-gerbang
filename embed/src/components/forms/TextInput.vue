<style scoped>
.form-control {
  padding: .5rem .75rem;
  border-radius: .5rem;
  width: 100%;
  line-height: 1.5rem;
  font-size: .875rem;
  border-width: 1px;
  display: block;
  background-color: var(--color-background);
  border-color: var(--color-heading);
  color: var(--color-heading);
}
</style>

<template>
  <input v-model="model" :type="typeInput" class="form-control" />
</template>

<script setup lang="ts">
interface Props {
  typeInput?: string
  labels?: string[]
}

// const props = 
withDefaults(defineProps<Props>(), {
  typeInput: 'text',
  labels: () => ['one', 'two']
})

const model = defineModel<any>();

</script>