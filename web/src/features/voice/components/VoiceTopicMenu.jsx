import { Show } from "solid-js";
import {
  Camera,
  CameraOff,
  Mic,
  MicOff,
  ScreenShare,
  ScreenShareOff,
  Settings,
} from "lucide-solid";

const ICON_CLASS = "size-8 my-2 mx-4";
const ICON_BUTTON_ON_CLASS = "hover:bg-emerald-300 bg-emerald-100";
const ICON_BUTTON_OFF_CLASS = "hover:bg-rose-300 bg-rose-100";

export default function VoiceTopicMenu(props) {
  return (
    <div class="py-4 w-full flex justify-center">
      <div class="flex border-2 neo-shadow rounded-lg w-fit">
        <button
          type="button"
          aria-label="Toggle microphone"
          class={`border-r-2 rounded-l-lg ${
            props.media.micActive()
              ? ICON_BUTTON_ON_CLASS
              : ICON_BUTTON_OFF_CLASS
          }`}
          onClick={() => props.media.toggleMic()}
        >
          <Show
            when={props.media.micActive()}
            fallback={<MicOff class={ICON_CLASS} />}
          >
            <Mic class={ICON_CLASS} />
          </Show>
        </button>
        <button
          type="button"
          aria-label="Toggle screen share"
          class={`border-r-2 ${
            props.media.screenActive()
              ? ICON_BUTTON_ON_CLASS
              : ICON_BUTTON_OFF_CLASS
          }`}
          onClick={() => props.media.toggleScreenShare()}
        >
          <Show
            when={props.media.screenActive()}
            fallback={<ScreenShareOff class={ICON_CLASS} />}
          >
            <ScreenShare class={ICON_CLASS} />
          </Show>
        </button>
        <button
          type="button"
          aria-label="Toggle camera"
          class={`border-r-2 ${
            props.media.cameraActive()
              ? ICON_BUTTON_ON_CLASS
              : ICON_BUTTON_OFF_CLASS
          }`}
          onClick={() => props.media.toggleCamera()}
        >
          <Show
            when={props.media.cameraActive()}
            fallback={<CameraOff class={ICON_CLASS} />}
          >
            <Camera class={ICON_CLASS} />
          </Show>
        </button>
        <button
          type="button"
          aria-label="Open call settings"
          class="rounded-r-lg bg-white hover:bg-sky-200"
          popoverTarget={props.settingsModalID}
        >
          <Settings class={ICON_CLASS} />
        </button>
      </div>
    </div>
  );
}
