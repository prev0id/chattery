import { createEffect, onCleanup, For, Show } from "solid-js";
import {
  messages,
  loadMoreMessages,
  currentChat,
  loadChatMessages,
} from "../stores/app";
import { sendTopicMessage } from "../stores/server";
import { sendDMMessage } from "../stores/dm";
import ChatInput from "./ChatInput";
import ChatMessage from "./ChatMessage";
import { ArrowLeftToLine, Loader } from "lucide-solid";

const DMsType = "dms";
const ServersType = "servers";
const TopicTypeText = "text";
const TopicTypeVoice = "voice";

export { DMsType, ServersType, TopicTypeText, TopicTypeVoice };

export function Chat(props) {
  let sentinelRef;
  let containerRef;

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
  };

  return (
    <Show when={hasChat()} fallback={selectChatMessage}>
      <div
        class="max-w-5xl min-w-sm w-full h-full mx-auto p-4 flex flex-col overflow-auto"
        ref={containerRef}
      >
        <div ref={sentinelRef} class="h-1" />
        <For each={messages} fallback={<div>No messages yet.</div>}>
          {(message, _) => <ChatMessage msg={message} />}
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
  return (
    <div class="m-auto flex items-center gap-4 text-2xl tracking-wider font-semibold">
      <Loader class="size-10 animate-spin" />
      Loading...
    </div>
  );
}
