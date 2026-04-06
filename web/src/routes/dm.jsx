import { useParams, createAsync } from "@solidjs/router";
import { createMemo, Show, Suspense, createSignal, onMount } from "solid-js";
import { fetchDMs } from "../lib/api";
import Header from "../components/Header";
import HeaderItem from "../components/HeaderItem";
import SidebarDM from "../components/SidebarDM";
import { Chat, DMsType } from "../components/Chat";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";
import Button from "../components/Button";

export default function DMs() {
  const params = useParams();

  const dmsQuery = createAsync(() => fetchDMs());

  const dmUsername = createMemo(() => {
    const dmId = parseInt(params.dmID, 10);
    if (!dmId || !dmsQuery()) return null;
    const dm = dmsQuery().find((d) => d.id === dmId);
    return dm?.user?.username || null;
  });

  return (
    <>
      <Sidebar fallback="Loading DMs...">
        <Button variant="amber" class="mx-4">
          Search users
        </Button>
        <For each={dmsQuery()}>{(dm) => <SidebarDM dm={dm} />}</For>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <Show when={dmUsername()} fallback={<Header />}>
          <Header icon="text">
            <HeaderItem>{dmUsername()}</HeaderItem>
          </Header>
        </Show>
        <Chat chatID={parseInt(params.dmID, 10)} type={DMsType} />
      </main>
      <Toasts />
    </>
  );
}
