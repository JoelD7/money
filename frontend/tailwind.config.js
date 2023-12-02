/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
    colors:{
      green: '#009821',
      orange: '#FF8042',
      red: '#D90707',
      darkGreen: '#024511',
      lightBlue: '#0088FE',
      gray:{
        100: '#D9D9D9',
        200: '#6F6F6F',
      }
    },
    fontFamily: {
      sans: ['Outfit', 'sans-serif'],
    }
  },
  plugins: [],
}

