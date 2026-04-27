import { Chat } from "~/components/Chat";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { WSChannelType } from "~/lib/ws";
import { UseDMContext } from "~/stores/dm";

export default function DM() {
  const { currentDM } = UseDMContext();

  return (
    <>
      <Header icon="text">
        <HeaderItem>{currentDM()?.user?.username}</HeaderItem>
      </Header>
      <Chat chatID={currentDM()?.id} type="dms" channelType={WSChannelType.DM} />
    </>
  );
}
