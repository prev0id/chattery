import HeaderItem from "~/components/HeaderItem";
import Header from "~/components/Header";
import { ArrowLeftToLine } from "lucide-solid";

export default function Select() {
  return (
    <>
      <Header>
        <HeaderItem>Choose a topic to get started</HeaderItem>
      </Header>
      <div class="m-auto flex items-center justify-center gap-4">
        <ArrowLeftToLine class="size-10" />
        <span class="text-2xl tracking-wider font-semibold">
          Join or create a server to start chatting
        </span>
      </div>
    </>
  );
}
