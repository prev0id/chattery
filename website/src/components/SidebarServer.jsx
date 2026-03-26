import { Mic, MessagesSquare, Settings2 } from "lucide-solid";

export default function SidebarServer(props) {
  const { name, topics = [], selectedTopicID, onTopicSelect } = props;

  const textTopics = topics.filter((t) => t.type === "text");
  const voiceTopics = topics.filter((t) => t.type === "voice");

  const isServerSelected = topics.some((t) => t.id === selectedTopicID);

  return (
    <details open class="border-2 rounded-lg p-1 bg-white">
      <summary
        class={`px-2 flex justify-between items-center rounded-lg border-2 transition-all duration-300 ease-in-out cursor-pointer ${
          isServerSelected
            ? "bg-emerald-200 border-black"
            : "hover:bg-emerald-200 hover:border-black border-white"
        }`}
      >
        <h2 class="text-lg font-semibold">{name}</h2>
        <Settings2 size={20} />
      </summary>
      <hr class="mt-1" />
      {textTopics.map((topic) => (
        <div
          key={topic.id}
          onClick={() => onTopicSelect?.(topic.id)}
          class={`flex items-center gap-1 px-2 my-1 py-0.5 border-2 rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out cursor-pointer ${
            topic.id === selectedTopicID
              ? "border-black bg-emerald-200"
              : "border-white hover:border-black"
          }`}
        >
          <MessagesSquare size={20} />
          <p>{topic.name}</p>
        </div>
      ))}
      {textTopics.length > 0 && voiceTopics.length > 0 && <hr />}
      {voiceTopics.map((topic) => (
        <div
          key={topic.id}
          onClick={() => onTopicSelect?.(topic.id)}
          class={`flex items-center gap-1 px-2 my-1 py-0.5 border-2 rounded-lg neo-shadow-white hover:neo-shadow transition-all duration-300 ease-in-out cursor-pointer ${
            topic.id === selectedTopicID
              ? "border-black bg-emerald-200"
              : "border-white hover:border-black"
          }`}
        >
          <Mic size={20} />
          <p>{topic.name}</p>
        </div>
      ))}
    </details>
  );
}
