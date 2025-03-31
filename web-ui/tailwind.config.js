/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'background-light': 'rgb(249 251 253)',
        'background-dark': 'rgb(10 16 23)',
        'background-glass': 'rgba(20, 30, 40, 0.7)',
        'background': {
          700: 'rgb(15, 23, 33)',
          800: 'rgb(10, 16, 23)',
          900: 'rgb(5, 10, 15)'
        },
        'primary': {
          400: 'rgb(0, 180, 240)',
          500: 'rgb(0, 140, 220)',
          600: 'rgb(0, 120, 200)',
          700: 'rgb(0, 100, 180)',
          800: 'rgb(10, 37, 64)',
          900: 'rgb(5, 20, 35)',
        },
        'secondary': {
          500: 'rgb(0, 207, 180)'
        },
        'accent-blue': 'rgb(0, 112, 243)',
        'accent-purple': 'rgb(112, 0, 255)',
        'accent-teal': 'rgb(0, 207, 180)',
      },
      backdropBlur: {
        'glass': '10px',
      },
      boxShadow: {
        'glass': '0 4px 30px rgba(0, 0, 0, 0.1)',
        'glow': '0 0 15px rgba(0, 120, 255, 0.5)',
      },
      borderWidth: {
        '3': '3px',
      },
    },
  },
  plugins: [],
}
