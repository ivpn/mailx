/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
    "./node_modules/preline/preline.js",
  ],
  theme: {
    extend: {
      colors: {
        bluish: {
          50: '#eff8ff',
          100: '#dbeefe',
          200: '#bee2ff',
          300: '#92d1fe',
          400: '#5fb6fb',
          500: '#449cf8',
          600: '#2378ed',
          700: '#1b62da',
          800: '#1d4fb0',
          900: '#1d468b',
          950: '#162c55',
        },
      }
    },
  },
  plugins: [
    require('preline/plugin'),
    require("@tailwindcss/forms")({
      strategy: 'class',
    }),
  ],
}
