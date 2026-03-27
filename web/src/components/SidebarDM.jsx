import { Show } from "solid-js";
import { selectDM, selectedDM } from "../stores/app";
import ProfilePicture from "./ProfilePicture";

export default function SidebarDM(props) {
  return (
    <button
      onClick={() => selectDM(props.dm())}
      class="flex gap-2 h-14 border-2  border-black rounded-lg neo-shadow-rose p-0.5 hover:neo-shadow hover:bg-emerald-200 hover:scale-105 transition-all duration-300 ease-in-out text-left"
      classList={{
        "bg-emerald-200": props.dm().id === selectedDM()?.id,
        "bg-white": props.dm().id !== selectedDM()?.id,
      }}
    >
      <ProfilePicture src={props.dm().profilePicture} class="size-9 my-auto" />
      <div class="flex-1 overflow-hidden grid grid-rows-2">
        <p class="font-semibold truncate">{props.dm().username}</p>
        <p class="truncate max-w-full text-sm">
          {props.dm().lastMessage?.message ?? ""}
        </p>
      </div>
      <Show when={props.dm().unread > 0}>
        <div class="my-auto rounded-full px-1 py-0.5 leading-none bg-red-600 text-white text-center text-sm">
          {props.dm().unread}
        </div>
      </Show>
    </button>
  );
}
