import { createSignal } from "solid-js";
import Toast from "../components/Toast";
import ChatHeader from "../components/ChatHeader";
import ChatMessage from "../components/ChatMessage";
import Button from "../components/Button";
import ProfilePicture from "../components/ProfilePicture";
import ChatInput from "../components/ChatInput";
import SidebarProfile from "../components/SidebarProfile";
import SidebarServer from "../components/SidebarServer";

export default function App() {
  const [selectedTopicID, setSelectedTopicID] = createSignal(null);
  const [selectedDM, setSelectedDM] = createSignal(0);

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

  const servers = [
    {
      name: "Server Name",
      topics: [
        { id: 1, name: "General", type: "text" },
        { id: 2, name: "Memes", type: "text" },
        { id: 3, name: "Voice Lounge", type: "voice" },
      ],
    },
    {
      name: "Gaming Hub",
      topics: [
        { id: 4, name: "Valorant", type: "text" },
        { id: 5, name: "Among Us", type: "voice" },
      ],
    },
  ];

  return (
    <>
      <aside class="h-full w-98 flex bg-rose-50">
        <div class="w-18 border-r-3 flex flex-col gap-4 p-4">
          <Button sideways variant="amber">
            Private
          </Button>
          <Button sideways variant="sky">
            Private
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
          {servers.map((server) => (
            <SidebarServer
              name={server.name}
              topics={server.topics}
              selectedTopicId={selectedTopicID()}
              onTopicSelect={setSelectedTopicID}
            />
          ))}
          <details class="border-2 rounded-lg p-1 bg-white" open>
            <summary class="px-2 flex justify-between items-center hover:bg-emerald-200 rounded-lg border-2 border-white hover:border-black transition-all duration-300 ease-in-out">
              <h2 class="text-lg font-semibold">Server Name</h2>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-settings2-icon lucide-settings-2"
              >
                <path d="M14 17H5" />
                <path d="M19 7h-9" />
                <circle cx="17" cy="17" r="3" />
                <circle cx="7" cy="7" r="3" />
              </svg>
            </summary>
            <hr class="mt-1" />
            <div class="flex items-center gap-1 px-2 my-1 py-0.5 border-2 border-white hover:border-black rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-messages-square-icon lucide-messages-square"
              >
                <path d="M16 10a2 2 0 0 1-2 2H6.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 2 14.286V4a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
                <path d="M20 9a2 2 0 0 1 2 2v10.286a.71.71 0 0 1-1.212.502l-2.202-2.202A2 2 0 0 0 17.172 19H10a2 2 0 0 1-2-2v-1" />
              </svg>
              <p>Topic name</p>
            </div>
            <div class="flex items-center gap-1 px-2 my-1 py-0.5 border-2 border-black neo-shadow-white rounded-lg hover:neo-shadow transition-all duration-300 ease-in-out bg-emerald-200">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-messages-square-icon lucide-messages-square"
              >
                <path d="M16 10a2 2 0 0 1-2 2H6.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 2 14.286V4a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
                <path d="M20 9a2 2 0 0 1 2 2v10.286a.71.71 0 0 1-1.212.502l-2.202-2.202A2 2 0 0 0 17.172 19H10a2 2 0 0 1-2-2v-1" />
              </svg>
              <p>Topic name</p>
            </div>
            <hr />
            <div class="flex items-center gap-1 px-2 my-1 py-0.5 border-2 border-white hover:border-black neo-shadow-white rounded-lg hover:neo-shadow transition-all duration-300 ease-in-out">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-mic-icon lucide-mic"
              >
                <path d="M12 19v3" />
                <path d="M19 10v2a7 7 0 0 1-14 0v-2" />
                <rect x="9" y="2" width="6" height="13" rx="3" />
              </svg>
              <p>Topic name</p>
            </div>
          </details>
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
          {messages.map((msg) => (
            <ChatMessage {...msg} />
          ))}
        </div>
        <ChatInput onSend={(text) => console.log(text)}></ChatInput>
      </main>
    </>
  );
}
