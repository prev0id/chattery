import { For } from "solid-js";
import Toast from "../components/Toast";
import ChatHeader from "../components/ChatHeader";
import ChatMessage from "../components/ChatMessage";
import Button from "../components/Button";
import ProfilePicture from "../components/ProfilePicture";
import ChatInput from "../components/ChatInput";
import SidebarProfile from "../components/SidebarProfile";
import SidebarServer from "../components/SidebarServer";
import {
  selectedTopicID,
  setSelectedTopicID,
  servers,
  setServers,
  selectedDM,
  setSelectedDM,
} from "../stores/app";

export default function App() {
  const currentTopic = "topic_name_123";

  const messages = [
    {
      avatar: "https://github.com/identicons/prev0id.png",
      time: "Today at 15:30",
      content:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin semper purus quis velit egestas gravida. Sed pellentesque eget lacus rhoncus sagittis. Proin ornare ac velit vitae facilisis. Sed et velit vitae diam pretium tristique eget quis purus. Vestibulum tellus neque, sodales in lobortis ac, laoreet nec tellus. Nunc semper dolor vel tortor varius, a tincidunt nulla sollicitudin.",
      isOwn: true,
    },
    {
      avatar: "https://github.com/identicons/prev0id.png",
      author: "user_name",
      time: "Today at 15:31",
      content: "123 some less long ass message.",
    },
  ];

  return (
    <>
      <aside class="h-full w-98 flex bg-rose-50">
        <div class="w-18 border-r-3 flex flex-col gap-4 p-4">
          <Button sideways variant="amber">
            Direct
          </Button>
          <Button sideways variant="sky">
            Servers
          </Button>
          <button class="mt-auto hover:scale-105 transition-all duration-300 ease-in-out">
            <ProfilePicture src="https://github.com/identicons/prev0id.png" />
          </button>
        </div>
        <div class="w-80 border-r-3 flex-1 flex flex-col gap-4 p-4 bg-rose-50">
          <Button variant="amber" class="mx-4">
            Search Server
          </Button>
          <Button variant="sky" class="mx-4">
            Create Server
          </Button>
          <For each={servers}>{(server, _) => <SidebarServer serverID={server.id} />}</For>
          <SidebarProfile
            avatar="https://github.com/identicons/prev0id.png"
            name="user_name"
            lastMessage="latest message really really long ass message fucking hate them."
            unread={5}
            selected={selectedDM() === 0}
            onClick={() => setSelectedDM(0)}
          />
        </div>
      </aside>
      <main class="flex-1 flex flex-col h-full">
        <ChatHeader topicName={currentTopic}></ChatHeader>
        <div class="max-w-5xl h-full mx-auto p-4 flex flex-col gap-4 overflow-auto">
          <For each={messages} fallback={<div>No messages yet.</div>}>
            {(message, _) => <ChatMessage {...message} />}
          </For>
        </div>
        <ChatInput onSend={(text) => console.log(text)}></ChatInput>
      </main>
    </>
  );
}
