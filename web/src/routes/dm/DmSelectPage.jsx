import { ArrowLeftToLine } from "lucide-solid";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";

export default function DmSelectPage() {
  return (
    <>
      <Header>
        <HeaderItem>Choose a chat to get started</HeaderItem>
      </Header>
      <div class="m-auto flex items-center justify-center gap-4">
        <ArrowLeftToLine class="size-10" />
        <span class="text-2xl tracking-wider font-semibold">
          Join or create a chat to start chatting
        </span>
      </div>
    </>
  );
}
