import { createSignal, createEffect } from "solid-js";
import { createStore } from "solid-js/store";

// TODO resource loading

export const [selectedTopicID, setSelectedTopicID] = createSignal(-1);

export const [selectedServerID, setSelectedServerID] = createSignal(-1);

export const [servers, setServers] = createStore([
  {
    id: 1,
    name: "Server Name",
    topics: [
      { id: 1, name: "General", type: "text" },
      { id: 2, name: "Memes", type: "text" },
      { id: 3, name: "Voice Lounge", type: "voice" },
    ],
  },
  {
    id: 2,
    name: "Gaming Hub",
    topics: [
      { id: 4, name: "Valorant", type: "text" },
      { id: 5, name: "Among Us", type: "voice" },
    ],
  },
]);

export const [selectedDM, setSelectedDM] = createSignal(-1);

export const [DMs, setDMs] = createStore(null);

export function selectTopic(topicID, serverID) {
  setSelectedTopicID(topicID);
  setSelectedServerID(serverID);
}
