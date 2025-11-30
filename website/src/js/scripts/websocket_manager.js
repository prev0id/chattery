import { DEFAULT_CONFIG } from "./consts.js";

class WebSocketManager {
  constructor(url, options = {}) {
    this.url = url;
    this.reconnectInterval =
      options.reconnectInterval || DEFAULT_CONFIG.RECONNECT_INTERVAL;
    this.reconnectAttempts = 0;
    this.subscriptions = new Map();
    this.isConnected = false;
    this.shouldReconnect = true;

    this.connect();
  }

  connect() {
    try {
      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        console.log("WebSocket connected");
        this.isConnected = true;
        this.reconnectAttempts = 0;
        this.onConnected?.();
      };

      this.ws.onmessage = (event) => {
        this.handleMessage(event.data);
      };

      this.ws.onclose = () => {
        console.log("WebSocket disconnected");
        this.isConnected = false;
        this.onDisconnected?.();
        this.handleReconnect();
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        this.onError?.(error);
      };
    } catch (error) {
      console.error("WebSocket connection failed:", error);
      this.handleReconnect();
    }
  }

  handleReconnect() {
    if (this.shouldReconnect) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect... (${this.reconnectAttempts})`);
      setTimeout(() => this.connect(), this.reconnectInterval);
    }
  }

  handleMessage(data) {
    try {
      const message = JSON.parse(data);

      const handlers = this.subscriptions.get(message.type) || [];
      handlers.forEach((handler) => {
        try {
          handler(message);
        } catch (error) {
          console.error("Error in message handler:", error);
        }
      });
    } catch (error) {
      console.error("Error parsing message:", error, data);
    }
  }

  sendMessage(message) {
    if (!this.isConnected) {
      console.warn("WebSocket not connected, message not sent:", message);
      return false;
    }

    try {
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error("Error sending message:", error);
      return false;
    }
  }

  subscribe(messageType, callback) {
    if (!this.subscriptions.has(messageType)) {
      this.subscriptions.set(messageType, []);
    }
    this.subscriptions.get(messageType).push(callback);

    return () => {
      const handlers = this.subscriptions.get(messageType);
      if (handlers) {
        const index = handlers.indexOf(callback);
        if (index > -1) {
          handlers.splice(index, 1);
        }
        if (handlers.length === 0) {
          this.subscriptions.delete(messageType);
        }
      }
    };
  }

  onConnected(callback) {
    this.onConnected = callback;
  }

  onDisconnected(callback) {
    this.onDisconnected = callback;
  }

  onError(callback) {
    this.onError = callback;
  }

  disconnect() {
    this.shouldReconnect = false;
    if (this.ws) {
      this.ws.close();
    }
  }
}

export default WebSocketManager;
