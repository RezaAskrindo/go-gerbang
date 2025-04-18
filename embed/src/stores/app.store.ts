import { computed } from 'vue';
import { useRoute } from 'vue-router';

export const pathQuery = computed(() => {
  const route = useRoute();
  const query = route.query;
  const queryString = Object.keys(query)
    .map(key => `${key}=${query[key]}`)
    .join('&');
  return `${queryString ? `?${queryString}` : ''}`;
})