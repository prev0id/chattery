import { For, Index, Match, Switch } from "solid-js";
import AppHeader from "../components/AppHeader";
import ChatMessage from "../components/ChatMessage";
import Button from "../components/Button";
import ProfilePicture from "../components/ProfilePicture";
import ChatInput from "../components/ChatInput";
import SidebarDM from "../components/SidebarDM";
import SidebarServer from "../components/SidebarServer";
import { servers, selectedTab, DMs, changeTab } from "../stores/app";

export default function App() {
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
          <Button sideways variant="amber" onClick={() => changeTab("direct")}>
            Direct
          </Button>
          <Button sideways variant="sky" onClick={() => changeTab("servers")}>
            Servers
          </Button>
          <button class="mt-auto hover:scale-105 transition-all duration-300 ease-in-out">
            <ProfilePicture src="https://github.com/identicons/prev0id.png" />
          </button>
        </div>
        <div class="w-80 border-r-3 flex-1 flex flex-col gap-4 p-4 bg-rose-50">
          <Switch>
            <Match when={selectedTab() === "servers"}>
              <Button variant="amber" class="mx-4">
                Search Server
              </Button>
              <Button variant="sky" class="mx-4">
                Create Server
              </Button>
              <Index each={servers}>
                {(server, _) => <SidebarServer server={server} />}
              </Index>
            </Match>
            <Match when={selectedTab() === "direct"}>
              <Button variant="amber" class="mx-4">
                Search users
              </Button>
              <Index each={DMs}>{(dm, _) => <SidebarDM dm={dm} />}</Index>
            </Match>
          </Switch>
        </div>
      </aside>
      <main class="flex-1 flex flex-col h-full">
        <AppHeader />
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
