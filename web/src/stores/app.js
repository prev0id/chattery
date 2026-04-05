import { createSignal } from "solid-js";
import { createStore } from "solid-js/store";
import { fetchTopicMessages } from "./server";
import { fetchDMMessages } from "./dm";
import { setSelectedTopic, setSelectedServer, refetchServers } from "./server";
import { setSelectedDM, refetchDMs } from "./dm";

export const [selectedTab, setSelectedTab] = createSignal("direct");

export const [messages, setMessages] = createStore([]);

export const [messagesCursor, setMessagesCursor] = createSignal(null);

export const [currentChat, setCurrentChat] = createSignal(null);

export async function loadChatMessages(chatId, chatType) {
  setCurrentChat({ id: chatId, type: chatType });
  setMessages([]);
  setMessagesCursor(null);

  let response;
  if (chatType === "topic") {
    response = await fetchTopicMessages(chatId);
  } else {
    response = await fetchDMMessages(chatId);
  }

  if (response && response.messages) {
    setMessages(response.messages);
    setMessagesCursor(response.cursor);
  } else {
    console.log("No messages in response or response is null");
  }
}

export async function loadMoreMessages() {
  const chat = currentChat();
  const cursor = messagesCursor();
  if (!chat || !cursor) return;

  let response;
  if (chat.type === "topic") {
    response = await fetchTopicMessages(chat.id, cursor);
  } else {
    response = await fetchDMMessages(chat.id, cursor);
  }

  if (response && response.messages && response.messages.length > 0) {
    setMessages((prev) => [...response.messages, ...prev]);
    setMessagesCursor(response.cursor);
  }
}

export function addMessage(message) {
  setMessages((prev) => [...prev, message]);
}

export function selectTopic(topic, server) {
  setSelectedTopic(topic);
  setSelectedServer(server);
  setSelectedDM(null);
  loadChatMessages(topic.id, "topic");
}

export function leaveTopic() {
  setSelectedTopic(null);
  setSelectedServer(null);
  setCurrentChat(null);
  setMessages([]);
  setMessagesCursor(null);
}

export function changeTab(tab) {
  setSelectedTab(tab);
  leaveTopic();
  setSelectedDM(null);
  if (tab === "servers") {
    refetchServers();
  }
  if (tab === "direct") {
    refetchDMs();
  }
}
