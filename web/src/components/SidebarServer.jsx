import { Index, Show, createMemo } from "solid-js";
import { A } from "@solidjs/router";
import { Mic, MessagesSquare } from "lucide-solid";
import Button from "./Button";
import CreateServerModal from "./CreateServerModal";
import EditServerModal from "./EditServerModal";
import { setSelectedServerForEdit } from "../stores/serverState";
import { TopicTypeText, TopicTypeVoice } from "./Chat";

export default function SidebarServer(props) {
  return (
    <>
      <Button variant="amber" class="mx-4">
        Search Server
      </Button>
      <Button variant="sky" class="mx-4" popovertarget="create-server-modal">
        Create Server
      </Button>
      <Index each={props.servers()}>
        {(server) => (
          <Server
            server={server}
            selectedServerID={props.selectedServer()?.id}
          />
        )}
      </Index>
      <CreateServerModal id="create-server-modal" />
      <EditServerModal id="edit-server-modal" />
    </>
  );
}

function Server(props) {
  const topics = createMemo(() => props.server().topics || []);

  const textTopics = createMemo(() =>
    topics().filter((topic) => topic.type === TopicTypeText),
  );

  const voiceTopics = createMemo(() =>
    topics().filter((topic) => topic.type === TopicTypeVoice),
  );

  return (
    <details open class="border-2 rounded-lg p-1 bg-white">
      <summary
        class="px-2 flex justify-between items-center rounded-lg border-2 transition-all duration-300 ease-in-out cursor-pointer hover:bg-emerald-200 hover:border-black border-white focus:outline-none focus:border-black"
        classList={{
          "bg-emerald-200": props.server().id === props.selectedServerID,
        }}
      >
        <h2 class="text-lg font-semibold">{props.server().name}</h2>
      </summary>
      {textTopics().length > 0 && voiceTopics().length > 0 && (
        <hr class="mt-1" />
      )}
      <Index each={textTopics()}>
        {(topic) => (
          <SidebarTopic
            href={`/server/${props.server().id}/${TopicTypeText}/${topic().id}`}
          >
            <MessagesSquare class="size-5" />
            <p>{topic().name}</p>
          </SidebarTopic>
        )}
      </Index>
      {textTopics().length > 0 && voiceTopics().length > 0 && <hr />}
      <Index each={voiceTopics()}>
        {(topic) => (
          <SidebarTopic
            href={`/server/${props.server().id}/${TopicTypeVoice}/${topic().id}`}
          >
            <Mic class="size-5" />
            <p>{topic().name}</p>
          </SidebarTopic>
        )}
      </Index>
      <Show when={props.server().role === "owner"}>
        <hr class="my-1" />
        <Button
          smallText
          variant="emerald"
          popovertarget="edit-server-modal"
          onClick={() => {
            setSelectedServerForEdit(props.server());
          }}
        >
          Update
        </Button>
      </Show>
    </details>
  );
}

function SidebarTopic(props) {
  return (
    <A
      href={props.href}
      class="flex items-center gap-1 px-2 my-1 py-0.5 border-2 rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out cursor-pointer border-white hover:border-black"
      activeClass="bg-emerald-200"
    >
      {props.children}
    </A>
  );
}
