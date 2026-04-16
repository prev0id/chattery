import {
  Camera,
  CameraOff,
  Mic,
  MicOff,
  ScreenShare,
  ScreenShareOff,
  Settings,
} from "lucide-solid";
import { createEffect, createSignal, onMount, Show } from "solid-js";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { UseServerContext } from "~/stores/server";

const IconClass = "size-8 my-2 mx-4";
const IconButtonOnClass = "hover:bg-emerald-300 bg-emerald-100";
const IconButtonOffClass = "hover:bg-rose-300 bg-rose-100";

export default function VoiceTopic() {
  const { currentServer, currentTopic } = UseServerContext();

  const [micOn, setMicOn] = createSignal(false);
  const [cameraOn, setCameraOn] = createSignal(false);
  const [screenShareOn, setScreenShareOn] = createSignal(false);

  var localStream = null;

  createEffect(() => {
    if (cameraOn()) {
      navigator.mediaDevices.getUserMedia({ video: true }).then((stream) => {
        localVideo.srcObject = stream;
        localStream = stream;
      });
    } else {
      localStream?.getTracks().forEach((track) => track.stop());
      localVideo.srcObject = null;
    }
  });

  return (
    <>
      <Header icon={currentTopic()?.type}>
        <HeaderItem>{currentServer()?.name}</HeaderItem>
        <HeaderItem>{currentTopic()?.name}</HeaderItem>
      </Header>
      <div class="flex flex-col h-full overflow-hidden">
        <div class="flex-1 flex gap-8 p-4 flex-wrap overflow-auto justify-center">
          <Vid id="localVideo" />
          <Vid id="remoteVideo" />
        </div>

        <div class="py-4 w-full flex justify-center">
          <div class="flex border-2 neo-shadow rounded-lg">
            <button
              class={`border-r-2 rounded-l-lg ${
                micOn() ? IconButtonOnClass : IconButtonOffClass
              }`}
              onClick={() => setMicOn(!micOn())}
            >
              <Show when={micOn()} fallback={<MicOff class={IconClass} />}>
                <Mic class={IconClass} />
              </Show>
            </button>
            <button
              class={`border-r-2 ${
                screenShareOn() ? IconButtonOnClass : IconButtonOffClass
              }`}
              onClick={() => setScreenShareOn(!screenShareOn())}
            >
              <Show
                when={screenShareOn()}
                fallback={<ScreenShareOff class={IconClass} />}
              >
                <ScreenShare class={IconClass} />
              </Show>
            </button>
            <button
              class={`border-r-2 ${cameraOn() ? IconButtonOnClass : IconButtonOffClass}`}
              onClick={() => setCameraOn(!cameraOn())}
            >
              <Show
                when={cameraOn()}
                fallback={<CameraOff class={IconClass} />}
              >
                <Camera class={IconClass} />
              </Show>
            </button>
            <button class="rounded-r-lg bg-white hover:bg-sky-200">
              <Settings class={IconClass} />
            </button>
          </div>
        </div>
      </div>
      {/*<div class="sticky bottom-0 py-4 w-full flex justify-center bg-linear-to-t from-white/10 to-transparent">*/}
    </>
  );
}

function Vid(props) {
  return (
    <div class="border-2 neo-shadow rounded-lg h-48 aspect-video overflow-hidden">
      <video
        class="w-full h-full object-cover"
        id={props.id}
        autoplay
        muted
      ></video>
    </div>
  );
}
