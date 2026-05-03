import { createSignal } from "solid-js";
import Button from "~/components/Button";

export default function ChatInput(props) {
  const { onSend } = props;
  const [message, setMessage] = createSignal("");

  const handleSubmit = async (event) => {
    event.preventDefault();
    const text = message().trim();
    if (!text || props.disabled) return;

    if ((await onSend?.(text)) !== false) {
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
        onInput={(event) => setMessage(event.currentTarget.value)}
        disabled={props.disabled}
        class="max-w-2xl min-w-sm w-full p-2 neo-shadow border-2 rounded-lg field-sizing-content focus:outline-none focus:border-sky-500 resize-none"
        placeholder="Write your message"
        rows="1"
      />
      <Button type="submit" variant="emerald" disabled={props.disabled}>
        Send
      </Button>
    </form>
  );
}
