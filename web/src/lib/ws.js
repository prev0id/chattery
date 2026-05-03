import { createSignal, onCleanup } from "solid-js";

const NORMAL_CLOSE_CODE = 1000;
const RECONNECT_DELAY_MS = 1200;

export const WS_EVENT_TYPE = {
  Ping: "ping",
  Pong: "pong",
  Join: "join",
  Leave: "leave",
  Message: "message",
  Error: "error",
  VoiceOffer: "voice_offer",
  VoiceAnswer: "voice_answer",
  VoiceICECandidate: "voice_ice_candidate",
  VoiceICECandidates: "voice_ice_candidates",
  VoiceState: "voice_state",
  VoiceJoined: "voice_joined",
  VoiceLeft: "voice_left",
};

export const WS_CHANNEL_TYPE = {
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
        socket.close(NORMAL_CLOSE_CODE, "reconnect");
      } catch {
        // Closing a half-open browser WebSocket can throw; reconnect continues with a fresh socket.
      }
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
          reconnectTimer = setTimeout(
            connect,
            options.reconnectDelay ?? RECONNECT_DELAY_MS,
          );
        }
      };
    } catch (err) {
      setError(err);
      setStatus("error");
    }
  }

  function disconnect(code = NORMAL_CLOSE_CODE, reason = "client disconnect") {
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
      } catch {
        // The socket is already being discarded, so close failures are safe to ignore.
      }
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
