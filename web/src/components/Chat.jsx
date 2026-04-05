import { createEffect, onCleanup, For } from "solid-js";
import {
  messages,
  loadMoreMessages,
  currentChat,
  sendTopicMessage,
  sendDMMessage,
} from "../stores/app";
import ChatInput from "./ChatInput";
import ChatMessage from "./ChatMessage";

export default function Chat(props) {
  let sentinelRef;
  let containerRef;

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
    <>
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
    </>
  );
}
