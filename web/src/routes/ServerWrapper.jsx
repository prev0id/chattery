import { useParams, createAsync } from "@solidjs/router";
import { createMemo, Suspense, Show, createSignal, onMount } from "solid-js";
import { fetchServers } from "../lib/api";

import Toasts from "../components/Toast";
import Sidebar from "../components/Sidebar";
import SidebarServer from "../components/SidebarServer";
import { Loader } from "lucide-solid";
import Button from "../components/Button";

export default function ServerPage(props) {
  const params = useParams();

  const serversQuery = createAsync(() => fetchServers());

  const state = createMemo(() => {
    const serverId = parseInt(params.serverID, 10);
    const topicId = parseInt(params.topicID, 10);

    if (!serverId && !topicId) return {};

    const server = serversQuery()?.find((s) => s.id === serverId);
    if (!server) return {};

    const topic = server.topics?.find((t) => t.id === topicId);

    return {
      server: server,
      topic: topic,
    };
  });

  return (
    <>
      <Sidebar>
        <Suspense fallback={<LoadingSpinner />}>
          <Button variant="amber" class="mx-4">
            Search Server
          </Button>
          <Button
            variant="sky"
            class="mx-4"
            popovertarget="create-server-modal"
          >
            Create Server
          </Button>
          <Index each={serversQuery()}>
            {(server) => (
              <SidebarServer
                server={server}
                selectedServerID={state().server?.id}
              />
            )}
          </Index>
        </Suspense>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">{props.children}</main>
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
