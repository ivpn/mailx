/// <reference types="vitest/config" />
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: [
      // Exact match only: src/lib/overlay.ts imports '@preline/overlay/index.js' itself.
      { find: /^@preline\/overlay$/, replacement: fileURLToPath(new URL('./src/lib/overlay.ts', import.meta.url)) },
    ],
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./vitest.setup.ts'],
  },
})
