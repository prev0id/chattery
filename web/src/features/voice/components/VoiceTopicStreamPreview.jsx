import { Maximize, Minimize, Volume2 } from "lucide-solid";
import { createEffect, createSignal, onCleanup, onMount, Show } from "solid-js";
import ProfilePicture from "~/shared/ui/ProfilePicture";

export default function VoiceTopicStreamPreview(props) {
  const [isFullscreen, setIsFullscreen] = createSignal(false);
  const [videoReady, setVideoReady] = createSignal(false);

  let cardRef;
  let videoRef;

  const hasVideo = () =>
    props.hasVideo ?? Boolean(props.stream?.getVideoTracks?.().length);
  const showVideo = () => hasVideo() && videoReady();

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

  createEffect(() => {
    const stream = props.stream;
    if (videoRef && videoRef.srcObject !== stream) {
      setVideoReady(false);
      videoRef.srcObject = stream || null;
      videoRef.play?.().catch(() => {
        // Browsers can block autoplay with audio; the stream remains attached for the next user gesture.
      });
    }
  });

  createEffect(() => {
    if (!hasVideo()) {
      setVideoReady(false);
    }
  });

  createEffect(() => {
    if (!videoRef) return;
    props.onVideoReady?.(videoRef);
  });

  createEffect(() => {
    if (!videoRef) return;
    videoRef.volume = props.volume ?? 1;
    videoRef.muted = Boolean(props.muted);
  });

  return (
    <div
      ref={cardRef}
      class="relative group border-2 neo-shadow hover:neo-shadow-lg rounded-lg h-48 aspect-video overflow-hidden transition-all duration-300 ease-in-out"
    >
      <div class="relative h-full w-full bg-slate-900">
        <video
          class={`w-full h-full object-cover bg-black transition-opacity ${
            showVideo() ? "opacity-100" : "opacity-0"
          }`}
          id={props.id}
          ref={(el) => {
            videoRef = el;
          }}
          onLoadedData={() => setVideoReady(true)}
          onPlaying={() => setVideoReady(true)}
          onWaiting={() => setVideoReady(false)}
          onEmptied={() => setVideoReady(false)}
          autoplay
          playsInline
        ></video>
        <Show when={!showVideo()}>
          <div class="absolute inset-0 flex items-center justify-center">
            <Show
              when={props.participant?.avatar}
              fallback={
                <div class="text-center text-white font-semibold tracking-wider">
                  {props.emptyText || "No video"}
                </div>
              }
            >
              <ProfilePicture
                src={props.participant.avatar}
                class="size-24"
                alt={props.label || props.participant?.username || "Participant"}
              />
            </Show>
          </div>
        </Show>
      </div>
      <div
        class={`absolute inset-x-0 bottom-0 px-2 pb-2 mx-auto flex gap-3 bg-linear-to-t from-black to-transparent ${
          isFullscreen() ? "justify-center gap-16" : "justify-between"
        }`}
      >
        <p
          class="text-white font-semibold tracking-wider text-neo-shadow"
          classList={{
            "text-3xl": isFullscreen(),
          }}
        >
          {props.label || props.participant?.username || "Participant"}
        </p>

        <Show when={!props.hideVolume}>
          <div class="flex items-center gap-1">
            <Volume2
              class={`text-white ${isFullscreen() ? "size-10" : "size-5"}`}
            />
            <input
              type="range"
              min="0"
              max="1"
              step="0.01"
              value={props.volume ?? 1}
              onInput={(event) =>
                props.onVolume?.(Number(event.currentTarget.value))
              }
              class={`accent-sky-500 ${isFullscreen() ? "w-32" : "w-20"}`}
            />
          </div>
        </Show>

        <button
          type="button"
          onClick={toggleFullscreen}
          aria-label="Toggle fullscreen"
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
