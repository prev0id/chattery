import { createResource, createSignal } from "solid-js";
import { toast } from "./toast";

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

export const [servers, { refetch: refetchServers }] =
  createResource(fetchServers);

export const [selectedServer, setSelectedServer] = createSignal(null);

export const [selectedServerForEdit, setSelectedServerForEdit] =
  createSignal(null);

export const [selectedTopic, setSelectedTopic] = createSignal(null);

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
