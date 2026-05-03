import { Chat } from "~/components/Chat";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { DM_CHAT_TARGET } from "~/features/dm/constants";
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
        type={DM_CHAT_TARGET}
        channelType={WSChannelType.DM}
      />
    </>
  );
}
