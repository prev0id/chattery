import { Show } from "solid-js";
import { Chat } from "../components/Chat";
import Header from "../components/Header";
import HeaderItem from "../components/HeaderItem";

export default function TextTopic() {
  return (
    <>
      <Show when={state()} fallback={<Header />}>
        <Header icon={state().topic?.type}>
          <HeaderItem>{state().server?.name}</HeaderItem>
          <HeaderItem>{state().topic?.name}</HeaderItem>
        </Header>
      </Show>
      <Show when={state().topic?.type === TopicTypeText}>
        <Chat chatID={parseInt(params.topicID, 10)} type={ServersType} />
      </Show>
    </>
  );
}
