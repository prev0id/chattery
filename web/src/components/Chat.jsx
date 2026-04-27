import { createEffect, For, onCleanup, Show, createSignal } from "solid-js";
import {
  fetchTopicMessages,
  sendTopicMessage,
  fetchDMMessages,
  sendDMMessage,
} from "../lib/api";
import { createChatWebSocketClient } from "../lib/ws";
import { toast } from "../stores/toast";
import ChatInput from "./ChatInput";
import ChatMessage from "./ChatMessage";
import { ArrowLeftToLine } from "lucide-solid";

const DMsType = "dms";
const ServersType = "servers";
const TopicTypeText = "text";
const TopicTypeVoice = "voice";

export { DMsType, ServersType, TopicTypeText, TopicTypeVoice };

export function Chat(props) {
  const [messages, setMessages] = createSignal([]);
  const [messagesCursor, setMessagesCursor] = createSignal(null);
  const [currentChat, setCurrentChat] = createSignal(null);
  const [loading, setLoading] = createSignal(false);
  const [sending, setSending] = createSignal(false);

  let containerRef;
  let bottomRef;

  const scrollToBottom = (smooth = false) => {
    requestAnimationFrame(() => {
      bottomRef?.scrollIntoView({
        behavior: smooth ? "smooth" : "auto",
        block: "end",
      });
    });
  };

  const isScrolledToBottom = () => {
    const el = containerRef;
    if (!el) return true;
    return el.scrollHeight - el.scrollTop - el.clientHeight < 80;
  };

  const normalizeMessages = (response) =>
    response?.messages?.slice().reverse() ?? [];

  const loadChatMessages = async (chatId, chatType) => {
    setCurrentChat({ id: chatId, type: chatType });
    setLoading(true);

    try {
      const response =
        chatType === "topic"
          ? await fetchTopicMessages(chatId)
          : await fetchDMMessages(chatId);

      setMessages(normalizeMessages(response));
      setMessagesCursor(response?.cursor ?? null);

      scrollToBottom(false);
    } finally {
      setLoading(false);
    }
  };

  const loadMoreMessages = async () => {
    const chat = currentChat();
    const cursor = messagesCursor();
    const el = containerRef;

    if (!chat || !cursor || loading() || !el) return;

    const prevScrollHeight = el.scrollHeight;
    const prevScrollTop = el.scrollTop;

    setLoading(true);

    try {
      const response =
        chat.type === "topic"
          ? await fetchTopicMessages(chat.id, cursor)
          : await fetchDMMessages(chat.id, cursor);

      const olderMessages = normalizeMessages(response);
      if (olderMessages.length > 0) {
        setMessages((prev) => [...olderMessages, ...prev]);
        setMessagesCursor(response?.cursor ?? null);

        requestAnimationFrame(() => {
          const newScrollHeight = el.scrollHeight;
          el.scrollTop = newScrollHeight - prevScrollHeight + prevScrollTop;
        });
      } else {
        setMessagesCursor(null);
      }
    } finally {
      setLoading(false);
    }
  };

  const addMessage = (message) => {
    if (!message) return;

    const shouldScroll = isScrolledToBottom();

    setMessages((prev) => {
      if (prev.some((item) => item.id === message.id)) {
        return prev;
      }
      return [...prev, message];
    });

    if (shouldScroll) {
      scrollToBottom(true);
    }
  };

  const handleScroll = () => {
    const el = containerRef;
    if (!el) return;

    if (el.scrollTop <= 50) {
      loadMoreMessages();
    }
  };

  createEffect(() => {
    const chatId = props.chatID;
    const channelType = props.channelType;
    if (!chatId || !channelType) return;

    const chatType = props.type === DMsType ? "dm" : "topic";
    loadChatMessages(chatId, chatType);

    const channel = { type: channelType, id: Number(chatId) };
    const ws = createChatWebSocketClient({
      channel,
      onMessage: ({ channel: eventChannel, payload }) => {
        if (
          eventChannel?.type !== channel.type ||
          eventChannel?.id !== channel.id
        ) {
          return;
        }
        addMessage(payload);
      },
      onError: (payload) => {
        toast.error(payload?.message ?? "WebSocket error");
      },
    });

    ws.connect();

    onCleanup(() => {
      ws.disconnect(1000, "leave chat");
    });
  });

  const handleSend = async (text) => {
    const chat = currentChat();
    if (!chat || sending()) return;

    setSending(true);
    try {
      if (chat.type === "topic") {
        return Boolean(await sendTopicMessage(chat.id, text));
      }
      return Boolean(await sendDMMessage(chat.id, text));
    } finally {
      setSending(false);
    }
  };

  return (
    <Show
      when={props.chatID}
      fallback={
        <div class="m-auto flex items-center justify-center gap-4">
          <ArrowLeftToLine class="size-10" />
          <span class="text-2xl tracking-wider font-semibold">
            Select a chat to start messaging
          </span>
        </div>
      }
    >
      <div class="w-full h-full flex flex-col min-h-0">
        <div
          ref={containerRef}
          onScroll={handleScroll}
          class="flex-1 min-h-0 overflow-y-auto p-4"
        >
          <Show when={loading()}>
            <div class="mb-4">Loading chat...</div>
          </Show>

          <For each={messages()}>
            {(message) => <ChatMessage msg={message} />}
          </For>

          <div ref={bottomRef} />
        </div>

        <ChatInput onSend={handleSend} disabled={sending()} />
      </div>
    </Show>
  );
}
