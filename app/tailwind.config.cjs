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
        ne: {
          1: '#0f0f0f',
          2: '#191919',
          3: '#222222',
          4: '#292929',
          5: '#313131',
          6: '#3a3a3a',
          7: '#484848',
          8: '#606060',
          9: '#6e6e6e',
          10: '#7b7b7b',
          11: '#b4b4b4',
          12: '#eeeeee',
        },
        bl: {
          1: '#09111d',
          2: '#0e1827',
          3: '#0b264c',
          4: '#032f68',
          5: '#093b7b',
          6: '#15488d',
          7: '#1f57a4',
          8: '#2567c2',
          9: '#007aff',
          10: '#006cf0',
          11: '#7ab6ff',
          12: '#cce3ff',
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
  darkMode: 'selector',
}
