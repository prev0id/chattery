import { ChevronRight, MessagesSquare, Mic, Sparkles } from "lucide-solid";
import { Match, Switch } from "solid-js";
import { TopicTypeText, TopicTypeVoice } from "./Chat";

const iconClasses = "size-8 mr-4";
const textClasses = "text-2xl font-semibold tracking-wider";

export default function AppHeader(props) {
  const hasServer = () => !!props.serverName;
  const hasDM = () => {
    const username =
      typeof props.dmUsername === "function"
        ? props.dmUsername()
        : props.dmUsername;
    return !!username;
  };

  const dmUsernameValue = () => {
    return typeof props.dmUsername === "function"
      ? props.dmUsername()
      : props.dmUsername;
  };

  return (
    <div class="border-b-3 px-4 py-2 flex items-center  bg-emerald-50">
      <Switch fallback={DefaultHeader()}>
        <Match when={hasServer()}>
          <Switch>
            <Match when={props.topicType === TopicTypeText}>
              <MessagesSquare class={iconClasses} />
            </Match>
            <Match when={props.topicType === TopicTypeVoice}>
              <Mic class={iconClasses} />
            </Match>
          </Switch>

          <h1 class={textClasses}>{props.serverName}</h1>

          <ChevronRight class="size-6" />

          <h1 class={textClasses}>{props.topicName}</h1>
        </Match>
        <Match when={hasDM()}>
          <MessagesSquare class={iconClasses} />
          <h1 class={textClasses}>{dmUsernameValue()}</h1>
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
