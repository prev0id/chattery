import { createAsync, useNavigate, useParams } from "@solidjs/router";
import { createMemo, createSignal, For, Show } from "solid-js";
import Button from "~/shared/ui/Button";
import ServerSidebarItem from "~/features/server/components/ServerSidebarItem";
import { getServersQuery } from "~/features/server/actions";
import { ServerContext } from "~/features/server/context";
import { getUserErrorMessage } from "~/shared/api/errors";
import { routes } from "~/shared/config/routes";
import { parseRouteId } from "~/shared/lib/route";
import AppSidebar from "~/shared/ui/AppSidebar";
import { SERVER_MESSAGES } from "~/features/server/constants";

export default function ServerLayout(props) {
  const navigate = useNavigate();
  const params = useParams();
  const [loadError, setLoadError] = createSignal("");
  const loadServers = async () => {
    setLoadError("");
    try {
      return await getServersQuery();
    } catch (error) {
      setLoadError(getUserErrorMessage(error, SERVER_MESSAGES.listFailed));
      return [];
    }
  };
  const servers = createAsync(loadServers);

  const currentServerId = () => parseRouteId(params.serverId);
  const currentTopicId = () => parseRouteId(params.topicId);

  const currentServer = createMemo(() =>
    servers()?.find((server) => server.id === currentServerId()),
  );
  const currentTopic = createMemo(() =>
    currentServer()?.topics.find((topic) => topic.id === currentTopicId()),
  );

  return (
    <>
      <AppSidebar fallback="Loading servers...">
        <Button
          variant="amber"
          class="mx-4"
          onClick={() => navigate(routes.server.manage())}
        >
          Manage Servers
        </Button>
        <Button
          variant="sky"
          class="mx-4"
          onClick={() => navigate(routes.server.create())}
        >
          Create Server
        </Button>
        <Show when={loadError()}>
          <p class="rounded-lg bg-red-200 px-3 py-2 text-sm font-semibold text-red-700">
            {loadError()}
          </p>
        </Show>
        <For each={servers()}>
          {(server) => (
            <ServerSidebarItem
              server={server}
              selectedServerId={currentServerId()}
            />
          )}
        </For>
      </AppSidebar>
      <main class="flex-1 flex flex-col h-full">
        <ServerContext.Provider
          value={{
            currentServer,
            currentTopic,
            servers,
          }}
        >
          {props.children}
        </ServerContext.Provider>
      </main>
    </>
  );
}
