import Toasts from "~/components/Toast";
import Sidebar from "~/components/Sidebar";
import SidebarServer from "~/components/SidebarServer";
import Button from "~/components/Button";
import { GetServers, ServerContext } from "~/stores/server";
import { createAsync, useParams } from "@solidjs/router";
import { createMemo } from "solid-js";

export default function Wrapper(props) {
  const params = useParams();
  const servers = createAsync(() => GetServers());

  console.log(servers());

  const currentServerID = () => parseInt(params.serverID, 10);
  const currentTopicID = () => parseInt(params.topicID, 10);

  const currentServer = createMemo(() =>
    servers()?.find((server) => server.id === currentServerID()),
  );
  const currentTopic = createMemo(() =>
    currentServer()?.topics.find((topic) => topic.id === currentTopicID()),
  );

  return (
    <>
      <Sidebar>
        <Button variant="amber" class="mx-4">
          Search Server
        </Button>
        <Button variant="sky" class="mx-4" popovertarget="create-server-modal">
          Create Server
        </Button>
        <Index each={servers()}>
          {(server) => (
            <SidebarServer
              server={server}
              selectedServerID={currentServerID()}
            />
          )}
        </Index>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <ServerContext.Provider
          value={{
            currentServer,
            currentServerID,
            currentTopic,
            currentTopicID,
          }}
        >
          {props.children}
        </ServerContext.Provider>
      </main>
      <Toasts />
    </>
  );
}
