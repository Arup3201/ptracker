import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{ts,tsx}", "./index.html"],
  darkMode: "class",
  theme: {
    extend: {
      fontFamily: {
        sans: ["Geist", "sans-serif"],
        mono: ["Geist Mono", "monospace"],
      },

      colors: {
        bg: {
          root: "var(--bg-root)",
          surface: "var(--bg-surface)",
          elevated: "var(--bg-elevated)",
          overlay: "var(--bg-overlay)",
        },
        border: {
          DEFAULT: "var(--border-default)",
          muted: "var(--border-muted)",
          strong: "var(--border-strong)",
        },
        text: {
          primary: "var(--text-primary)",
          secondary: "var(--text-secondary)",
          muted: "var(--text-muted)",
          disabled: "var(--text-disabled)",
        },
        primary: {
          DEFAULT: "var(--primary)",
          hover: "var(--primary-hover)",
          muted: "var(--primary-muted)",
          foreground: "var(--primary-foreground)",
        },
        success: {
          DEFAULT: "var(--success)",
          muted: "var(--success-muted)",
        },
        warning: {
          DEFAULT: "var(--warning)",
          muted: "var(--warning-muted)",
        },
        danger: {
          DEFAULT: "var(--danger)",
          muted: "var(--danger-muted)",
        },
      },

      borderRadius: {
        xs: "4px",
        sm: "6px",
        md: "8px",
        lg: "10px",
        xl: "14px",
        "2xl": "18px",
        full: "9999px",
      },

      boxShadow: {
        xs: "0 1px 2px rgba(0,0,0,0.5)",
        sm: "0 1px 4px rgba(0,0,0,0.5), 0 0 0 0.5px rgba(255,255,255,0.04)",
        md: "0 4px 12px rgba(0,0,0,0.6), 0 0 0 0.5px rgba(255,255,255,0.04)",
        lg: "0 8px 24px rgba(0,0,0,0.7), 0 0 0 0.5px rgba(255,255,255,0.04)",
        inset: "inset 0 1px 0 rgba(255,255,255,0.06)",
        "focus-primary": "0 0 0 2px var(--primary-muted)",
        "focus-danger": "0 0 0 2px var(--danger-muted)",
      },

      fontSize: {
        "2xs": ["10px", { lineHeight: "1.3", letterSpacing: "0.02em" }],
        xs: ["11px", { lineHeight: "1.4", letterSpacing: "0.01em" }],
        sm: ["12px", { lineHeight: "1.5" }],
        base: ["13px", { lineHeight: "1.6" }],
        lg: ["14px", { lineHeight: "1.5" }],
        xl: ["16px", { lineHeight: "1.4" }],
        "2xl": ["20px", { lineHeight: "1.3" }],
        "3xl": ["24px", { lineHeight: "1.2", letterSpacing: "-0.02em" }],
      },

      fontWeight: {
        normal: "400",
        medium: "500",
        semibold: "600",
      },

      letterSpacing: {
        tight: "-0.02em",
        snug: "-0.01em",
        normal: "0em",
        wide: "0.02em",
        wider: "0.04em",
      },

      spacing: {
        0.5: "2px",
        1: "4px",
        1.5: "6px",
        2: "8px",
        2.5: "10px",
        3: "12px",
        3.5: "14px",
        4: "16px",
        5: "20px",
        6: "24px",
        7: "28px",
        8: "32px",
        9: "36px",
        10: "40px",
        11: "44px",
        12: "48px",
        14: "56px",
        16: "64px",
        20: "80px",
        24: "96px",
      },

      size: {
        4: "16px",
        5: "20px",
        6: "24px",
        7: "28px",
        8: "32px",
        9: "36px",
        10: "40px",
      },

      transitionDuration: {
        fast: "100ms",
        DEFAULT: "150ms",
        slow: "250ms",
      },

      transitionTimingFunction: {
        DEFAULT: "cubic-bezier(0.16, 1, 0.3, 1)",
        linear: "linear",
        "in-out": "cubic-bezier(0.4, 0, 0.2, 1)",
      },

      animation: {
        "fade-in": "fadeIn 150ms cubic-bezier(0.16, 1, 0.3, 1)",
        "slide-up": "slideUp 150ms cubic-bezier(0.16, 1, 0.3, 1)",
        "slide-down": "slideDown 150ms cubic-bezier(0.16, 1, 0.3, 1)",
        "scale-in": "scaleIn 150ms cubic-bezier(0.16, 1, 0.3, 1)",
      },

      keyframes: {
        fadeIn: {
          from: { opacity: "0" },
          to: { opacity: "1" },
        },
        slideUp: {
          from: { opacity: "0", transform: "translateY(4px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
        slideDown: {
          from: { opacity: "0", transform: "translateY(-4px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
        scaleIn: {
          from: { opacity: "0", transform: "scale(0.96)" },
          to: { opacity: "1", transform: "scale(1)" },
        },
      },

      screens: {
        xs: "480px",
        sm: "640px",
        md: "768px",
        lg: "1024px",
        xl: "1280px",
        "2xl": "1440px",
      },
    },
  },
  plugins: [],
};

export default config;
