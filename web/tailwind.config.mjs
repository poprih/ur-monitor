/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#921841",
          50: "#fdf2f6",
          100: "#fce7ef",
          200: "#f9c6d7",
          300: "#f594b3",
          400: "#ed5786",
          500: "#df2759",
          600: "#921841",
          700: "#7d1038",
          800: "#690e30",
          900: "#5a0e2a",
        },
        secondary: {
          DEFAULT: "#2c446d",
          50: "#f4f7fb",
          100: "#e9eff5",
          200: "#cfdae8",
          300: "#a4bcda",
          400: "#6f97c2",
          500: "#4678aa",
          600: "#2c446d",
          700: "#24375a",
          800: "#1e2f4d",
          900: "#1c2a42",
        },
      },
      fontFamily: {
        sans: ["Inter", "sans-serif"],
      },
    },
  },
  plugins: [],
};
