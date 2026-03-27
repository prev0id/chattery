import { ChevronRight, MessagesSquare, Mic, Sparkles } from "lucide-solid";
import { Match, Switch } from "solid-js";
import { selectedDM, selectedServer, selectedTopic } from "../stores/app";

const iconClasses = "size-8 mr-4";
const textClasses = "text-2xl font-semibold tracking-wider";

export default function AppHeader(props) {
  return (
    <div class="border-b-3 px-4 py-2 flex items-center  bg-emerald-50">
      <Switch fallback={DefaultHeader()}>
        <Match when={selectedTopic() !== null}>
          <Switch>
            <Match when={selectedTopic()?.type === "text"}>
              <MessagesSquare class={iconClasses} />
            </Match>
            <Match when={selectedTopic()?.type === "voice"}>
              <Mic class={iconClasses} />
            </Match>
          </Switch>

          <h1 class={textClasses}>{selectedServer()?.name}</h1>

          <ChevronRight class="size-6" />

          <h1 class={textClasses}>{selectedTopic()?.name}</h1>
        </Match>
        <Match when={selectedDM() !== null}>
          <MessagesSquare class={iconClasses} />
          <h1 class={textClasses}>{selectedDM()?.username}</h1>
        </Match>
      </Switch>
    </div>
  );
}

function DefaultHeader() {
  return (
    <>
      <Sparkles class={iconClasses} />
      <h1 class={textClasses}>{"Join a chat to get started!"}</h1>
    </>
  );
}
