interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary";
}

export function Button({
  variant = "primary",
  className = "",
  ...props
}: ButtonProps) {
  const base =
    "h-8 px-3 text-sm font-medium rounded-xs transition-colors focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed";

  const variants = {
    primary: "bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white",
    secondary:
      "bg-transparent border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-elevated)]",
  };

  return (
    <button
      className={`${base} ${variants[variant]} ${className}`}
      {...props}
    />
  );
}
