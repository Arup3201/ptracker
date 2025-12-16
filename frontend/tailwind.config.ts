import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{ts,tsx}", "./index.html"],
  theme: {
    extend: {
      colors: {
        bg: {
          root: "var(--bg-root)",
          surface: "var(--bg-surface)",
          elevated: "var(--bg-elevated)",
        },
        border: {
          DEFAULT: "var(--border-default)",
          muted: "var(--border-muted)",
        },
        text: {
          primary: "var(--text-primary)",
          secondary: "var(--text-secondary)",
          muted: "var(--text-muted)",
        },
        primary: {
          DEFAULT: "var(--primary)",
          hover: "var(--primary-hover)",
        },
        success: "var(--success)",
        warning: "var(--warning)",
        danger: "var(--danger)",
      },
      borderRadius: {
        xs: "4px",
        sm: "6px",
        md: "8px",
      },
      boxShadow: {
        sm: "0 1px 2px rgba(0,0,0,0.4)",
        md: "0 4px 12px rgba(0,0,0,0.6)",
      },
      spacing: {
        3: "12px",
        5: "20px",
      },
      fontSize: {
        xs: ["11px", { lineHeight: "1.3" }],
        sm: ["12px", { lineHeight: "1.4" }],
        base: ["13px", { lineHeight: "1.5" }],
        lg: ["14px", { lineHeight: "1.4" }],
        xl: ["20px", { lineHeight: "1.3" }],
      },
    },
  },
  plugins: [],
};

export default config;
