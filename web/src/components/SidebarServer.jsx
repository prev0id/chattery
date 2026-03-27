import { For, createMemo } from "solid-js";
import { Mic, MessagesSquare, Settings2 } from "lucide-solid";
import { selectedTopicID, servers, selectTopic, selectedServerID } from "../stores/app";

export default function SidebarServer(props) {
  const server = createMemo(() => servers.find((server) => server.id === props.serverID));

  const topics = createMemo(() => server().topics);

  const textTopics = createMemo(() => topics().filter((topic) => topic.type === "text"));

  const voiceTopics = createMemo(() => topics().filter((topic) => topic.type === "voice"));

  console.log("topics", topics());
  console.log("text topics", textTopics());
  console.log("voice topics", voiceTopics());

  return (
    <details open class="border-2 rounded-lg p-1 bg-white">
      <summary
        class={`px-2 flex justify-between items-center rounded-lg border-2 transition-all duration-300 ease-in-out cursor-pointer ${
          server().id == selectedServerID()
            ? "bg-emerald-200 border-black"
            : "hover:bg-emerald-200 hover:border-black border-white"
        }`}
      >
        <h2 class="text-lg font-semibold">{server().name}</h2>
        <Settings2 size={20} />
      </summary>
      {textTopics().length > 0 && voiceTopics().length > 0 && <hr class="mt-1" />}
      <For each={textTopics()}>
        {(topic, _) => (
          <div
            onClick={() => selectTopic(topic.id, server().id)}
            class={`flex items-center gap-1 px-2 my-1 py-0.5 border-2 rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out cursor-pointer ${
              topic.id === selectedTopicID()
                ? "border-black bg-emerald-200"
                : "border-white hover:border-black"
            }`}
          >
            <MessagesSquare size={20} />
            <p>{topic.name}</p>
          </div>
        )}
      </For>
      {textTopics().length > 0 && voiceTopics().length > 0 && <hr />}
      <For each={voiceTopics()}>
        {(topic, _) => (
          <div
            onClick={() => selectTopic(topic.id, server().id)}
            class={`flex items-center gap-1 px-2 my-1 py-0.5 border-2 rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out cursor-pointer ${
              topic.id === selectedTopicID()
                ? "border-black bg-emerald-200"
                : "border-white hover:border-black"
            }`}
          >
            <Mic size={20} />
            <p>{topic.name}</p>
          </div>
        )}
      </For>
    </details>
  );
}
