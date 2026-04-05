export function TopBar({
  title,
  actions,
}: {
  title: string;
  actions?: React.ReactNode;
}) {
  return (
    <div className="flex h-14 items-center justify-between border-b border-border px-6 shrink-0">
      <h1 className="text-lg font-semibold text-text-primary tracking-snug">
        {title}
      </h1>
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  );
}
