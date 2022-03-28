const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: [
    "./src/pages/**/*.{tsx}",
    "./src/components/**/*.{tsx}",
    "./src/theme/**/*.{tsx}",
    "./docs/**/**/*.{md,mdx}",
  ],
  corePlugins: {
    preflight: false,
  },
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Source Sans"', defaultTheme.fontFamily.sans],
        mono: ['"Fira Code"', defaultTheme.fontFamily.mono],
      },
    },
  },
  plugins: [],
};
