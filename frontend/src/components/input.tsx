import { useState, type ChangeEvent, type InputHTMLAttributes } from "react";
import { cn } from "../utils/cn";

interface InputProps extends Omit<
  InputHTMLAttributes<HTMLInputElement>,
  "onChange"
> {
  /** Native change handler — receives the full event (compatible with React Hook Form, etc.) */
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  /** Legacy string-only handler kept for backwards compatibility */
  onValueChange?: (value: string) => void;
}

export function Input({
  placeholder = "",
  onChange,
  onValueChange,
  value: controlledValue,
  className,
  ...props
}: InputProps) {
  // Uncontrolled internal state used only when no `value` prop is provided
  const [internalValue, setInternalValue] = useState("");

  const isControlled = controlledValue !== undefined;
  const value = isControlled ? controlledValue : internalValue;

  function handleChange(e: ChangeEvent<HTMLInputElement>) {
    if (!isControlled) {
      setInternalValue(e.target.value);
    }
    onChange?.(e);
    onValueChange?.(e.target.value);
  }

  return (
    <input
      value={value}
      onChange={handleChange}
      placeholder={placeholder}
      className={cn(
        "h-9 w-full rounded-md border border-border bg-bg-elevated px-3 text-sm text-text-primary",
        "placeholder:text-text-muted",
        "transition duration-fast",
        "focus:border-primary focus:outline-none focus:shadow-focus-primary",
        "disabled:opacity-40 disabled:cursor-not-allowed",
        className,
      )}
      {...props}
    />
  );
}
