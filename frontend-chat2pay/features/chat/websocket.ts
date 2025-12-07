type MessageHandler = (data: string) => void;
type ConnectionHandler = () => void;

interface WebSocketConfig {
  url: string;
  userId: string;
  onMessage: MessageHandler;
  onConnect?: ConnectionHandler;
  onDisconnect?: ConnectionHandler;
  onError?: (error: Event) => void;
}

class ChatWebSocket {
  private ws: WebSocket | null = null;
  private config: WebSocketConfig | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private isIntentionallyClosed = false;

  connect(config: WebSocketConfig): void {
    this.config = config;
    this.isIntentionallyClosed = false;
    this.createConnection();
  }

  private createConnection(): void {
    if (!this.config) return;

    const wsUrl = `${this.config.url}/ws/chat/${this.config.userId}`;
    
    try {
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.reconnectAttempts = 0;
        this.config?.onConnect?.();
      };

      this.ws.onmessage = (event) => {
        this.config?.onMessage(event.data);
      };

      this.ws.onclose = () => {
        if (!this.isIntentionallyClosed && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.reconnectAttempts++;
          setTimeout(() => this.createConnection(), this.reconnectDelay * this.reconnectAttempts);
        }
        this.config?.onDisconnect?.();
      };

      this.ws.onerror = (error) => {
        this.config?.onError?.(error);
      };
    } catch (error) {
      console.error("WebSocket connection error:", error);
    }
  }

  send(message: string): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(message);
    } else {
      console.error("WebSocket is not connected");
    }
  }

  disconnect(): void {
    this.isIntentionallyClosed = true;
    this.ws?.close();
    this.ws = null;
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export const chatWebSocket = new ChatWebSocket();

export function getWebSocketBaseUrl(): string {
  return process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:9005";
}
