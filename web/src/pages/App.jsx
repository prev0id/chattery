import { Index, Match, Switch } from "solid-js";
import AppHeader from "../components/AppHeader";
import Button from "../components/Button";
import ProfilePicture from "../components/ProfilePicture";
import SidebarDM from "../components/SidebarDM";
import SidebarServer from "../components/SidebarServer";
import Chat from "../components/Chat";
import {
  servers,
  selectedTab,
  DMs,
  changeTab,
  selectedTopic,
  selectedDM,
  userData,
} from "../stores/app";
import { ProfileSettingsModal } from "../components/ModalProfileSettings";
import Toasts from "../components/Toast";

export default function App() {
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
          <button
            class="mt-auto hover:scale-105 transition-all duration-300 ease-in-out"
            popovertarget="profile-settings-popover"
          >
            <ProfilePicture src={userData()?.avatar} />
          </button>
          <ProfileSettingsModal id="profile-settings-popover" />
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
