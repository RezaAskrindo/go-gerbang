<template>
  <button @click="toggleDropdown" ref="reference" type="button" class="flex justify-center items-center rounded-full cursor-pointer w-10 h-10 hover:bg-gray-100">
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M12 13a3 3 0 1 0 0 -6a3 3 0 0 0 0 6z" /><path d="M12 3c7.2 0 9 1.8 9 9s-1.8 9 -9 9s-9 -1.8 -9 -9s1.8 -9 9 -9z" /><path d="M6 20.05v-.05a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v.05" /></svg>
  </button>
  <div v-show="isOpen" ref="floating" :style="floatingStyles" class="bg-white shadow-sm p-2">
    <div class="">Header</div>
    <ul>
      <li>Test</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { useFloating } from '@floating-ui/vue';

const isOpen = ref(false);

const reference = ref(null);
const floating = ref(null);
const { floatingStyles, update } = useFloating(reference, floating);

const toggleDropdown = () => {
  isOpen.value = !isOpen.value;
  handleResize()
};

const handleResize = () => isOpen.value && update();

onMounted(() => {
  window.addEventListener("resize", handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
});

</script>