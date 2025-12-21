export function TopBar() {
  return (
    <div className="flex h-12 items-center justify-between border-b border-(--border-default) px-4">
      <div className="text-sm font-semibold">Dashboard</div>

      <button className="h-8 rounded-xs bg-(--primary) px-3 text-sm font-medium hover:bg-(--primary-hover)">
        New Project
      </button>
    </div>
  );
}
