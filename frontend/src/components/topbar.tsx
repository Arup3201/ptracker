export function TopBar({
  title,
  actions,
}: {
  title: string;
  actions?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col md:flex-row md:items-center md:justify-between border-b border-border px-4 md:px-6 py-3 md:h-14 shrink-0 gap-3 md:gap-0">
      <h1 className="text-base md:text-lg font-semibold text-text-primary">
        {title}
      </h1>

      {actions && (
        <div className="flex items-center justify-end gap-2">{actions}</div>
      )}
    </div>
  );
}
