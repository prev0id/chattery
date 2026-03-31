import { Show } from "solid-js";

export default function ProfilePicture(props) {
  return (
    <div class={`relative ${props.class ?? ""}`}>
      <img {...props} class="border-2 neo-shadow rounded-lg w-full h-full" />
      <Show when={props.unread}>
        <div class="absolute bottom-0 right-0 w-3 h-3 bg-rose-600 border-2 neo-shadow-sm rounded-full transform translate-x-0.5 translate-y-0.5"></div>
      </Show>
    </div>
  );
}
