import { createSignal, onCleanup } from "solid-js";

export function createWebSocketClient(getUrl, options = {}) {
  const [status, setStatus] = createSignal("idle"); // idle | connecting | open | closed | error
  const [error, setError] = createSignal(null);
  const [lastMessage, setLastMessage] = createSignal(null);

  let socket = null;
  let reconnectTimer = null;
  let closedManually = false;
  const listeners = new Set();

  function resolveUrl() {
    return typeof getUrl === "function" ? getUrl() : getUrl;
  }

  function clearReconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  }

  function sendRaw(data) {
    if (!socket || socket.readyState !== WebSocket.OPEN) return false;
    socket.send(data);
    return true;
  }

  function sendJson(payload) {
    return sendRaw(JSON.stringify(payload));
  }

  function subscribe(handler) {
    listeners.add(handler);
    return () => listeners.delete(handler);
  }

  function connect() {
    const url = resolveUrl();

    if (!url) {
      setError(new Error("WebSocket URL is empty"));
      setStatus("error");
      return;
    }

    clearReconnect();
    closedManually = false;
    if (socket) {
      try {
        socket.close(1000, "reconnect");
      } catch {}
      socket = null;
    }

    setStatus("connecting");
    setError(null);

    try {
      const ws = new WebSocket(url, options.protocols);
      socket = ws;

      ws.onopen = () => {
        setStatus("open");
        options.onOpen?.(ws);
      };

      ws.onmessage = (event) => {
        let payload = event.data;

        if (typeof payload === "string") {
          try {
            payload = JSON.parse(payload);
          } catch {
            // keep as string
          }
        }

        setLastMessage(payload);
        listeners.forEach((fn) => fn(payload, event, ws));
        options.onMessage?.(payload, event, ws);
      };

      ws.onerror = (event) => {
        setError(event);
        setStatus("error");
        options.onError?.(event, ws);
      };

      ws.onclose = (event) => {
        setStatus("closed");
        options.onClose?.(event, ws);
        socket = null;

        if (!closedManually && options.reconnect) {
          reconnectTimer = setTimeout(connect, options.reconnectDelay ?? 1200);
        }
      };
    } catch (err) {
      setError(err);
      setStatus("error");
    }
  }

  function disconnect(code = 1000, reason = "client disconnect") {
    closedManually = true;
    clearReconnect();

    const ws = socket;
    socket = null;

    if (
      ws &&
      (ws.readyState === WebSocket.CONNECTING ||
        ws.readyState === WebSocket.OPEN)
    ) {
      try {
        ws.close(code, reason);
      } catch {}
    }
  }

  onCleanup(() => disconnect());

  return {
    connect,
    disconnect,
    sendRaw,
    sendJson,
    subscribe,
    status,
    error,
    lastMessage,
    socket: () => socket,
  };
}
