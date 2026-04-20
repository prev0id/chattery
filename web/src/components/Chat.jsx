import { createEffect, For, Show, createSignal } from "solid-js";
import {
  fetchTopicMessages,
  sendTopicMessage,
  fetchDMMessages,
  sendDMMessage,
} from "../lib/api";
import ChatInput from "./ChatInput";
import ChatMessage from "./ChatMessage";
import { ArrowLeftToLine, Loader } from "lucide-solid";

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

  const loadChatMessages = async (chatId, chatType) => {
    setCurrentChat({ id: chatId, type: chatType });
    setLoading(true);

    try {
      const response =
        chatType === "topic"
          ? await fetchTopicMessages(chatId)
          : await fetchDMMessages(chatId);

      setMessages(response?.messages.reverse() ?? []);
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

      const olderMessages = response?.messages.reverse() ?? [];
      console.log(cursor);
      if (olderMessages.length > 0) {
        setMessages((prev) => [...olderMessages, ...prev]);
        setMessagesCursor(response?.cursor ?? null);

        requestAnimationFrame(() => {
          const newScrollHeight = el.scrollHeight;
          el.scrollTop = newScrollHeight - prevScrollHeight + prevScrollTop;
        });
      }
    } finally {
      setLoading(false);
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
    if (!chatId) return;

    const chatType = props.type === DMsType ? "dm" : "topic";
    loadChatMessages(chatId, chatType);
  });

  const handleSend = async (text) => {
    const chat = currentChat();
    if (!chat) return;

    if (chat.type === "topic") {
      await sendTopicMessage(chat.id, text);
    } else {
      await sendDMMessage(chat.id, text);
    }

    await loadChatMessages(chat.id, chat.type);
    scrollToBottom(true);
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

        <ChatInput onSend={handleSend} />
      </div>
    </Show>
  );
}
