import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/backend': {
        target: 'http://localhost:9000',
        // target: 'https://auth.siskor.web.id/backend',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/backend/, ''),
      },
    }
  },
  preview: {
    proxy: {
      '/backend': {
        target: 'http://localhost:9000',
        // target: 'https://apigateway.siskor.web.id',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/backend/, ''),
      },
    }
  }
})
