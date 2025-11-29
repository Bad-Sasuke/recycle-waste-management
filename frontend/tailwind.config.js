import { createRequire } from 'module';
const require = createRequire(import.meta.url);

/** @type {import('tailwindcss').Config} */
import daisyui from 'daisyui'
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        'kanit': ['Kanit', 'sans-serif'],
      },
      colors: {
        // Map Tailwind colors to CSS variables
        primary: 'var(--color-primary)',
        'primary-content': 'var(--color-primary-content)',
        secondary: 'var(--color-secondary)',
        accent: 'var(--color-accent)',
        success: 'var(--color-success)',
        info: 'var(--color-info)',
        'info-dark': 'var(--color-info-dark)',
        warning: 'var(--color-warning)',
        'warning-dark': 'var(--color-warning-dark)',
        danger: 'var(--color-danger)',
        'danger-dark': 'var(--color-danger-dark)',
        disabled: 'var(--color-disabled)',

        // Backgrounds
        'body': 'var(--bg-body)',
        'surface': 'var(--bg-surface)',
        'footer': 'var(--bg-footer)',
        'navbar': 'var(--bg-navbar)',

        // Text
        'text-primary': 'var(--text-primary)',
        'text-secondary': 'var(--text-secondary)',
        'text-muted': 'var(--text-muted)',
        'text-footer': 'var(--text-footer)',
        'text-footer-hover': 'var(--text-footer-hover)',
        'text-inverse': 'var(--text-inverse)',

        // Borders
        'border-navbar': 'var(--border-navbar)',
      },
      backgroundImage: {
        'gradient-primary': 'var(--gradient-primary)',
        'gradient-info': 'var(--gradient-info)',
        'gradient-warning': 'var(--gradient-warning)',
        'gradient-danger': 'var(--gradient-danger)',
      },
      boxShadow: {
        'card': 'var(--shadow-card)',
        'blue': 'var(--shadow-blue)',
        'orange': 'var(--shadow-orange)',
        'red': 'var(--shadow-red)',
      },
      borderRadius: {
        'sm': 'var(--radius-sm)',
        'md': 'var(--radius-md)',
        'lg': 'var(--radius-lg)',
        'xl': 'var(--radius-xl)',
        'full': 'var(--radius-full)',
      },
    },
    screens: {
      'sm': '640px',
      // => @media (min-width: 640px) { ... }

      'md': '960px',

      'lg': '1440px',

      'xl': '1280px',
      // => @media (min-width: 1280px) { ... }

      '2xl': '1536px',
      // => @media (min-width: 1536px) { ... }
    },
  },
  plugins: [daisyui],
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["light"],
          "primary": "#22c55e", // green-500
          "primary-content": "#ffffff",
          "secondary": "#10b981", // emerald-500
          "accent": "#059669", // emerald-600
        },
      },
    ],
  }
}
