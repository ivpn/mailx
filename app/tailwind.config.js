/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
    "./node_modules/preline/preline.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('preline/plugin'),
    require("@tailwindcss/forms")({
      // strategy: 'base',
      // strategy: 'class',
    }),
  ],
}
