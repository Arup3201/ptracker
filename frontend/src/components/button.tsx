import { cn } from "../utils/cn";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "danger";
  /**
   * "md" (default) — h-8, used throughout the app
   * "lg"           — h-9, better proportion next to Input fields (auth pages)
   */
  size?: "md" | "lg";
}

export function Button({
  variant = "primary",
  size = "md",
  className,
  ...props
}: ButtonProps) {
  const base =
    "inline-flex items-center justify-center px-3 text-xs font-medium rounded-lg transition-all duration-150 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-1 focus-visible:ring-offset-zinc-950 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed select-none";

  const sizes = {
    md: "h-8",
    lg: "h-9",
  };

  const variants = {
    primary:
      "bg-emerald-500 hover:bg-emerald-400 active:bg-emerald-600 text-zinc-950 font-semibold shadow-sm",
    secondary:
      "bg-transparent border border-zinc-700 text-zinc-300 hover:bg-zinc-800 hover:text-zinc-100",
    danger:
      "bg-transparent text-red-400 hover:bg-red-950/50 hover:text-red-300",
  };

  return (
    <button
      className={cn(base, sizes[size], variants[variant], className)}
      {...props}
    />
  );
}
