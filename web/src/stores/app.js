import { createSignal, createEffect } from "solid-js";
import { createStore } from "solid-js/store";

// TODO resource loading

export const [selectedTopic, setSelectedTopic] = createSignal(null);

export const [selectedServer, setSelectedServer] = createSignal(null);

export const [selectedTab, setSelectedTab] = createSignal("direct");

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

export const [selectedDM, setSelectedDM] = createSignal(null);

export const [DMs, setDMs] = createStore([
  {
    id: 2,
    username: "user_name_2",
    profilePicture: "https://github.com/identicons/prev0id.png",
    unread: 0,
    lastMessage: {
      date: "Today, 12:30",
      message: "hello! slksjf slkjfsla bllka sfsfiuhjfklsd",
    },
  },
  {
    id: 1,
    username: "user_name_1",
    profilePicture: "https://github.com/identicons/prev0id.png",
    unread: 5,
    lastMessage: {
      date: "Today, 12:30",
      message: "hello! slksjf slkjfsla bllka sfsfiuhjfklsd",
    },
  },
  {
    id: 3,
    username: "user_name_3",
    unread: 0,
    latestMessage: null,
    profilePicture: "https://github.com/identicons/prev0id.png",
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
  setDMs((dm) => dm.id === selectedDM.id, "unread", 0);
  leaveTopic();
}

export function changeTab(tab) {
  setSelectedTab(tab);
  leaveTopic();
  setSelectedDM(null);
}
