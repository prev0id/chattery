import { User } from "lucide-solid";
import { createEffect, createSignal, Show } from "solid-js";

export default function ProfilePicture(props) {
  const [hasError, setHasError] = createSignal(false);

  createEffect(() => {
    props.src;
    setHasError(false);
  });

  return (
    <div class={`relative aspect-square ${props.class ?? ""}`}>
      <Show
        when={props.src && !hasError()}
        fallback={
          <div class="flex h-full w-full items-center justify-center border-2 neo-shadow rounded-lg bg-white">
            <User class="size-1/2 text-gray-500" />
          </div>
        }
      >
        <img
          {...props}
          class="border-2 neo-shadow rounded-lg w-full h-full"
          onError={() => setHasError(true)}
        />
      </Show>
      <Show when={props.unread}>
        <div class="absolute bottom-0 right-0 w-3 h-3 bg-rose-600 border-2 neo-shadow-sm rounded-full transform translate-x-0.5 translate-y-0.5"></div>
      </Show>
    </div>
  );
}
