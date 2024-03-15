/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.templ'],
  theme: {
    extend: {
      boxShadow: {
        'm': '2px 2px 0px 1px rgba(0 0 0)',
        'xl': '5px 5px 0px 1px rgba(0 0 0)',
      },
      fontFamily: {
        Inria: "'Inria Sans'",
        Lilita: "'Lilita One'",
      },
      colors: {
        primary: '#a8f5ff',
        secondary: '#d5a4ed',
        tri: '#ffb2e6'
      },
    },
  },
  plugins: [],
}
