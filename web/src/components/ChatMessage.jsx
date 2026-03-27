import { userData } from "../stores/app";
import ProfilePicture from "./ProfilePicture";

export default function ChatMessage(props) {
  return (
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
  );
}
