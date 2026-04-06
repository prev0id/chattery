import { Index, Show, createMemo } from "solid-js";
import { A } from "@solidjs/router";
import { Mic, MessagesSquare, Settings2 } from "lucide-solid";
import { TopicTypeText, TopicTypeVoice } from "./Chat";

export default function SidebarServer(props) {
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
            href={`/server/${props.server().id}/${topic().type}/${topic().id}`}
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
            href={`/server/${props.server().id}/${topic().type}/${topic().id}`}
          >
            <Mic class="size-5" />
            <p>{topic().name}</p>
          </SidebarTopic>
        )}
      </Index>
      <Show when={props.server().role === "owner"}>
        <hr class="my-1" />
        <SidebarTopic href={`/server/${props.server().id}/edit`}>
          <Settings2 class="size-5" />
          <p>Edit</p>
        </SidebarTopic>
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
