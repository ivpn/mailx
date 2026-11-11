/// <reference types="vitest/config" />
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // The @preline/* 2.x packages are CommonJS bundles that set `exports.default`. Vite 8 resolves a
  // default import of CommonJS the way Node does for a "type": "module" package - to the whole
  // `module.exports` - so `import dropdown from '@preline/dropdown'` stopped yielding the class
  // (dropdown.autoInit is not a function). This restores the Vite <= 7 interop those imports rely on.
  legacy: {
    inconsistentCjsInterop: true,
  },
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
