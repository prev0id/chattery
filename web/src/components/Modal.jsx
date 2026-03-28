import { X } from "lucide-solid";

export default function Modal(props) {
  return (
    <div
      popover
      id={props.id}
      class="border-3 w-md neo-shadow-lg rounded-xl m-auto p-4 backdrop:backdrop-blur-sm backdrop:bg-rose-50/20"
    >
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold tracking-wider">{props.name}</h1>
        <button
          onclick={() => document.getElementById(props.id).hidePopover()}
          class="flex items-center justify-center rounded-full p-1 hover:bg-red-600 hover:text-white"
        >
          <X class="size-5"></X>
        </button>
      </div>
      {props.children}
    </div>
  );
}
