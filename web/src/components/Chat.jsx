import { messages } from "../stores/app";
import ChatInput from "./ChatInput";
import ChatMessage from "./ChatMessage";

export default function Chat(props) {
  return (
    <>
      <div>This is chat {props.chatID}</div>
      <div class="max-w-5xl h-full mx-auto p-4 flex flex-col overflow-auto">
        <For each={messages} fallback={<div>No messages yet.</div>}>
          {(message, _) => <ChatMessage msg={message} />}
        </For>
      </div>
      <ChatInput onSend={(text) => console.log(text)}></ChatInput>
    </>
  );
}
