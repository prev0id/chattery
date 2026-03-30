import { Show } from "solid-js";
import { selectDM, selectedDM } from "../stores/app";
import ProfilePicture from "./ProfilePicture";
import { Info } from "lucide-solid";

export default function SidebarDM(props) {
  return (
    <button
      onClick={() => selectDM(props.dm())}
      class="flex gap-2 h-14 border-2  border-black rounded-lg neo-shadow-rose p-0.5 hover:neo-shadow hover:bg-emerald-200 hover:scale-105 transition-all duration-300 ease-in-out text-left"
      classList={{
        "bg-emerald-200": props.dm().user.id === selectedDM()?.user.id,
        "bg-white": props.dm().user.id !== selectedDM()?.user.id,
      }}
    >
      <ProfilePicture
        src={props.dm().user.profilePicture}
        class="size-9 my-auto"
      />
      <div class="flex-1 overflow-hidden grid grid-rows-2">
        <p class="font-semibold truncate">{props.dm().user.username}</p>
        <p class="truncate max-w-full text-sm">
          {props.dm().message?.content ?? ""}
        </p>
      </div>
      <Show when={props.dm().unread > 0}>
        <Info
          stroke-width="3"
          class="my-auto h-4 leading-none text-amber-600"
        />
      </Show>
    </button>
  );
}
