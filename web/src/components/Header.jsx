import { children, Match, Show, Switch } from "solid-js";
import { ChevronRight, MessagesSquare, Mic, Sparkles } from "lucide-solid";

const iconClasses = "size-8 mr-4";

export default function Header(props) {
  const items = children(() => props.children);

  console.log(items.toArray().length);

  return (
    <header class="border-b-3 px-4 py-2 flex items-center bg-emerald-50">
      <Switch fallback={<Sparkles class={iconClasses} />}>
        <Match when={props.icon === "text"}>
          <MessagesSquare class={iconClasses} />
        </Match>
        <Match when={props.icon === "voice"}>
          <Mic class={iconClasses} />
        </Match>
      </Switch>
      <For each={items.toArray()}>
        {(item, index) => (
          <>
            {item}
            <Show when={index() + 1 !== items.toArray().length}>
              <ChevronRight class="size-6" />
            </Show>
          </>
        )}
      </For>
    </header>
  );
}
