import { createSignal, Show } from "solid-js";
import { userData } from "../stores/app";
import ProfilePicture from "./ProfilePicture";

export default function ChatMessage(props) {
  const [showLastSeenMessage, setShowLastSeenMessage] = createSignal(
    props.msg.lastSeenMessage ?? false,
  );

  if (showLastSeenMessage()) {
    setTimeout(() => {
      setShowLastSeenMessage(false);
    }, 10000);
  }

  return (
    <>
      <div class="flex gap-4">
        <ProfilePicture
          src={props.msg.user.profilePicture}
          class="mt-1 size-10"
        />
        <div class="flex-1">
          <div class="flex gap-2 items-center">
            <div class="font-semibold text-xl">
              {props.msg.user.id === userData().id
                ? "You"
                : props.msg.user.username}
            </div>
            <div class="text-sm">{props.msg.message.date}</div>
          </div>
          <p>{props.msg.message.content}</p>
        </div>
      </div>
      <div class="flex items-center h-6">
        <Show when={showLastSeenMessage()}>
          <span class="px-4 text-sm text-sky-700">New messages</span>
          <div class="flex-1 border-t border-2 border-sky-700"></div>
        </Show>
      </div>
    </>
  );
}
