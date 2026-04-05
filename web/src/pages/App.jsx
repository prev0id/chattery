import { Match, Switch } from "solid-js";
import AppHeader from "../components/AppHeader";
import SidebarDM from "../components/SidebarDM";
import SidebarServer from "../components/SidebarServer";
import Chat from "../components/Chat";
import { selectedTab } from "../stores/app";
import { selectedTopic } from "../stores/server";
import { selectedDM } from "../stores/dm";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";

export default function App() {
  return (
    <>
      <Sidebar>
        <Switch>
          <Match when={selectedTab() === "servers"}>
            <SidebarServer />
          </Match>
          <Match when={selectedTab() === "direct"}>
            <SidebarDM />
          </Match>
        </Switch>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <AppHeader />
        <Switch>
          <Match when={selectedTopic()?.type === "text"}>
            <Chat chatID={selectedTopic()?.id} />
          </Match>
          <Match when={selectedDM()}>
            <Chat chatID={selectedDM().id} />
          </Match>
          <Match when={selectedTopic()?.type === "voice"}>
            <p>Voice topic {selectedTopic().id} </p>
          </Match>
        </Switch>
      </main>
      <Toasts />
    </>
  );
}
