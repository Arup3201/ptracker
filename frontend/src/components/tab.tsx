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
        "pb-2.5 text-sm font-medium transition duration-fast cursor-pointer focus:outline-none " +
        (active
          ? "border-b-2 border-primary text-text-primary"
          : "border-b-2 border-transparent text-text-muted hover:text-text-secondary")
      }
    >
      {label}
    </button>
  );
}
