import ProfilePicture from "./ProfilePicture";

export default function ChatMessage(props) {
  const { avatar, author, time, content, isOwn = false } = props;

  return (
    <div class="flex gap-4">
      <ProfilePicture src={avatar} class="mt-1 size-10" />
      <div class="flex-1">
        <div class="flex gap-2 items-center">
          <div class="font-semibold text-xl">{isOwn ? "You" : author}</div>
          <div class="text-sm">{time}</div>
        </div>
        <p>{content}</p>
      </div>
    </div>
  );
}
