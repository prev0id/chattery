import { DMs, selectedDM } from "../stores/dm";
import { A } from "@solidjs/router";
import ProfilePicture from "./ProfilePicture";
import Button from "./Button";

export default function SidebarDM() {
  return (
    <>
      <Button variant="amber" class="mx-4">
        Search users
      </Button>
      <Index each={DMs()}>{(dm, _) => <DM dm={dm} />}</Index>
    </>
  );
}

function DM(props) {
  const dmID = () => props.dm().id;
  const username = () => props.dm().user.username;

  return (
    <A
      href={`/dm/${dmID()}`}
      class="flex gap-2 h-14 border-2 border-black rounded-lg neo-shadow-rose p-0.5 hover:neo-shadow hover:bg-emerald-200 hover:scale-105 transition-all duration-300 ease-in-out text-left"
      activeClass="bg-emerald-200"
      inactiveClass="bg-white"
    >
      <ProfilePicture
        src={props.dm().user.avatar}
        unread={props.dm().unread}
        class="size-9 my-auto"
      />
      <div class="flex-1 overflow-hidden grid grid-rows-2">
        <p class="font-semibold truncate">{username()}</p>
        <p class="truncate max-w-full text-sm">
          {props.dm().message?.content ?? ""}
        </p>
      </div>
    </A>
  );
}
