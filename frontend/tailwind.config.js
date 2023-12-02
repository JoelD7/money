/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
    colors:{
      white: {
        100: '#FFFFFF',
        200: '#e6e6e6',
        300: '#cccccc',
      },
      green: {
        100: '#009821',
        200: '#00851d',
        300: '#005212'
      },
      orange: {
        100: '#FF8042',
        200: '#ee6520',
        300: '#d44c08'
      },
      red: {
        100: '#D90707',
        200: '#ad0101',
        300: '#7a0101',
      },
      darkGreen: '#024511',
      blue: {
        100: '#0088FE',
        200: '#006dcc',
        300: '#004d99',
      },
      gray:{
        100: '#D9D9D9',
        200: '#6F6F6F',
        300: '#4D4D4D',
      }
    },
    fontFamily: {
      sans: ['Outfit', 'sans-serif'],
    }
  },
  plugins: [],
}

