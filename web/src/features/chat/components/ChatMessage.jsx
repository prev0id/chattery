import { Show } from "solid-js";
import ProfilePicture from "~/shared/ui/ProfilePicture";
import { userData } from "~/stores/auth";

export default function ChatMessage(props) {
  return (
    <Show when={props.message}>
      <div class="flex gap-4 py-2 max-w-3xl min-w-sm w-full mx-auto">
        <ProfilePicture
          src={props.message?.sender?.avatar}
          class="mt-1 size-10"
        />
        <div class="flex-1">
          <div class="flex gap-2 items-center">
            <div class="font-semibold text-xl">
              {props.message?.sender?.id === userData()?.id
                ? "You"
                : props.message?.sender?.username}
            </div>
            <div class="text-sm">{props.message?.createdAt}</div>
          </div>
          <p>{props.message?.text}</p>
        </div>
      </div>
    </Show>
  );
}
