import { Maximize, Minimize, Volume2 } from "lucide-solid";
import { createSignal, onCleanup, onMount, Show } from "solid-js";

export default function VoiceTopicStreamPreview(props) {
  const [isFullscreen, setIsFullscreen] = createSignal(false);

  let cardRef;

  const toggleFullscreen = async () => {
    if (!cardRef) return;

    if (document.fullscreenElement) {
      await document.exitFullscreen();
    } else {
      await cardRef.requestFullscreen();
    }
  };

  const handleFullscreenChange = () => {
    setIsFullscreen(document.fullscreenElement === cardRef);
  };

  onMount(() => {
    document.addEventListener("fullscreenchange", handleFullscreenChange);
  });

  onCleanup(() => {
    document.removeEventListener("fullscreenchange", handleFullscreenChange);
  });

  return (
    <div
      ref={cardRef}
      class="relative group border-2 neo-shadow hover:neo-shadow-lg rounded-lg h-48 aspect-video overflow-hidden transition-all duration-300 ease-in-out"
    >
      <video
        class="w-full h-full object-cover"
        id={props.id}
        ref={props.ref}
        autoplay
        muted
      ></video>
      {/* <div class="absolute inset-x-0 bottom-0 px-2 flex items-end justify-between gap-3">
        <p class="text-white font-semibold tracking-wider text-neo-shadow">
          {props.participant?.username}
        </p>
      </div>*/}
      <div
        class={`absolute inset-x-0 bottom-0 px-2 pb-2 mx-auto flex gap-3 bg-linear-to-t from-black to-transparent ${isFullscreen() ? "justify-center gap-16" : "justify-between"}`}
      >
        <p
          class="text-white font-semibold tracking-wider text-neo-shadow"
          classList={{
            "text-3xl": isFullscreen(),
          }}
        >
          {props.participant?.username}
        </p>

        <div class="flex items-center gap-1">
          <Volume2
            class={`text-white ${isFullscreen() ? "size-10" : "size-5"}`}
          />
          <input
            type="range"
            min="0"
            max="1"
            step="0.01"
            // value={volume()}
            // onInput={handleVolume}
            class={`accent-sky-500 ${isFullscreen() ? "w-32" : "w-20"}`}
          />
        </div>

        <button
          type="button"
          onClick={toggleFullscreen}
          class="text-white hover:opacity-80 transition"
        >
          <Show
            when={isFullscreen()}
            fallback={<Maximize class="size-5 text-white" />}
          >
            <Minimize class="size-10 text-white" />
          </Show>
        </button>
      </div>
    </div>
  );
}
