import { Index, createMemo } from "solid-js";
import { Mic, MessagesSquare, Settings2 } from "lucide-solid";
import {
  selectedTopic,
  selectTopic,
  selectedServer,
  leaveTopic,
} from "../stores/app";

export default function SidebarServer(props) {
  const topics = createMemo(() => props.server().topics);

  const textTopics = createMemo(() =>
    topics().filter((topic) => topic.type === "text"),
  );

  const voiceTopics = createMemo(() =>
    topics().filter((topic) => topic.type === "voice"),
  );

  const toggleVoiceTopic = (topic) => {
    if (selectedTopic()?.id === topic.id) {
      leaveTopic();
      return;
    }
    selectTopic(topic, props.server());
  };

  return (
    <details open class="border-2 rounded-lg p-1 bg-white">
      <summary
        class="px-2 flex justify-between items-center rounded-lg border-2 transition-all duration-300 ease-in-out cursor-pointer hover:bg-emerald-200 hover:border-black border-white focus:outline-none focus:border-black"
        classList={{
          "bg-emerald-200": props.server().id == selectedServer()?.id,
        }}
      >
        <h2 class="text-lg font-semibold">{props.server().name}</h2>
        <Settings2 size={20} />
      </summary>
      {textTopics().length > 0 && voiceTopics().length > 0 && (
        <hr class="mt-1" />
      )}
      <Index each={textTopics()}>
        {(topic, _) => (
          <SidebarTopic
            topic={topic()}
            onClick={() => selectTopic(topic(), props.server())}
          >
            <MessagesSquare class="size-5" />
          </SidebarTopic>
        )}
      </Index>
      {textTopics().length > 0 && voiceTopics().length > 0 && <hr />}
      <Index each={voiceTopics()}>
        {(topic, _) => (
          <SidebarTopic
            topic={topic()}
            onClick={() => toggleVoiceTopic(topic())}
          >
            <Mic class="size-5" />
          </SidebarTopic>
        )}
      </Index>
    </details>
  );
}

function SidebarTopic(props) {
  return (
    <div
      onClick={() => props.onClick()}
      class="flex items-center gap-1 px-2 my-1 py-0.5 border-2 rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out cursor-pointer border-white hover:border-black"
      classList={{
        "bg-emerald-200": props.topic.id === selectedTopic()?.id,
      }}
    >
      {props.children}
      <p>{props.topic.name}</p>
    </div>
  );
}
