import { createSignal, onCleanup } from "solid-js";

export const WSEventType = {
  Ping: "ping",
  Pong: "pong",
  Join: "join",
  Leave: "leave",
  Message: "message",
  Error: "error",
  VoiceOffer: "voice_offer",
  VoiceAnswer: "voice_answer",
  VoiceICECandidate: "voice_ice_candidate",
  VoiceJoined: "voice_joined",
  VoiceLeft: "voice_left",
};

export const WSChannelType = {
  DM: "dm",
  TextTopic: "text_topic",
  VoiceTopic: "voice_topic",
};

export function websocketURL(path = "/ws/") {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}${path}`;
}

export function createWebSocketClient(getUrl, options = {}) {
  const [status, setStatus] = createSignal("idle");
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
            // keep raw payload
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

    if (ws?.readyState === WebSocket.OPEN) {
      options.onBeforeDisconnect?.(ws);
    }

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

function parsePayload(payload) {
  if (typeof payload !== "string") return payload;

  try {
    return JSON.parse(payload);
  } catch {
    return payload;
  }
}

export function createChatWebSocketClient({ channel, onMessage, onError }) {
  const client = createWebSocketClient(websocketURL, {
    reconnect: true,
    onOpen: () => {
      client.sendJson({
        type: WSEventType.Join,
        payload: channel,
      });
    },
    onBeforeDisconnect: () => {
      client.sendJson({ type: WSEventType.Leave });
    },
    onMessage: (event) => {
      if (!event || typeof event !== "object") return;

      if (event.type === WSEventType.Ping) {
        client.sendJson({ type: WSEventType.Pong });
        return;
      }

      if (event.type === WSEventType.Error) {
        onError?.(parsePayload(event.payload));
        return;
      }

      if (event.type === WSEventType.Message) {
        onMessage?.({
          channel: event.channel,
          payload: parsePayload(event.payload),
        });
      }
    },
  });

  return client;
}
