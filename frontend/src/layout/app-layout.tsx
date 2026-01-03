import { Outlet } from "react-router";
import { Sidebar } from "../components/sidebar";

export function AppLayout() {
  return (
    <div className="flex h-screen overflow-y-hidden bg-(--bg-root) text-(--text-primary)">
      <Sidebar />
      <main className="flex flex-1 flex-col">
        <Outlet />
      </main>
    </div>
  );
}
