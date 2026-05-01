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
  const eventHandlers = new Set();
  const messageHandlers = new Set();
  const errorHandlers = new Set();
  const pendingEvents = [];

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
      while (pendingEvents.length > 0) {
        client.sendJson(pendingEvents.shift());
      }
    },
    onBeforeDisconnect: sendLeave,
    onMessage: (event) => {
      if (!event || typeof event !== "object") return;
      eventHandlers.forEach((handler) => handler(event));

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

    for (let i = pendingEvents.length - 1; i >= 0; i -= 1) {
      if (sameChannel(pendingEvents[i]?.channel, channel)) {
        pendingEvents.splice(i, 1);
      }
    }

    if (untrack(client.status) === "open") {
      sendLeave();
    }
    setActiveChannel(null);
  }

  function subscribeMessage(handler) {
    messageHandlers.add(handler);
    return () => messageHandlers.delete(handler);
  }

  function sendEvent(type, channel, payload) {
    connect();
    const event = {
      type,
      channel,
      payload,
    };

    if (untrack(client.status) !== "open") {
      pendingEvents.push(event);
      return;
    }
    client.sendJson(event);
  }

  function subscribeEvent(handler) {
    eventHandlers.add(handler);
    return () => eventHandlers.delete(handler);
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
    sendEvent,
    status: client.status,
    subscribeEvent,
    subscribeMessage,
    subscribeError,
  };
}

export const appWebSocket = createRoot(createAppWebSocket);
