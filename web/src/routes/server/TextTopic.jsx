import { Chat, ServersType } from "~/components/Chat";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { WSChannelType } from "~/lib/ws";
import { UseServerContext } from "~/stores/server";

export default function TextTopic() {
  const { currentServer, currentTopic } = UseServerContext();

  return (
    <>
      <Header icon={currentTopic()?.type}>
        <HeaderItem>{currentServer()?.name}</HeaderItem>
        <HeaderItem>{currentTopic()?.name}</HeaderItem>
      </Header>
      <Chat
        chatID={currentTopic()?.id}
        type={ServersType}
        channelType={WSChannelType.TextTopic}
      />
    </>
  );
}
