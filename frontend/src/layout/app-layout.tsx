import { Outlet } from "react-router";

export function AppLayout() {
  return (
    <div className="flex h-screen bg-(--bg-root) text-(--text-primary)">
      <Outlet />
    </div>
  );
}
