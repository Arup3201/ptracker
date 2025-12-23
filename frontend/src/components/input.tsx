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
      className="h-8 w-full rounded-xs border border-(--border-default) bg-(--bg-surface) px-3 text-sm placeholder:text-(--text-muted) focus:border-(--primary) focus:outline-none"
    />
  );
}
