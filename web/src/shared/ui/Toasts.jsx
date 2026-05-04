import { X } from "lucide-solid";
import { For, Match, Switch } from "solid-js";
import ProfilePicture from "~/shared/ui/ProfilePicture";
import { removeToast, toasts } from "~/shared/stores/toast";

export default function Toasts() {
  return (
    <div class="fixed top-4 right-4 z-[60] flex flex-col gap-3 items-end pointer-events-none isolate">
      <For each={toasts}>
        {(toast) => (
          <Toast
            id={toast.id}
            variant={toast.variant}
            message={toast.message}
            data={toast.data}
          />
        )}
      </For>
    </div>
  );
}

function Toast(props) {
  const { id, variant = "error", message, data } = props;

  return (
    <div class="pointer-events-auto border-2 neo-shadow rounded-xl px-4 py-2 bg-white max-w-xs">
      <div class="flex items-start gap-3">
        <Switch
          fallback={
            <div class="flex-1">
              <Switch>
                <Match when={variant === "error"}>
                  <div class="font-bold text-red-600 mb-0.5 tracking-widest">
                    ERROR
                  </div>
                </Match>
                <Match when={variant === "warning"}>
                  <div class="font-bold text-amber-600 mb-0.5 tracking-widest">
                    WARN
                  </div>
                </Match>
                <Match when={variant === "info"}>
                  <div class="font-bold text-sky-600 mb-0.5 tracking-widest">
                    INFO
                  </div>
                </Match>
                <Match when={variant === "success"}>
                  <div class="font-bold text-green-600 mb-0.5 tracking-widest">
                    SUCCESS
                  </div>
                </Match>
              </Switch>
              <p>{message || "An error occurred"}</p>
            </div>
          }
        >
          <Match when={variant === "dm-message"}>
            <ProfilePicture
              src={data?.sender?.avatar}
              class="size-10 mt-0.5 shrink-0"
            />
            <div class="min-w-0 flex-1">
              <div class="font-semibold leading-tight truncate">
                {data?.sender?.username ?? "Direct message"}
              </div>
              <p class="text-sm leading-snug break-words">
                {message || "New message"}
              </p>
            </div>
          </Match>
        </Switch>

        <button
          type="button"
          onClick={() => removeToast(id)}
          aria-label="Close toast"
          class="flex items-center justify-center rounded-full p-1 hover:bg-red-600 hover:text-white transition-all duration-300 ease-in-out"
        >
          <X class="size-4"></X>
        </button>
      </div>
    </div>
  );
}
