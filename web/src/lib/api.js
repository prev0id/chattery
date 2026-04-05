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

export async function updateServer(serverId, name) {
  try {
    const res = await fetch("/v1/server/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId, name }),
    });
    const data = await handleResponse(res, "Failed to update server");
    if (data) toast.success("Server updated!");
    return data;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function deleteServer(serverId) {
  try {
    const res = await fetch("/v1/server/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId }),
    });
    const data = await handleResponse(res, "Failed to delete server");
    if (data) toast.success("Server deleted!");
    return data;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function createTopic(serverId, name, type) {
  try {
    const res = await fetch("/v1/server/topic/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ server_id: serverId, name, type }),
    });
    const data = await handleResponse(res, "Failed to create topic");
    if (data) toast.success("Topic created!");
    return data;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function updateTopic(topicId, name) {
  try {
    const res = await fetch("/v1/server/topic/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId, name }),
    });
    const data = await handleResponse(res, "Failed to update topic");
    if (data) toast.success("Topic updated!");
    return data;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export async function deleteTopic(topicId) {
  try {
    const res = await fetch("/v1/server/topic/delete", {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ topic_id: topicId }),
    });
    const data = await handleResponse(res, "Failed to delete topic");
    if (data) toast.success("Topic deleted!");
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
