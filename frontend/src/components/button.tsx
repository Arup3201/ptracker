interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "danger";
}

export function Button({
  variant = "primary",
  className = "",
  ...props
}: ButtonProps) {
  const base =
    "h-8 px-3 text-sm font-medium rounded-xs transition-colors focus:outline-none cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed";

  const variants = {
    primary: "bg-(--primary) hover:bg-(--primary-hover) text-white",
    secondary:
      "bg-transparent border border-(--border-default) text-(--text-primary) hover:bg-(--bg-elevated)",
    danger: "text-(--danger) hover:bg-(--danger)/10",
  };

  return (
    <button
      className={`${base} ${variants[variant]} ${className}`}
      {...props}
    />
  );
}
