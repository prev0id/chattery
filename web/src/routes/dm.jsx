import { useParams } from "@solidjs/router";
import { createMemo } from "solid-js";
import AppHeader from "../components/AppHeader";
import SidebarDM from "../components/SidebarDM";
import { Chat, DMsType } from "../components/Chat";
import { DMs as storedDMs } from "../stores/dm";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";

export default function DMs() {
  const params = useParams();

  const dmUsername = createMemo(() => {
    const dmId = parseInt(params.dmID, 10);
    if (!dmId || !storedDMs()) return null;
    const dm = storedDMs().find((d) => d.id === dmId);
    return dm?.user?.username || null;
  });

  return (
    <>
      <Sidebar>
        <SidebarDM />
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <AppHeader dmUsername={dmUsername} />
        <Chat chatID={parseInt(params.dmID, 10)} type={DMsType} />
      </main>
      <Toasts />
    </>
  );
}
