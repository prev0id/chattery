import { createResource, createSignal } from "solid-js";
import { toast } from "./toast";
import { leaveTopic, loadChatMessages } from "./app";

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

export const [DMs, { refetch: refetchDMs }] = createResource(fetchDMs);

export const [selectedDM, setSelectedDM] = createSignal(null);

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

export async function sendDMMessage(dmID, text) {
  try {
    const res = await fetch("/v1/dm/message", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ dm_id: dmID, text }),
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

export function selectDM(selectedDM) {
  setSelectedDM(selectedDM);
  leaveTopic();
  loadChatMessages(selectedDM.id, "dm");
}
