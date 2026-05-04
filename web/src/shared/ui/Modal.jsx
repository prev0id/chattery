import { X } from "lucide-solid";
import { createEffect, onCleanup, Show } from "solid-js";

export default function Modal(props) {
  createEffect(() => {
    if (!props.open) return;

    const handleKeyDown = (event) => {
      if (event.key === "Escape") {
        props.onClose?.();
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    onCleanup(() => document.removeEventListener("keydown", handleKeyDown));
  });

  return (
    <Show when={props.open}>
      <div class="fixed inset-0 z-40 bg-rose-50/20 backdrop-blur-sm" />
      <div
        id={props.id}
        role="dialog"
        aria-modal="true"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        onMouseDown={(event) => {
          if (event.target === event.currentTarget) {
            props.onClose?.();
          }
        }}
      >
        <div class="border-3 w-md neo-shadow-lg rounded-xl bg-white p-4">
          <div class="flex justify-between items-center">
            <h1 class="text-2xl font-semibold tracking-wider">{props.name}</h1>
            <button
              type="button"
              aria-label="Close modal"
              class="flex items-center justify-center rounded-full p-1 hover:bg-red-600 hover:text-white"
              onClick={() => props.onClose?.()}
            >
              <X class="size-5"></X>
            </button>
          </div>
          {props.children}
        </div>
      </div>
    </Show>
  );
}
