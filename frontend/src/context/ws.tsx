import { createContext, useContext, useState } from "react";
import { useWebsocket } from "../hooks/ws";
import type { NotificationMessage } from "../types/message";

interface WebSocketContextValue {
  isConnected: boolean;
  notifications: NotificationMessage[];
  clearNotifications: () => void;
  reconnect: () => void;
}

const WebSocketContext = createContext<WebSocketContextValue | null>(null);

export function useWebSocketContext() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error(
      "useWebSocketContext must be used within WebSocketProvider",
    );
  }
  return context;
}

interface WebSocketProviderProps {
  children: React.ReactNode;
  url: string;
}

export function WebSocketProvider({ children, url }: WebSocketProviderProps) {
  const [notifications, setNotifications] = useState<NotificationMessage[]>([]);

  const { isConnected, reconnect } = useWebsocket({
    url,
    onMessage: (message) => {
      const id = Date.now() + "" + Math.random();
      // Add notification with timestamp
      setNotifications((prev) => [
        ...prev,
        { id: id, message: message, timestamp: new Date().getTime() },
      ]);
    },
    onOpen: () => {
      console.log("WebSocket connection established");
    },
    onClose: () => {
      console.log("WebSocket connection closed");
    },
  });

  const clearNotifications = () => {
    setNotifications([]);
  };

  return (
    <WebSocketContext.Provider
      value={{
        isConnected,
        notifications,
        clearNotifications,
        reconnect,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
}
