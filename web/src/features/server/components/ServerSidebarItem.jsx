import { A } from "@solidjs/router";
import { createMemo, For, Show } from "solid-js";
import { Mic, MessagesSquare } from "lucide-solid";
import { SERVER_TOPIC_TYPE } from "~/features/server/constants";
import { routes } from "~/shared/config/routes";

/**
 * @param {{server: import("~/features/server/model").Server, selectedServerId: number}} props
 */
export default function ServerSidebarItem(props) {
  const topics = createMemo(() => props.server.topics || []);

  const textTopics = createMemo(() =>
    topics().filter((topic) => topic.type === SERVER_TOPIC_TYPE.text),
  );

  const voiceTopics = createMemo(() =>
    topics().filter((topic) => topic.type === SERVER_TOPIC_TYPE.voice),
  );

  return (
    <details open class="border-2 rounded-lg p-1 bg-white">
      <summary
        class="px-2 flex justify-between items-center rounded-lg border-2 transition-all duration-300 ease-in-out cursor-pointer hover:bg-emerald-200 hover:border-black border-white focus:outline-none focus:border-black"
        classList={{
          "bg-emerald-200": props.server.id === props.selectedServerId,
        }}
      >
        <h2 class="text-lg font-semibold">{props.server.name}</h2>
      </summary>
      <Show when={textTopics().length > 0 && voiceTopics().length > 0}>
        <hr class="mt-1" />
      </Show>
      <For each={textTopics()}>
        {(topic) => (
          <SidebarTopic href={routes.server.textTopic(props.server.id, topic.id)}>
            <MessagesSquare class="size-5" />
            <p>{topic.name}</p>
          </SidebarTopic>
        )}
      </For>
      <Show when={textTopics().length > 0 && voiceTopics().length > 0}>
        <hr />
      </Show>
      <For each={voiceTopics()}>
        {(topic) => (
          <SidebarTopic
            href={routes.server.voiceTopic(props.server.id, topic.id)}
          >
            <Mic class="size-5" />
            <p>{topic.name}</p>
          </SidebarTopic>
        )}
      </For>
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
