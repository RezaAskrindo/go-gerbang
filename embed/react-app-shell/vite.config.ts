import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from "@tailwindcss/vite"
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    cssInjectedByJsPlugin(),
    federation({
      name: 'app-shell',
      filename: 'remoteEntry.js',
      exposes: {
        './AppShell': './src/components/app-shell.tsx',
        './Sidebar': './src/components/ui/sidebar.tsx',
        './UseHelper': './src/components/ui/useHelper.tsx',
        './ThemeProvider': './src/components/theme-provider.tsx',
      },
      shared: ['react']
    })
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 3000
  },
  // base: 'http://localhost:2001',
  preview: {
    port: 2001,
  },
})
