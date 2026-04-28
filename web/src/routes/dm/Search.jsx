import { createMemo, createResource, createSignal, For, Show } from "solid-js";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { MessageSquarePlus, Search } from "lucide-solid";
import Button from "~/components/Button";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import ProfilePicture from "~/components/ProfilePicture";
import { createDM, searchUsers } from "~/lib/api";
import { UseDMContext } from "~/stores/dm";

export default function SearchDM() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const { refetchDMs } = UseDMContext();
  const [pendingUsers, setPendingUsers] = createSignal(new Set());

  const query = createMemo(() => (searchParams.query || "").trim());
  const [users, { refetch: refetchUsers }] = createResource(query, (value) =>
    value ? searchUsers(value) : [],
  );

  const isPending = (userID) => pendingUsers().has(userID);

  const setUserPending = (userID, pending) => {
    setPendingUsers((current) => {
      const next = new Set(current);
      if (pending) {
        next.add(userID);
      } else {
        next.delete(userID);
      }
      return next;
    });
  };

  const handleCreateDM = async (userID) => {
    setUserPending(userID, true);
    const result = await createDM(userID);
    if (result?.id) {
      await refetchDMs?.();
      refetchUsers();
      navigate(`/dm/${result.id}`);
      return;
    }
    setUserPending(userID, false);
  };

  return (
    <>
      <Header icon="settings">
        <HeaderItem>Search Users</HeaderItem>
      </Header>
      <div class="m-auto w-full max-w-2xl min-w-sm rounded-xl mt-4 p-4 overflow-auto flex flex-col gap-8">
        <SearchInput
          query={searchParams.query || ""}
          setSearchParams={setSearchParams}
        />

        <Show
          when={query().length > 0}
          fallback={
            <p class="text-center text-lg font-semibold">
              Type a username to start searching
            </p>
          }
        >
          <UserList loading={users.loading} users={users() || []}>
            {(user) => (
              <Button
                variant="emerald"
                disabled={isPending(user.id)}
                onClick={() => handleCreateDM(user.id)}
              >
                <span class="flex items-center gap-1">
                  <MessageSquarePlus class="size-5" />
                  Create
                </span>
              </Button>
            )}
          </UserList>
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
        name="search-user"
        placeholder="Search users"
        onInput={(event) => {
          props.setSearchParams({ query: event.currentTarget.value || null });
        }}
        value={props.query}
      />
    </div>
  );
}

function UserList(props) {
  return (
    <Show
      when={!props.loading}
      fallback={<p class="text-center text-lg font-semibold">Loading...</p>}
    >
      <Show
        when={props.users.length > 0}
        fallback={
          <p class="text-center text-lg font-semibold">No users found</p>
        }
      >
        <div class="flex flex-col gap-4">
          <For each={props.users}>
            {(user) => (
              <UserPreview user={user}>{props.children(user)}</UserPreview>
            )}
          </For>
        </div>
      </Show>
    </Show>
  );
}

function UserPreview(props) {
  return (
    <div class="flex items-center justify-between gap-3 px-3 py-2 border-2 neo-shadow rounded-lg bg-white">
      <div class="min-w-0 flex flex-1 items-center gap-2">
        <ProfilePicture src={props.user.avatar} class="size-10 shrink-0" />
        <p class="truncate font-semibold text-xl tracking-wider">
          {props.user.username}
        </p>
      </div>
      <div class="shrink-0">{props.children}</div>
    </div>
  );
}
