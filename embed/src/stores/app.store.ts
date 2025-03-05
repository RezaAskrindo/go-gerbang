// import { ref, computed } from 'vue'
// import { defineStore } from 'pinia'
import { useFetch } from '@vueuse/core'

export const baseURL =  "http://localhost:9000";


export const { data: getCSRFToken } = useFetch(`${baseURL}/secure-gateway-c`, { credentials: 'include' }).get().json();

// export const useAppStore = defineStore('app-store', () => {
//   const count = ref(0)
//   const doubleCount = computed(() => count.value * 2)
//   function increment() {
//     count.value++
//   }

//   return { count, doubleCount, increment }
// })
