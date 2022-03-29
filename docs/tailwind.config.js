const defaultTheme = require('tailwindcss/defaultTheme');

module.exports = {
  content: [
    './src/pages/**/*.{js,jsx,tsx}',
    './src/components/**/*.{js,jsx,tsx}',
    './src/theme/**/*.{js,jsx,tsx}',
    './docs/**/**/*.{md,mdx}',
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
      screens: {
        lg: '997px',
      },
      colors: {
        primary: {
          DEFAULT: 'var(--docs-color-primary)',
          100: 'var(--docs-color-primary-100)',
        },
        text: {
          DEFAULT: 'var(--docs-color-text)',
          100: 'var(--docs-color-text-100)',
        },
        border: 'var(--docs-color-border)',
        background: {
          DEFAULT: 'var(--docs-color-background)',
          100: 'var(--docs-color-background-100)',
          200: 'var(--docs-color-background-200)',
          300: 'var(--docs-color-background-300)',
        },
      },
    },
  },
  plugins: [],
};
