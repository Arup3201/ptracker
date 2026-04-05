import { Outlet } from "react-router";
import { Sidebar } from "../components/sidebar";
import { WebSocketProvider } from "../context/ws";
import { WEBSOCKET_URL } from "../api/ws";

export function AppLayout() {
  return (
    <WebSocketProvider url={WEBSOCKET_URL}>
      <div className="flex h-screen overflow-hidden bg-bg-root text-text-primary">
        <Sidebar />
        <main className="flex flex-1 flex-col overflow-hidden pb-16 md:pb-0">
          <Outlet />
        </main>
      </div>
    </WebSocketProvider>
  );
}
