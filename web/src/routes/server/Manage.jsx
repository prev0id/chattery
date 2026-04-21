import Button from "~/components/Button";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { useSearchParams } from "@solidjs/router";
import { createSignal, For } from "solid-js";
import { UseServerContext } from "~/stores/server";

export default function Search() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { servers } = UseServerContext();

  return (
    <>
      <Header icon="settings">
        <HeaderItem>Find Server to Join!</HeaderItem>
      </Header>
      <div class="m-auto w-full max-w-2xl min-w-sm rounded-xl mt-4 p-4 overflow-auto flex flex-col gap-2">
        <SearchInput
          searchParams={searchParams}
          setSearchParams={setSearchParams}
        />
        <p>{searchParams.query}</p>
        <For each={servers()}>
          {(server) => <ServerPreview server={server} />}
        </For>
      </div>
    </>
  );
}

function SearchInput(props) {
  return (
    <>
      <input
        class="px-3 py-1 border-3 neo-shadow-lg w-full text-lg rounded-lg focus:outline-none focus:border-sky-500 flex-1"
        type="text"
        name="search-server"
        onChange={(event) => {
          event.preventDefault();
          props.setSearchParams({ query: event.target.value });
        }}
        value={props.searchParams.query || ""}
      />
    </>
  );
}

function ServerPreview(props) {
  const [joined, setJoined] = createSignal(false);
  return (
    <div class="flex justify-between items-center px-2 py-1 border-2 neo-shadow rounded-lg">
      <p class="font-semibold text-xl tracking-wider">{props.server?.name}</p>
      <Button
        variant={joined() ? "amber" : "emerald"}
        onClick={() => setJoined(!joined())}
      >
        {joined() ? "Leave" : "Join"}
      </Button>
    </div>
  );
}
