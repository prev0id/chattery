import { Chat } from "~/components/Chat";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { SERVER_CHAT_TARGET } from "~/features/server/constants";
import { useServerContext } from "~/features/server/context";
import { WSChannelType } from "~/lib/ws";

export default function ServerTextTopicPage() {
  const { currentServer, currentTopic } = useServerContext();

  return (
    <>
      <Header icon={currentTopic()?.type}>
        <HeaderItem>{currentServer()?.name}</HeaderItem>
        <HeaderItem>{currentTopic()?.name}</HeaderItem>
      </Header>
      <Chat
        chatID={currentTopic()?.id}
        type={SERVER_CHAT_TARGET}
        channelType={WSChannelType.TextTopic}
      />
    </>
  );
}
