import { useParams, createAsync } from "@solidjs/router";
import { createMemo, Suspense, Show, createSignal, onMount } from "solid-js";
import { fetchServers } from "../lib/api";
import AppHeader from "../components/AppHeader";
import { Chat, ServersType, TopicTypeText } from "../components/Chat";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";
import SidebarServer from "../components/SidebarServer";
import { Loader } from "lucide-solid";

export default function Servers() {
  const params = useParams();

  const serversQuery = createAsync(() => fetchServers());

  const serverData = createMemo(() => {
    const serverId = parseInt(params.serverID, 10);
    const topicId = parseInt(params.topicID, 10);
    const topicType = params.topicType;

    if (!serverId || !topicId || !topicType) return null;

    const server = serversQuery()?.find((s) => s.id === serverId);
    if (!server) return null;

    const topic = server.topics?.find((t) => t.id === topicId);
    if (!topic) return null;

    return {
      serverName: server.name,
      serverID: server.id,
      topicName: topic.name,
      topicType: topic.type,
    };
  });

  return (
    <>
      <Sidebar>
        <Suspense fallback={<LoadingSpinner />}>
          <SidebarServer servers={serversQuery} selectedServer={serverData} />
        </Suspense>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <Show when={serverData()} fallback={<AppHeader />}>
          <AppHeader
            serverName={serverData().serverName}
            topicName={serverData().topicName}
            topicType={serverData().topicType}
          />
        </Show>
        <Show when={params.topicType === TopicTypeText}>
          <Chat chatID={parseInt(params.topicID, 10)} type={ServersType} />
        </Show>
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
      <div class="mx-auto mt-8 flex items-center gap-2">
        <Loader class="size-5 animate-spin" />
        <span class="tracking-wider text-lg font-semibold">
          Loading servers...
        </span>
      </div>
    </Show>
  );
}
