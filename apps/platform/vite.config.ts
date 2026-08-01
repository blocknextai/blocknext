import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import removeConsole from 'vite-plugin-remove-console'
import path from 'path'

const srcPath = path.resolve(new URL('./src', import.meta.url).pathname)

export default defineConfig(() => ({
  envDir: '../../',
  base: '/',
  server: {
    port: 4000,
  },
  plugins: [react(), tailwindcss(), removeConsole()],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (
            id.includes('node_modules/react/') ||
            id.includes('node_modules/react-dom/') ||
            id.includes('node_modules/react-router/') ||
            id.includes('node_modules/react-is/') ||
            id.includes('node_modules/scheduler/') ||
            id.includes('node_modules/clsx/') ||
            id.includes('node_modules/tailwind-merge/') ||
            id.includes('node_modules/class-variance-authority/') ||
            id.includes('/src/lib/utils.ts')
          ) {
            return 'react-vendor'
          }

          if (id.includes('node_modules/@xyflow/')) {
            return 'xyflow'
          }
          if (id.includes('node_modules/d3-')) {
            return 'd3'
          }
          if (
            id.includes('node_modules/ethers/') ||
            id.includes('node_modules/@adraffy/') ||
            id.includes('node_modules/@noble/')
          ) {
            return 'web3'
          }
          if (id.includes('node_modules/@dicebear/')) {
            return 'dicebear'
          }
        },
      },
    },
  },
  resolve: {
    alias: {
      '@': srcPath,
    },
  },
}))
