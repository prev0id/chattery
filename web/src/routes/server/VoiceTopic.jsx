import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { UseServerContext } from "~/stores/server";

export default function VoiceTopic() {
  const { currentServer, currentTopic } = UseServerContext();

  return (
    <>
      <Header icon={currentTopic()?.type}>
        <HeaderItem>{currentServer()?.name}</HeaderItem>
        <HeaderItem>{currentTopic()?.name}</HeaderItem>
      </Header>
      <p>Voice topic {currentTopic()?.id}</p>
      <p>TODO</p>
    </>
  );
}
