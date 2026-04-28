import { createMemo, createResource, createSignal, For, Show } from "solid-js";
import { revalidate, useNavigate, useSearchParams } from "@solidjs/router";
import { LogOut, Search, Settings2, UserPlus } from "lucide-solid";
import Button from "~/components/Button";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { joinServer, leaveServer, searchServers } from "~/lib/api";
import { GetServers, UseServerContext } from "~/stores/server";

export default function Manage() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const { servers } = UseServerContext();
  const [pendingServers, setPendingServers] = createSignal(new Set());

  const query = createMemo(() => (searchParams.query || "").trim());
  const isSearching = createMemo(() => query().length > 0);
  const [searchResults, { refetch: refetchSearchResults }] = createResource(
    query,
    (value) => (value ? searchServers(value) : []),
  );

  const isPending = (serverID) => pendingServers().has(serverID);

  const setServerPending = (serverID, pending) => {
    setPendingServers((current) => {
      const next = new Set(current);
      if (pending) {
        next.add(serverID);
      } else {
        next.delete(serverID);
      }
      return next;
    });
  };

  const refreshServers = async () => {
    await revalidate(GetServers.key);
  };

  const handleJoin = async (serverID) => {
    setServerPending(serverID, true);
    const result = await joinServer(serverID);
    if (result) {
      await refreshServers();
      refetchSearchResults();
    }
    setServerPending(serverID, false);
  };

  const handleLeave = async (serverID) => {
    if (!confirm("Are you sure you want to leave this server?")) return;

    setServerPending(serverID, true);
    const result = await leaveServer(serverID);
    if (result) {
      await refreshServers();
    }
    setServerPending(serverID, false);
  };

  return (
    <>
      <Header icon="settings">
        <HeaderItem>Manage Servers</HeaderItem>
      </Header>
      <div class="m-auto w-full max-w-2xl min-w-sm rounded-xl mt-4 p-4 overflow-auto flex flex-col gap-8">
        <SearchInput
          query={searchParams.query || ""}
          setSearchParams={setSearchParams}
        />

        <Show
          when={isSearching()}
          fallback={
            <ServerList
              emptyText="You have not joined any servers yet"
              servers={servers() || []}
            >
              {(server) => (
                <MyServerActions
                  server={server}
                  isPending={isPending(server.id)}
                  onEdit={() => navigate(`/server/${server.id}/edit`)}
                  onLeave={() => handleLeave(server.id)}
                />
              )}
            </ServerList>
          }
        >
          <ServerList
            emptyText="No servers found"
            loading={searchResults.loading}
            servers={searchResults() || []}
          >
            {(server) => (
              <Button
                variant="emerald"
                disabled={isPending(server.id)}
                onClick={() => handleJoin(server.id)}
              >
                <span class="flex items-center gap-1">
                  <UserPlus class="size-5" />
                  Join
                </span>
              </Button>
            )}
          </ServerList>
        </Show>
      </div>
    </>
  );
}

function SearchInput(props) {
  return (
    <div class="relative">
      <Search class="absolute left-3 top-1/2 size-5 -translate-y-1/2" />
      <input
        class="px-10 py-1 border-3 neo-shadow-lg w-full text-lg rounded-lg focus:outline-none focus:border-sky-500"
        type="search"
        name="search-server"
        placeholder="Search servers"
        onInput={(event) => {
          props.setSearchParams({ query: event.currentTarget.value || null });
        }}
        value={props.query}
      />
    </div>
  );
}

function ServerList(props) {
  return (
    <Show
      when={!props.loading}
      fallback={<p class="text-center text-lg font-semibold">Loading...</p>}
    >
      <Show
        when={props.servers.length > 0}
        fallback={
          <p class="text-center text-lg font-semibold">{props.emptyText}</p>
        }
      >
        <div class="flex flex-col gap-4">
          <For each={props.servers}>
            {(server) => (
              <ServerPreview server={server}>
                {props.children(server)}
              </ServerPreview>
            )}
          </For>
        </div>
      </Show>
    </Show>
  );
}

function ServerPreview(props) {
  return (
    <div class="flex items-center justify-between gap-3 px-3 py-2 border-2 neo-shadow rounded-lg bg-white">
      <p class="min-w-0 flex-1 truncate font-semibold text-xl tracking-wider">
        {props.server?.name}
      </p>
      <div class="shrink-0">{props.children}</div>
    </div>
  );
}

function MyServerActions(props) {
  return (
    <Show
      when={props.server.role === "owner"}
      fallback={
        <Button
          variant="rose"
          disabled={props.isPending}
          onClick={props.onLeave}
        >
          <span class="flex items-center gap-1">
            <LogOut class="size-5" />
            Leave
          </span>
        </Button>
      }
    >
      <Button variant="sky" onClick={props.onEdit}>
        <span class="flex items-center gap-1">
          <Settings2 class="size-5" />
          Edit
        </span>
      </Button>
    </Show>
  );
}
