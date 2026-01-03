export function TopBar({
  title,
  actions,
}: {
  title: string;
  actions?: React.ReactNode;
}) {
  return (
    <div className="flex h-12 items-center justify-between border-b border-(--border-default) p-4">
      <div className="text-sm font-semibold">{title}</div>
      {actions}
    </div>
  );
}
