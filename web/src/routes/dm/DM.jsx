import { Chat } from "~/components/Chat";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { UseDMContext } from "~/stores/dm";

export default function DM() {
  const { currentDM } = UseDMContext();

  return (
    <>
      <Header icon="text">
        <HeaderItem>{currentDM()?.name}</HeaderItem>
      </Header>
      <Chat chatID={currentDM()?.id} type="dms" />
    </>
  );
}
