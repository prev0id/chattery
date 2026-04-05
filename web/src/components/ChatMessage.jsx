import { Show } from "solid-js";
import { userData } from "../stores/auth";
import ProfilePicture from "./ProfilePicture";

export default function ChatMessage(props) {
  return (
    <Show when={props.msg}>
      <div class="flex gap-4 py-2">
        <ProfilePicture src={props.msg?.sender?.avatar} class="mt-1 size-10" />
        <div class="flex-1">
          <div class="flex gap-2 items-center">
            <div class="font-semibold text-xl">
              {props.msg?.sender?.id === userData()?.id
                ? "You"
                : props.msg?.sender?.username}
            </div>
            <div class="text-sm">{props.msg?.created_at}</div>
          </div>
          <p>{props.msg?.text}</p>
        </div>
      </div>
    </Show>
  );
}
