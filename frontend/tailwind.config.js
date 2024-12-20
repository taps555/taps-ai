/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}", // Scan all JavaScript/TypeScript files in src/
  ],

  theme: {
    extend: {},
  },
  plugins: [require("tailwind-hamburgers")],
};
