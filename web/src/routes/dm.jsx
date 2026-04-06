import { useParams, createAsync } from "@solidjs/router";
import { createMemo, Show, Suspense, createSignal, onMount } from "solid-js";
import { fetchDMs } from "../lib/api";
import Header from "../components/Header";
import HeaderItem from "../components/HeaderItem";
import SidebarDM from "../components/SidebarDM";
import { Chat, DMsType } from "../components/Chat";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";
import { Loader } from "lucide-solid";

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
      <Sidebar>
        <Suspense fallback={<LoadingSpinner />}>
          <SidebarDM dms={dmsQuery} />
        </Suspense>
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

function LoadingSpinner() {
  const [show, setShow] = createSignal(false);

  onMount(() => {
    const timer = setTimeout(() => setShow(true), 300);
    return () => clearTimeout(timer);
  });

  return (
    <Show when={show()}>
      return (
      <div class="mx-auto mt-8 flex items-center gap-2">
        <Loader class="size-5 animate-spin" />
        <span class="tracking-wider text-lg font-semibold">Loading DMs...</span>
      </div>
    </Show>
  );
}
