import { createResource, createSignal } from "solid-js";
import { createStore } from "solid-js/store";
import { toast } from "./toast";

async function fetchUserData() {
  try {
    const res = await fetch("/v1/user/me");
    if (res.status === 401) {
      window.location.href = "/login";
      return null;
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to load user data");
      return null;
    }
    const data = await res.json();
    return data.me;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

async function fetchServers() {
  try {
    const res = await fetch("/v1/server/list");
    if (res.status === 401) {
      window.location.href = "/login";
      return [];
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to load servers");
      return [];
    }
    const data = await res.json();
    return data.servers || [];
  } catch (err) {
    toast.error("Network error – please check your connection");
    return [];
  }
}

async function fetchDMs() {
  try {
    const res = await fetch("/v1/dm/list");
    if (res.status === 401) {
      window.location.href = "/login";
      return [];
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to load DMs");
      return [];
    }
    const data = await res.json();
    return data.dms || [];
  } catch (err) {
    console.log(err);
    toast.error("Network error – please check your connection");
    return [];
  }
}

export const [userData, { refetch: refetchUserData }] =
  createResource(fetchUserData);

export const [servers, { refetch: refetchServers }] =
  createResource(fetchServers);

export const [DMs, { refetch: refetchDMs }] = createResource(fetchDMs);

export const [selectedTopic, setSelectedTopic] = createSignal(null);

export const [selectedServer, setSelectedServer] = createSignal(null);

export const [selectedServerForEdit, setSelectedServerForEdit] =
  createSignal(null);

export const [selectedTab, setSelectedTab] = createSignal("direct");

export const [selectedDM, setSelectedDM] = createSignal(null);

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

export function selectDM(selectedDM) {
  setSelectedDM(selectedDM);
  leaveTopic();
  loadChatMessages(selectedDM.id, "dm");
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

export async function createServer(name) {
  try {
    const res = await fetch("/v1/server/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to create server");
      return false;
    }

    toast.success("Server created!");
    return true;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return false;
  }
}

export async function updateServer(serverId, name) {
  try {
    const res = await fetch("/v1/server/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId, name }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to update server");
      return false;
    }

    setSelectedServerForEdit((prev) => ({ ...prev, name }));
    toast.success("Server updated!");
    return true;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return false;
  }
}

export async function deleteServer(serverId) {
  try {
    const res = await fetch("/v1/server/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to delete server");
      return false;
    }

    toast.success("Server deleted!");
    return true;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return false;
  }
}

export async function createTopic(serverId, name, type) {
  try {
    const res = await fetch("/v1/server/topic/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId, name, type }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to create topic");
      return false;
    }

    const data = await res.json();
    setSelectedServerForEdit((prev) => ({
      ...prev,
      topics: [{ id: data.id, name, type }, ...prev.topics],
    }));
    toast.success("Topic created!");
    return true;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return false;
  }
}

export async function updateTopic(topicId, name) {
  try {
    const res = await fetch("/v1/server/topic/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId, name }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to update topic");
      return false;
    }

    setSelectedServerForEdit((prev) => ({
      ...prev,
      topics: prev.topics.map((t) => (t.id === topicId ? { ...t, name } : t)),
    }));
    toast.success("Topic updated!");
    return true;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return false;
  }
}

export async function deleteTopic(topicId) {
  try {
    const res = await fetch("/v1/server/topic/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to delete topic");
      return false;
    }

    setSelectedServerForEdit((prev) => ({
      ...prev,
      topics: prev.topics.filter((t) => t.id !== topicId),
    }));
    toast.success("Topic deleted!");
    return true;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return false;
  }
}

export async function fetchTopicMessages(topicId, cursor = null) {
  try {
    const body = {
      cursor: cursor
        ? {
            topic_id: topicId,
            message_id: cursor.message_id,
            timestamp: cursor.timestamp,
          }
        : { topic_id: topicId },
    };

    const res = await fetch("/v1/server/topic/messages", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    if (res.status === 401) {
      window.location.href = "/login";
      return null;
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to load messages");
      return null;
    }
    return await res.json();
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function fetchDMMessages(dmId, cursor = null) {
  try {
    const body = {
      cursor: cursor
        ? {
            dm_id: dmId,
            message_id: cursor.message_id,
            timestamp: cursor.timestamp,
          }
        : { dm_id: dmId },
    };

    const res = await fetch("/v1/dm/messages", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    if (res.status === 401) {
      window.location.href = "/login";
      return null;
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to load messages");
      return null;
    }
    return await res.json();
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function sendTopicMessage(topicId, text) {
  try {
    const res = await fetch("/v1/server/topic/message", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId, text }),
    });

    if (res.status === 401) {
      window.location.href = "/login";
      return null;
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to send message");
      return null;
    }
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function sendDMMessage(dmId, text) {
  try {
    const res = await fetch("/v1/dm/message", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ dm_id: dmId, text }),
    });

    if (res.status === 401) {
      window.location.href = "/login";
      return null;
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to send message");
      return null;
    }
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}
