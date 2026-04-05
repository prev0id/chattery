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

export const [messages, setMessages] = createStore([
  {
    user: {
      id: 12,
      username: "prevoid",
      avatar: "https://github.com/identicons/prev0id.png",
    },
    message: {
      date: "Today at 15:31",
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin semper purus quis velit egestas gravida. Sed pellentesque eget lacus rhoncus sagittis. Proin ornare ac velit vitae facilisis. Sed et velit vitae diam pretium tristique eget quis purus. Vestibulum tellus neque, sodales in lobortis ac, laoreet nec tellus. Nunc semper dolor vel tortor varius, a tincidunt nulla sollicitudin.",
    },
    lastSeenMessage: true,
  },
  {
    user: {
      id: 123,
      username: "prevoid",
      avatar: "https://github.com/identicons/prev0id.png",
    },
    message: {
      date: "Today at 15:30",
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin semper purus quis velit egestas gravida. Sed pellentesque eget lacus rhoncus sagittis. Proin ornare ac velit vitae facilisis. Sed et velit vitae diam pretium tristique eget quis purus. Vestibulum tellus neque, sodales in lobortis ac, laoreet nec tellus. Nunc semper dolor vel tortor varius, a tincidunt nulla sollicitudin.",
    },
  },
  {
    user: {
      id: 312,
      username: "user_name",
      avatar: "https://github.com/identicons/prev0id.png",
    },
    message: {
      date: "Today at 15:31",
      content: "123 some less long ass message.",
    },
  },
]);

export function selectTopic(topic, server) {
  setSelectedTopic(topic);
  setSelectedServer(server);
  setSelectedDM(null);
}

export function leaveTopic() {
  setSelectedTopic(null);
  setSelectedServer(null);
}

export function selectDM(selectedDM) {
  setSelectedDM(selectedDM);
  leaveTopic();
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
