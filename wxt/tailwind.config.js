/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./entrypoints/**/*.{js,ts,jsx,tsx,vue,html}", // include WXT entrypoints
    "./components/**/*.{js,ts,jsx,tsx,vue,html}", // if you have shared components
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
