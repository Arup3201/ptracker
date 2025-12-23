import { Outlet } from "react-router";
import { Sidebar } from "../components/sidebar";

export function AppLayout() {
  return (
    <div className="flex h-screen overflow-y-hidden bg-(--bg-root) text-(--text-primary)">
      <Sidebar />
      <Outlet />
    </div>
  );
}
