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

const IconClass = "size-8 my-2 mx-4";
const IconButtonOnClass = "hover:bg-emerald-300 bg-emerald-100";
const IconButtonOffClass = "hover:bg-rose-300 bg-rose-100";

export default function VoiceTopicMenu(props) {
  return (
    <div class="py-4 w-full flex justify-center">
      <div class="flex border-2 neo-shadow rounded-lg w-fit">
        <button
          class={`border-r-2 rounded-l-lg ${
            props.media.micActive() ? IconButtonOnClass : IconButtonOffClass
          }`}
          onClick={() => props.media.toggleMic()}
        >
          <Show
            when={props.media.micActive()}
            fallback={<MicOff class={IconClass} />}
          >
            <Mic class={IconClass} />
          </Show>
        </button>
        <button
          class={`border-r-2 ${
            props.media.screenActive() ? IconButtonOnClass : IconButtonOffClass
          }`}
          onClick={() => props.media.toggleScreenShare()}
        >
          <Show
            when={props.media.screenActive()}
            fallback={<ScreenShareOff class={IconClass} />}
          >
            <ScreenShare class={IconClass} />
          </Show>
        </button>
        <button
          class={`border-r-2 ${props.media.cameraActive() ? IconButtonOnClass : IconButtonOffClass}`}
          onClick={() => props.media.toggleCamera()}
        >
          <Show
            when={props.media.cameraActive()}
            fallback={<CameraOff class={IconClass} />}
          >
            <Camera class={IconClass} />
          </Show>
        </button>
        <button
          class="rounded-r-lg bg-white hover:bg-sky-200"
          popoverTarget={props.settingsModalID}
        >
          <Settings class={IconClass} />
        </button>
      </div>
    </div>
  );
}
