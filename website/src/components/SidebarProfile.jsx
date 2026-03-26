import ProfilePicture from "./ProfilePicture";

export default function SidebarProfile(props) {
  const {
    avatar,
    name,
    lastMessage,
    unread = 0,
    selected = false,
    onClick,
  } = props;

  return (
    <div
      onClick={onClick}
      class={`flex gap-2 border-2  border-black rounded-lg neo-shadow-rose px-0.5 hover:neo-shadow hover:bg-emerald-200 hover:scale-105 transition-all duration-300 ease-in-out cursor-pointer ${
        selected ? "bg-emerald-200" : "bg-white"
      }`}
    >
      <ProfilePicture src={avatar} class="size-9 my-auto" />
      <div class="flex-1 overflow-hidden">
        <div class="font-semibold truncate">{name}</div>
        <p class="truncate max-w-full text-sm">{lastMessage}</p>
      </div>
      {unread > 0 && (
        <div class="my-auto rounded-full px-1 py-0.5 leading-none bg-red-600 text-white text-center text-sm">
          {unread}
        </div>
      )}
    </div>
  );
}
