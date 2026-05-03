import { createRoot, createSignal, untrack } from "solid-js";
import {
  createWebSocketClient,
  websocketURL,
  WS_EVENT_TYPE,
} from "~/lib/ws";

const MAX_PENDING_EVENTS = 100;

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

function channelKey(channel) {
  return channel ? `${channel.type}:${Number(channel.id)}` : "";
}

function createAppWebSocket() {
  const [activeChannel, setActiveChannel] = createSignal(null);
  const eventHandlers = new Set();
  const messageHandlers = new Set();
  const errorHandlers = new Set();
  const pendingEvents = [];
  const activeChannels = new Map();

  let client;

  function sendJoin(channel) {
    if (!channel) return;
    client.sendJson({
      type: WS_EVENT_TYPE.Join,
      payload: channel,
    });
  }

  function sendLeave(channel) {
    client.sendJson(
      channel
        ? { type: WS_EVENT_TYPE.Leave, payload: channel }
        : { type: WS_EVENT_TYPE.Leave },
    );
  }

  function syncActiveChannelSignal() {
    setActiveChannel(activeChannels.values().next().value ?? null);
  }

  client = createWebSocketClient(websocketURL, {
    reconnect: true,
    onOpen: () => {
      activeChannels.forEach((channel) => sendJoin(channel));
      while (pendingEvents.length > 0) {
        client.sendJson(pendingEvents.shift());
      }
    },
    onBeforeDisconnect: () => {
      activeChannels.forEach((channel) => sendLeave(channel));
    },
    onMessage: (event) => {
      if (!event || typeof event !== "object") return;
      eventHandlers.forEach((handler) => {
        Promise.resolve(handler(event)).catch((error) => {
          errorHandlers.forEach((errorHandler) => errorHandler(error));
        });
      });

      if (event.type === WS_EVENT_TYPE.Ping) {
        client.sendJson({ type: WS_EVENT_TYPE.Pong });
        return;
      }

      if (event.type === WS_EVENT_TYPE.Error) {
        const payload = parsePayload(event.payload);
        errorHandlers.forEach((handler) => handler(payload));
        return;
      }

      if (event.type === WS_EVENT_TYPE.Message) {
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
    const key = channelKey(channel);
    if (!key || activeChannels.has(key)) return;

    activeChannels.set(key, channel);
    syncActiveChannelSignal();
    if (untrack(client.status) === "open") {
      sendJoin(channel);
    }
  }

  function leave(channel) {
    const key = channelKey(channel);
    if (!key || !activeChannels.has(key)) return;

    for (let i = pendingEvents.length - 1; i >= 0; i -= 1) {
      if (sameChannel(pendingEvents[i]?.channel, channel)) {
        pendingEvents.splice(i, 1);
      }
    }

    activeChannels.delete(key);
    if (untrack(client.status) === "open") {
      sendLeave(channel);
    }
    syncActiveChannelSignal();
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
      if (pendingEvents.length >= MAX_PENDING_EVENTS) {
        pendingEvents.shift();
      }
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
