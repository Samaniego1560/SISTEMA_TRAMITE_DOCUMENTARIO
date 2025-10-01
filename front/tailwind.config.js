/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,ts}",
  ],
  theme: {
    extend: {
      colors :{
        'brand': '#0F75BC',
        'brand-2': '#1ABB9C'
      }
    },
  },
  plugins: [],
}
