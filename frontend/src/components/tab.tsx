export function Tab({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={
        "pb-2 text-sm font-medium transition-colors " +
        (active
          ? "border-b-2 border-(--primary) text-(--text-primary)"
          : "text-(--text-muted) hover:text-(--text-primary) cursor-pointer")
      }
    >
      {label}
    </button>
  );
}
