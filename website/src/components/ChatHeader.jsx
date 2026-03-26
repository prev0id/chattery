import { MessagesSquare } from "lucide-solid";

export default function ChatHeader(props) {
  const { topicName } = props;

  return (
    <div class="border-b-3 px-4 py-2 flex items-center gap-4 bg-emerald-50">
      <MessagesSquare size={30} />
      <h1 class="text-2xl font-semibold tracking-wider">{topicName}</h1>
    </div>
  );
}
