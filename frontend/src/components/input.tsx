import { useState, type ChangeEvent } from "react";

export function Input({
  placeholder = "",
  onChange = () => {},
}: {
  placeholder?: string;
  onChange?: (value: string) => void;
}) {
  const [value, setValue] = useState("");

  function handleChange(e: ChangeEvent<HTMLInputElement>) {
    const text = e.target.value;
    setValue(text);
    onChange(text);
  }

  return (
    <input
      value={value}
      onChange={(e) => handleChange(e)}
      placeholder={placeholder}
      className="h-9 w-full rounded-md border border-border bg-bg-elevated px-3 text-sm text-text-primary placeholder:text-text-muted transition duration-fast focus:border-primary focus:outline-none focus:shadow-focus-primary"
    />
  );
}
