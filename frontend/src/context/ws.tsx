import { createContext, useContext } from "react";
import { useWebsocket } from "../hooks/ws";
import toast from "react-hot-toast";
import {
  Ban,
  ClipboardList,
  Info,
  MessageSquareIcon,
  UserCheck,
  UserRoundPlus,
} from "lucide-react";

interface WebSocketContextValue {
  isConnected: boolean;
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
  const { isConnected, reconnect } = useWebsocket({
    url,
    onMessage: (message) => {
      // show notifications
      switch (message.type) {
        case "join_accepted":
          showCustomToast(
            `Your request to join "${message.data["project_name"]}" is accepted`,
            <UserCheck color="#00563b" size={18} />,
          );
          break;
        case "join_rejected":
          showCustomToast(
            `Your request to join "${message.data["project_name"]}" is rejected`,
            <Ban color="#e23d28" size={18} />,
          );
          break;
        case "task_updated":
          showCustomToast(
            `Task "${message.data["title"]}" updated`,
            <ClipboardList color="#ffbf00" size={18} />,
          );
          break;
        case "assignee_added":
          showCustomToast(
            `You have been assigned to task "${message.data["title"]}"`,
            <UserRoundPlus color="#e4d004" size={18} />,
          );
          break;
        case "assignee_removed":
          showCustomToast(
            `You have been removed as assignee from task "${message.data["title"]}"`,
            <UserRoundPlus color="#e23d28" size={18} />,
          );
          break;
        case "comment_added":
          showCustomToast(
            `New comment added to task "${message.data["title"]}"`,
            <MessageSquareIcon size={18} />,
          );
          break;
        default:
          showCustomToast(`New notification recieved`, <Info size={18} />);
      }
    },
    onOpen: () => {
      console.log("WebSocket connection established");
    },
    onClose: () => {
      console.log("WebSocket connection closed");
    },
  });

  return (
    <WebSocketContext.Provider
      value={{
        isConnected,
        reconnect,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
}

interface CustomToastProps {
  message: string;
  icon?: React.ReactNode;
  id: string;
}

const CustomToast: React.FC<CustomToastProps> = ({ message, icon, id }) => {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: "12px",
        padding: "12px 16px",
        background: "#1f2937",
        border: "1px solid #374151",
        borderRadius: "8px",
        boxShadow:
          "0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -1px rgba(0, 0, 0, 0.2)",
        minWidth: "300px",
        maxWidth: "400px",
      }}
    >
      {/* Icon */}
      {icon && <div style={{ flexShrink: 0, fontSize: "20px" }}>{icon}</div>}

      {/* Message */}
      <div
        style={{
          flex: 1,
          color: "#f9fafb",
          fontSize: "14px",
          lineHeight: "1.5",
        }}
      >
        {message}
      </div>

      {/* Close Button */}
      <button
        onClick={() => toast.dismiss(id)}
        style={{
          flexShrink: 0,
          background: "transparent",
          border: "none",
          color: "#9ca3af",
          cursor: "pointer",
          padding: "4px",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          borderRadius: "4px",
          transition: "all 0.2s",
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.background = "#374151";
          e.currentTarget.style.color = "#f9fafb";
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.background = "transparent";
          e.currentTarget.style.color = "#9ca3af";
        }}
        aria-label="Close notification"
      >
        <svg
          width="18"
          height="18"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18" />
          <line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
    </div>
  );
};

// Helper function to show custom toast
const showCustomToast = (message: string, icon?: React.ReactNode) => {
  toast.custom((t) => <CustomToast message={message} icon={icon} id={t.id} />, {
    duration: Infinity,
    position: "bottom-right",
  });
};
