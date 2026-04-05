import { createSignal } from "solid-js";
import Button from "./Button";

export default function ChatInput(props) {
  const { onSend } = props;
  const [message, setMessage] = createSignal("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (message().trim()) {
      onSend?.(message());
      setMessage("");
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      class="sticky bottom-0 m-4 flex items-center justify-center gap-4"
    >
      <textarea
        value={message()}
        onInput={(e) => setMessage(e.currentTarget.value)}
        class="max-w-2xl min-w-sm w-full p-2 neo-shadow border-2 rounded-lg field-sizing-content focus:outline-none focus:border-sky-500 resize-none"
        placeholder="Write your message"
        rows="1"
      />
      <Button type="submit" variant="emerald">
        Send
      </Button>
    </form>
  );
}
