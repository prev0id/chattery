import { Show } from "solid-js";

export default function VoiceTopicStatus(props) {
  return (
    <Show when={props.status() !== "connected" || props.error()}>
      <div class="mx-4 mt-4 rounded-lg border-2 neo-shadow bg-white px-3 py-2 text-sm font-semibold tracking-wider">
        <span>Status: {props.status()}</span>
        <Show when={props.error()}>
          <span class="ml-3 text-rose-700">{props.error()}</span>
        </Show>
      </div>
    </Show>
  );
}
