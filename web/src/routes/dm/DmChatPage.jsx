import { Chat } from "~/features/chat/components/Chat";
import Header from "~/shared/ui/Header";
import HeaderItem from "~/shared/ui/HeaderItem";
import { CHAT_TARGET } from "~/features/chat/constants";
import { useDmContext } from "~/features/dm/context";
import { WSChannelType } from "~/lib/ws";

export default function DmChatPage() {
  const { currentDm, currentDmId } = useDmContext();

  return (
    <>
      <Header icon="text">
        <HeaderItem>{currentDm()?.user?.username}</HeaderItem>
      </Header>
      <Chat
        chatID={currentDmId()}
        type={CHAT_TARGET.dm}
        channelType={WSChannelType.DM}
      />
    </>
  );
}
