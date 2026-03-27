import { onCleanup } from "solid-js";

export default function Toast(props) {
  const { status, message, onClose } = props;

  const timer = setTimeout(() => {
    onClose?.();
  }, 5000);

  onCleanup(() => clearTimeout(timer));

  return (
    <div class="fixed top-4 right-4 z-50 border-3 neo-shadow-lg rounded-xl p-4 bg-white max-w-xs">
      <div class="flex items-start gap-3">
        <div class="flex-1">
          {status && (
            <div class="font-bold text-red-500 mb-0.5 tracking-widest">
              {status}
            </div>
          )}
          <p>{message || "An error occurred"}</p>
        </div>

        <button
          type="button"
          onClick={onClose}
          class="rounded-full px-1 py-0.5 leading-none hover:bg-red-600 transition-all duration-300 ease-in-out"
        >
          ✕
        </button>
      </div>
    </div>
  );
}
