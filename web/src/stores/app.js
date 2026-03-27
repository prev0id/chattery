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
    unread: 0,
    user: {
      id: 2,
      username: "user_name_2",
      profilePicture: "https://github.com/identicons/prev0id.png",
    },
    message: {
      date: "Today, 12:30",
      content: "hello! slksjf slkjfsla bllka sfsfiuhjfklsd",
    },
  },
  {
    user: {
      id: 1,
      username: "user_name_1",
      profilePicture: "https://github.com/identicons/prev0id.png",
    },
    unread: 5,
    message: {
      date: "Today, 12:30",
      content: "hello! slksjf slkjfsla bllka sfsfiuhjfklsd",
    },
  },
  {
    user: {
      id: 3,
      username: "user_name_3",
      profilePicture: "https://github.com/identicons/prev0id.png",
    },
    unread: 0,
    message: null,
  },
]);

export const [userData, setUserData] = createSignal({
  id: 123,
  username: "prevoid",
  email: "email@exmaple.com",
  profilePicture: "https://github.com/identicons/prev0id.png",
});

export const [messages, setMessages] = createStore([
  {
    user: {
      id: 123,
      username: "prevoid",
      profilePicture: "https://github.com/identicons/prev0id.png",
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
      profilePicture: "https://github.com/identicons/prev0id.png",
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
  setDMs((dm) => dm.user.id === selectedDM.user.id, "unread", 0);
  leaveTopic();
}

export function changeTab(tab) {
  setSelectedTab(tab);
  leaveTopic();
  setSelectedDM(null);
}
