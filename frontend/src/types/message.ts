export interface WebSocketMessage {
  type: string;
  data: Record<string, any>;
}

export interface NotificationMessage {
  id: string;
  message: WebSocketMessage;
  timestamp: number;
}
