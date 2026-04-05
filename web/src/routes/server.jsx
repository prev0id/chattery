import { useParams } from "@solidjs/router";
import AppHeader from "../components/AppHeader";
import { Chat, ServersType, TopicTypeText } from "../components/Chat";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";
import SidebarServer from "../components/SidebarServer";
import { Match, Switch, createMemo, Show } from "solid-js";
import { servers } from "../stores/server";

export default function Servers() {
  const params = useParams();

  const serverData = createMemo(() => {
    const serverId = parseInt(params.serverID, 10);
    const topicId = parseInt(params.topicID, 10);
    const topicType = params.topicType;

    if (!serverId || !topicId || !topicType) return null;

    const server = servers()?.find((s) => s.id === serverId);
    if (!server) return null;

    const topic = server.topics?.find((t) => t.id === topicId);
    if (!topic) return null;

    return {
      serverName: server.name,
      topicName: topic.name,
      topicType: topic.type,
    };
  });

  return (
    <>
      <Sidebar>
        <SidebarServer />
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <Show when={serverData()} fallback={<AppHeader />}>
          <AppHeader
            serverName={serverData().serverName}
            topicName={serverData().topicName}
            topicType={serverData().topicType}
          />
        </Show>
        <Switch>
          <Match when={params.topicType === TopicTypeText}>
            <Chat chatID={parseInt(params.topicID, 10)} type={ServersType} />
          </Match>
        </Switch>
      </main>
      <Toasts />
    </>
  );
}
