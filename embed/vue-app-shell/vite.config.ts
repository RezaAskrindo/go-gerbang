import path from 'node:path'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js'
import federation from '@originjs/vite-plugin-federation'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(), 
    tailwindcss(),
    cssInjectedByJsPlugin(),
    federation({
      name: 'app-shell',
      filename: 'remoteEntry.js',
      exposes: {
        './AppShell': './src/components/app-shell.vue',
        './AppLogin': './src/components/app-login.vue',
        './Sidebar': './src/components/ui/sidebar/index.ts',
      },
      shared: ['vue']
    })
  ],
  server: {
    port: 3000
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
