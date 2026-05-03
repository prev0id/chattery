import {
  createEffect,
  For,
  on,
  onCleanup,
  Show,
  createSignal,
  untrack,
} from "solid-js";
import {
  getDmMessages,
  markDmRead,
  sendDmMessage,
} from "~/features/dm/api";
import { DM_CHAT_TARGET, DM_MESSAGES } from "~/features/dm/constants";
import {
  getTopicMessages,
  sendTopicMessage,
} from "~/features/server/api";
import {
  SERVER_CHAT_TARGET,
  SERVER_MESSAGES,
} from "~/features/server/constants";
import { appWebSocket } from "~/stores/websocket";
import { getUserErrorMessage } from "~/shared/api/errors";
import { toast } from "../stores/toast";
import ChatInput from "./ChatInput";
import ChatMessage from "./ChatMessage";
import { ArrowLeftToLine } from "lucide-solid";

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
          ? await getTopicMessages(chatId)
          : await getDmMessages(chatId);

      setMessages(normalizeMessages(response));
      setMessagesCursor(response?.cursor ?? null);

      scrollToBottom(false);
    } catch (error) {
      const fallback =
        chatType === "topic"
          ? SERVER_MESSAGES.topicMessagesFailed
          : DM_MESSAGES.messagesFailed;
      toast.error(getUserErrorMessage(error, fallback));
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
          ? await getTopicMessages(chat.id, cursor)
          : await getDmMessages(chat.id, cursor);

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
    } catch (error) {
      const fallback =
        chat.type === "topic"
          ? SERVER_MESSAGES.topicMessagesFailed
          : DM_MESSAGES.messagesFailed;
      toast.error(getUserErrorMessage(error, fallback));
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

  createEffect(
    on(
      () => {
        const chatId = props.chatID;
        const channelType = props.channelType;
        const chatType = props.type === DM_CHAT_TARGET ? "dm" : "topic";
        return chatId && channelType
          ? `${chatType}:${channelType}:${Number(chatId)}`
          : "";
      },
      () => {
        const chatId = untrack(() => props.chatID);
        const channelType = untrack(() => props.channelType);
        if (!chatId || !channelType) return;

        const chatType = untrack(() =>
          props.type === DM_CHAT_TARGET ? "dm" : "topic",
        );
        loadChatMessages(chatId, chatType);

        const channel = { type: channelType, id: Number(chatId) };
        appWebSocket.join(channel);

        const unsubscribeMessage = appWebSocket.subscribeMessage(
          ({ channel: eventChannel, payload }) => {
            if (
              eventChannel?.type !== channel.type ||
              eventChannel?.id !== channel.id
            ) {
              return;
            }
            addMessage(payload);
            props.onMessage?.(payload);

            if (chatType === "dm" && payload?.id) {
              markDmRead(chatId, payload.id).catch((error) => {
                toast.error(
                  getUserErrorMessage(error, DM_MESSAGES.markReadFailed),
                );
              });
            }
          },
        );
        const unsubscribeError = appWebSocket.subscribeError((payload) => {
          toast.error(payload?.message ?? "WebSocket error");
        });

        onCleanup(() => {
          unsubscribeMessage();
          unsubscribeError();
          appWebSocket.leave(channel);
        });
      },
    ),
  );

  const handleSend = async (text) => {
    const chat = currentChat();
    if (!chat || sending()) return;

    setSending(true);
    try {
      if (chat.type === "topic") {
        return Boolean(await sendTopicMessage(chat.id, text));
      }
      return Boolean(await sendDmMessage(chat.id, text));
    } catch (error) {
      const fallback =
        chat.type === "topic"
          ? SERVER_MESSAGES.sendTopicMessageFailed
          : DM_MESSAGES.sendFailed;
      toast.error(getUserErrorMessage(error, fallback));
      return false;
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
