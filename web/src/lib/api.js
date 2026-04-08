import { action } from "@solidjs/router";
import { toast } from "../stores/toast";

async function handleResponse(res, errorMsg = "Request failed") {
  if (res.status === 401) {
    window.location.href = "/login";
    return null;
  }
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    toast.error(data.message ?? errorMsg);
    return null;
  }
  return res.json();
}

async function updateServer(serverId, name) {
  try {
    const res = await fetch("/v1/server/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId, name }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      return { error: data.message ?? "Failed to update server" };
    }

    return { ok: true };
  } catch {
    return { error: "Network error – please check your connection" };
  }
}

export const updateServerAction = action(
  async (serverID, serverName, formData) => {
    const newName = formData.get("name");

    if (newName === serverName) {
      return { ok: false, error: "Name didn't changed" };
    }

    console.log("here");

    const result = await updateServer(serverID, newName);
    if (result?.error) {
      return { ok: false, error: result.error };
    }
    return { ok: true };
  },
  "update_server_action",
);

async function createTopic(serverId, name, type) {
  try {
    const res = await fetch("/v1/server/topic/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId, name, type }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      return { error: data.message ?? "Failed to add topic" };
    }

    return { ok: true };
  } catch {
    return { error: "Network error – please check your connection" };
  }
}

export const addTopicAction = action(async (serverID, formData) => {
  const name = formData.get("name");
  const type = formData.get("type");

  const result = await createTopic(serverID, name, type);
  if (result?.error) {
    return { ok: false, error: result.error };
  }
  return { ok: true };
}, "add_topic_action");

async function updateTopic(topicId, name) {
  try {
    const res = await fetch("/v1/server/topic/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId, name }),
    });
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      return { error: data.message ?? "Failed to update topic" };
    }

    return { ok: true };
  } catch {
    return { error: "Network error – please check your connection" };
  }
}

export const updateTopicAction = action(async (topicID, formData) => {
  const name = formData.get("name");

  const result = await updateTopic(topicID, name);
  if (result?.error) {
    return { ok: false, error: result.error };
  }
  return { ok: true };
}, "update_topic_action");

async function deleteTopic(topicId) {
  try {
    const res = await fetch("/v1/server/topic/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId }),
    });
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      return { error: data.message ?? "Failed to delete topic" };
    }

    return { ok: true };
  } catch {
    return { error: "Network error – please check your connection" };
  }
}

export const deleteTopicAction = action(async (topicID) => {
  if (!confirm("Are you sure you want to delete this topic?"))
    return { ok: false };

  const result = await deleteTopic(topicID);
  if (result?.error) {
    return { ok: false, error: result.error };
  }
  return { ok: true };
}, "delete_topic_action");

async function deleteServer(serverId) {
  try {
    const res = await fetch("/v1/server/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId }),
    });
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      return { error: data.message ?? "Failed to update topic" };
    }
    return { ok: true };
  } catch {
    return { error: "Network error – please check your connection" };
  }
}

export const deleteServerAction = action(async (serverID) => {
  if (!confirm("Are you sure you want to delete this server?"))
    return { ok: false };

  const result = await deleteServer(serverID);
  if (result?.error) {
    return { ok: false, error: result.error };
  }
  return { ok: true };
}, "delete_server_action");

export async function fetchDMs() {
  try {
    const res = await fetch("/v1/dm/list");
    const data = await handleResponse(res, "Failed to load DMs");
    return data?.dms || [];
  } catch (err) {
    console.log(err);
    toast.error("Network error – please check your connection");
    return [];
  }
}

export async function fetchDMMessages(dmID, cursor = null) {
  try {
    const body = {
      cursor: cursor
        ? {
            dm_id: dmID,
            message_id: cursor.message_id,
            timestamp: cursor.timestamp,
          }
        : { dm_id: dmID },
    };
    const res = await fetch("/v1/dm/messages", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    return handleResponse(res, "Failed to load messages");
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function sendDMMessage(dmID, text) {
  try {
    const res = await fetch("/v1/dm/message", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ dm_id: dmID, text }),
    });
    return handleResponse(res, "Failed to send message");
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function fetchServers() {
  try {
    const res = await fetch("/v1/server/list");
    const data = await handleResponse(res, "Failed to load servers");
    return data?.servers || [];
  } catch (err) {
    toast.error("Network error – please check your connection");
    return [];
  }
}

export async function createServer(name) {
  try {
    const res = await fetch("/v1/server/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });
    const data = await handleResponse(res, "Failed to create server");
    if (data) toast.success("Server created!");
    return data;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
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
    return handleResponse(res, "Failed to load messages");
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
    return handleResponse(res, "Failed to send message");
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}
