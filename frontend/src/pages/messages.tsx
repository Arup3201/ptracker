import { useEffect, useState } from "react";
import { TopBar } from "../components/topbar.tsx";
import { ApiFetch } from "../utils/api";
import { renderLocalTime } from "../utils.ts";
import { Button } from "../components/button.tsx";

interface Message {
  id: string;
  user_id: string;
  type: string;
  body: {
    project?: {
      id: string;
      name: string;
    };
    task?: {
      id: string;
      title: string;
    };
    updates?: Array<{
      field: string;
      to: string;
    }>;
    updater?: {
      username: string;
      display_name: string;
    };
    assignee?: {
      username: string;
      display_name: string;
    };
    requestor?: {
      username: string;
      display_name: string;
    };
    responder?: {
      username: string;
      display_name: string;
    };
    commenter?: {
      username: string;
      display_name: string;
    };
    status?: string;
  };
  read: boolean;
  created_at?: string;
}

interface MessagesApiResponse {
  messages: Message[];
}

const getNotificationDisplayName = (name: string | undefined): string => {
  return name || "Unknown";
};

const NotificationContent: React.FC<{ notification: Message }> = ({
  notification,
}) => {
  const { type, body } = notification;
  const projectName = body.project?.name || "Unknown Project";
  const taskTitle = body.task?.title || "Unknown Task";

  switch (type) {
    case "task_added":
      return (
        <div>
          <p className="font-medium text-text-primary">New Task Added</p>
          <p className="text-sm text-text-muted">
            A new task <span className="font-semibold">{taskTitle}</span> was
            added to <span className="font-semibold">{projectName}</span>
          </p>
        </div>
      );

    case "task_updated":
      return (
        <div>
          <p className="font-medium text-text-primary">Task Updated</p>
          <p className="text-sm text-text-muted">
            <span className="font-semibold">
              {getNotificationDisplayName(body.updater?.display_name)}
            </span>{" "}
            updated task <span className="font-semibold">{taskTitle}</span>
            {body.updates && body.updates.length > 0 && (
              <span> ({body.updates.map((u) => u.field).join(", ")})</span>
            )}
          </p>
        </div>
      );

    case "assignee_added":
      return (
        <div>
          <p className="font-medium text-text-primary">Assigned to Task</p>
          <p className="text-sm text-text-muted">
            <span className="font-semibold">
              {getNotificationDisplayName(body.assignee?.display_name)}
            </span>{" "}
            was assigned to task{" "}
            <span className="font-semibold">{taskTitle}</span> in{" "}
            <span className="font-semibold">{projectName}</span>
          </p>
        </div>
      );

    case "assignee_removed":
      return (
        <div>
          <p className="font-medium text-text-primary">Removed from Task</p>
          <p className="text-sm text-text-muted">
            <span className="font-semibold">
              {getNotificationDisplayName(body.assignee?.display_name)}
            </span>{" "}
            was removed from task{" "}
            <span className="font-semibold">{taskTitle}</span>
          </p>
        </div>
      );

    case "join_requested":
      return (
        <div>
          <p className="font-medium text-text-primary">Join Request</p>
          <p className="text-sm text-text-muted">
            <span className="font-semibold">
              {getNotificationDisplayName(body.requestor?.display_name)}
            </span>{" "}
            requested to join{" "}
            <span className="font-semibold">{projectName}</span>
          </p>
        </div>
      );

    case "join_responded":
      return (
        <div>
          <p className="font-medium text-text-primary">Join Request Response</p>
          <p className="text-sm text-text-muted">
            Your request to join{" "}
            <span className="font-semibold">{projectName}</span> was{" "}
            <span className="font-semibold">{body.status}</span> by{" "}
            <span className="font-semibold">
              {getNotificationDisplayName(body.responder?.display_name)}
            </span>
          </p>
        </div>
      );

    case "comment_added":
      return (
        <div>
          <p className="font-medium text-text-primary">New Comment</p>
          <p className="text-sm text-text-muted">
            <span className="font-semibold">
              {getNotificationDisplayName(body.commenter?.display_name)}
            </span>{" "}
            commented on task <span className="font-semibold">{taskTitle}</span>{" "}
            in <span className="font-semibold">{projectName}</span>
          </p>
        </div>
      );

    default:
      return (
        <div>
          <p className="font-medium text-text-primary">Notification</p>
          <p className="text-sm text-text-muted">You have a new notification</p>
        </div>
      );
  }
};

export function MessagesPage() {
  const [notifications, setNotifications] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  const getMessages = async () => {
    try {
      const response = await ApiFetch("/messages");
      if (response.ok) {
        const responseData = await response.json();
        const data: MessagesApiResponse = responseData.data;
        if (data?.messages) {
          setNotifications(data.messages);
        }
      } else {
        throw new Error("Failed to get messages.");
      }
    } catch (err) {
      console.log(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getMessages();
  }, []);

  const sortedNotifications = [...notifications].sort((a, b) => {
    if (!a.created_at || !b.created_at) return 0;
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
  });

  const handleMarkAsRead = async (
    notificationId: string,
    currentReadStatus: boolean,
  ) => {
    try {
      const response = await ApiFetch(`/messages/${notificationId}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ read: !currentReadStatus }),
      });

      if (response.ok) {
        setNotifications((prevNotifications) =>
          prevNotifications.map((n) =>
            n.id === notificationId ? { ...n, read: !n.read } : n,
          ),
        );
      }
    } catch (err) {
      console.log("Failed to update notification read status:", err);
    }
  };

  return (
    <>
      <TopBar title="Messages" />

      <div className="flex-1 overflow-y-auto p-6">
        <div className="max-w-4xl mx-auto">
          {loading ? (
            <div className="text-center text-text-muted py-8">
              Loading messages...
            </div>
          ) : notifications.length === 0 ? (
            <div className="text-center text-text-muted py-8">
              No messages yet
            </div>
          ) : (
            <div className="space-y-3">
              {sortedNotifications.map((notification) => (
                <div
                  key={notification.id}
                  className={`rounded-lg p-3 md:p-4 hover:shadow-sm transition-shadow duration-fast relative ${
                    notification.read
                      ? "bg-bg-root border border-border text-text-muted"
                      : "bg-bg-surface border border-primary-light text-text-primary"
                  }`}
                >
                  {!notification.read && (
                    <Button
                      onClick={() =>
                        handleMarkAsRead(notification.id, notification.read)
                      }
                      className="absolute top-3 right-3 p-1"
                      title="Mark as read"
                    >
                      <svg
                        className="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                    </Button>
                  )}
                  <div className="flex flex-col gap-3">
                    <div className="flex-1 min-w-0">
                      <NotificationContent notification={notification} />
                    </div>
                    {notification.created_at && (
                      <div className="text-xs text-text-muted">
                        {renderLocalTime(notification.created_at)}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </>
  );
}
