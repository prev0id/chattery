import { action, redirect } from "@solidjs/router";

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
  return redirect("/server");
}, "delete_server_action");
