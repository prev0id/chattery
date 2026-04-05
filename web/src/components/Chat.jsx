import {
  createEffect,
  onCleanup,
  For,
  Show,
  createSignal,
  onMount,
} from "solid-js";
import { createStore } from "solid-js/store";
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
  const [messages, setMessages] = createStore([]);
  const [messagesCursor, setMessagesCursor] = createSignal(null);
  const [currentChat, setCurrentChat] = createSignal(null);
  const [loading, setLoading] = createSignal(false);

  let sentinelRef;
  let containerRef;

  const loadChatMessages = async (chatId, chatType) => {
    setCurrentChat({ id: chatId, type: chatType });
    setMessages([]);
    setMessagesCursor(null);
    setLoading(true);

    let response;
    if (chatType === "topic") {
      response = await fetchTopicMessages(chatId);
    } else {
      response = await fetchDMMessages(chatId);
    }

    setLoading(false);
    if (response && response.messages) {
      setMessages(response.messages);
      setMessagesCursor(response.cursor);
    }
  };

  const loadMoreMessages = async () => {
    const chat = currentChat();
    const cursor = messagesCursor();
    if (!chat || !cursor || loading()) return;

    setLoading(true);
    let response;
    if (chat.type === "topic") {
      response = await fetchTopicMessages(chat.id, cursor);
    } else {
      response = await fetchDMMessages(chat.id, cursor);
    }
    setLoading(false);

    if (response && response.messages && response.messages.length > 0) {
      setMessages((prev) => [...response.messages, ...prev]);
      setMessagesCursor(response.cursor);
    }
  };

  createEffect(() => {
    const chatId = props.chatID;
    const chatType = props.type;
    if (!chatId) return;

    const messageType = chatType === DMsType ? "dm" : "topic";
    loadChatMessages(chatId, messageType);
  });

  const hasChat = () => !!props.chatID;

  createEffect(() => {
    if (!sentinelRef) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && currentChat()) {
          loadMoreMessages();
        }
      },
      { root: containerRef, threshold: 0.1 },
    );

    observer.observe(sentinelRef);

    onCleanup(() => observer.disconnect());
  });

  const handleSend = async (text) => {
    const chat = currentChat();
    if (!chat) return;

    if (chat.type === "topic") {
      await sendTopicMessage(chat.id, text);
    } else {
      await sendDMMessage(chat.id, text);
    }
    loadChatMessages(chat.id, chat.type);
  };

  return (
    <Show when={hasChat()} fallback={selectChatMessage}>
      <div
        class="max-w-5xl min-w-sm w-full h-full mx-auto p-4 flex flex-col overflow-auto"
        ref={containerRef}
      >
        <div ref={sentinelRef} class="h-1" />
        <Show when={loading()}>
          <LoadingSpinner />
        </Show>
        <For each={messages} fallback={<div>No messages yet.</div>}>
          {(message) => <ChatMessage msg={message} />}
        </For>
      </div>
      <ChatInput onSend={handleSend} />
    </Show>
  );
}

function selectChatMessage() {
  return (
    <div class="m-auto flex items-center justify-center gap-4">
      <ArrowLeftToLine class="size-10" />
      <span class="text-2xl tracking-wider font-semibold">
        Select a chat to start messaging
      </span>
    </div>
  );
}

function LoadingSpinner() {
  const [show, setShow] = createSignal(false);

  onMount(() => {
    const timer = setTimeout(() => setShow(true), 300);
    return () => clearTimeout(timer);
  });

  return (
    <Show when={show()}>
      {props.children}
      <div class="m-auto flex items-center gap-4">
        <Loader class="size-10 animate-spin" />
        <span class="tracking-wider text-2xl font-semibold">
          Loading chat...
        </span>
      </div>
    </Show>
  );
}
