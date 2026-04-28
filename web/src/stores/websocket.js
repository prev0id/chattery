import { createRoot, createSignal, untrack } from "solid-js";
import {
  createWebSocketClient,
  websocketURL,
  WSEventType,
} from "~/lib/ws";

function parsePayload(payload) {
  if (typeof payload !== "string") return payload;

  try {
    return JSON.parse(payload);
  } catch {
    return payload;
  }
}

function sameChannel(left, right) {
  return (
    left?.type === right?.type && Number(left?.id) === Number(right?.id)
  );
}

function createAppWebSocket() {
  const [activeChannel, setActiveChannel] = createSignal(null);
  const messageHandlers = new Set();
  const errorHandlers = new Set();

  let client;

  function sendJoin(channel) {
    if (!channel) return;
    client.sendJson({
      type: WSEventType.Join,
      payload: channel,
    });
  }

  function sendLeave() {
    client.sendJson({ type: WSEventType.Leave });
  }

  client = createWebSocketClient(websocketURL, {
    reconnect: true,
    onOpen: () => {
      sendJoin(activeChannel());
    },
    onBeforeDisconnect: sendLeave,
    onMessage: (event) => {
      if (!event || typeof event !== "object") return;

      if (event.type === WSEventType.Ping) {
        client.sendJson({ type: WSEventType.Pong });
        return;
      }

      if (event.type === WSEventType.Error) {
        const payload = parsePayload(event.payload);
        errorHandlers.forEach((handler) => handler(payload));
        return;
      }

      if (event.type === WSEventType.Message) {
        const message = {
          channel: event.channel,
          payload: parsePayload(event.payload),
        };
        messageHandlers.forEach((handler) => handler(message));
      }
    },
  });

  function connect() {
    const status = untrack(client.status);
    if (status === "open" || status === "connecting") {
      return;
    }
    client.connect();
  }

  function join(channel) {
    connect();
    const previous = untrack(activeChannel);
    if (sameChannel(previous, channel)) return;

    if (previous && untrack(client.status) === "open") {
      sendLeave();
    }

    setActiveChannel(channel);

    if (untrack(client.status) === "open") {
      sendJoin(channel);
    }
  }

  function leave(channel) {
    const current = untrack(activeChannel);
    if (!sameChannel(current, channel)) return;

    if (untrack(client.status) === "open") {
      sendLeave();
    }
    setActiveChannel(null);
  }

  function subscribeMessage(handler) {
    messageHandlers.add(handler);
    return () => messageHandlers.delete(handler);
  }

  function subscribeError(handler) {
    errorHandlers.add(handler);
    return () => errorHandlers.delete(handler);
  }

  connect();

  return {
    activeChannel,
    connect,
    join,
    leave,
    status: client.status,
    subscribeMessage,
    subscribeError,
  };
}

export const appWebSocket = createRoot(createAppWebSocket);
