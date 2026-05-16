import { createEffect, For, Show } from "solid-js";
import Modal from "~/shared/ui/Modal";

export default function VoiceTopicSettingsModal(props) {
  createEffect(() => {
    if (!props.open) return;
    void props.media.ensureDeviceLabels();
  });

  return (
    <Modal
      id={props.id}
      name="Settings"
      open={props.open}
      onClose={props.onClose}
    >
      <div class="space-y-4 mt-2">
        <div class="space-y-2">
          <label class="block font-semibold tracking-wider" for="select_camera">
            Camera device
          </label>
          <select
            id="select_camera"
            class="bg-amber-200 px-2 py-0.5 border-2 neo-shadow rounded-lg focus:outline-none focus:border-amber-500 w-full"
            value={props.media.selectedCameraId()}
            onChange={(event) =>
              props.media.changeCamera(event.currentTarget.value)
            }
          >
            <For each={props.media.devices().videoInputs}>
              {(device, index) => (
                <option value={device.deviceId}>
                  {device.label ||
                    props.media.getHiddenDeviceLabel("videoinput", index())}
                </option>
              )}
            </For>
          </select>
        </div>

        <div class="space-y-2">
          <label
            class="block font-semibold tracking-wider"
            for="select_microphone"
          >
            Microphone device
          </label>
          <select
            id="select_microphone"
            class="bg-sky-200 px-2 py-0.5 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
            value={props.media.selectedMicId()}
            onChange={(event) => props.media.changeMic(event.currentTarget.value)}
          >
            <For each={props.media.devices().audioInputs}>
              {(device, index) => (
                <option value={device.deviceId}>
                  {device.label ||
                    props.media.getHiddenDeviceLabel("audioinput", index())}
                </option>
              )}
            </For>
          </select>
        </div>

        <div class="space-y-2">
          <label class="block font-semibold tracking-wider" for="select_speaker">
            Speaker device
          </label>
          <select
            id="select_speaker"
            class="bg-emerald-200 px-2 py-0.5 border-2 neo-shadow rounded-lg focus:outline-none focus:border-emerald-500 w-full"
            value={props.media.selectedSpeakerId()}
            disabled={!props.media.supportsSinkID}
            onChange={(event) =>
              props.media.changeSpeaker(event.currentTarget.value)
            }
          >
            <For each={props.media.devices().audioOutputs}>
              {(device, index) => (
                <option value={device.deviceId}>
                  {device.label ||
                    props.media.getHiddenDeviceLabel("audiooutput", index())}
                </option>
              )}
            </For>
          </select>
          {!props.media.supportsSinkID && (
            <p class="text-sm text-rose-700">
              Output device selection is not supported by this browser.
            </p>
          )}
        </div>

        <Show when={props.media.deviceLabelsPending()}>
          <p class="text-sm text-slate-600">
            Loading real device names from browser permissions...
          </p>
        </Show>
        <Show when={props.media.deviceLabelsError()}>
          <p class="text-sm text-rose-700">{props.media.deviceLabelsError()}</p>
        </Show>

        <div class="mt-2">
          <div class="font-semibold tracking-wider mb-2">Media errors</div>
          <div>Camera: {props.media.errors().camera || "none"}</div>
          <div>Mic: {props.media.errors().mic || "none"}</div>
          <div>Speaker: {props.media.errors().speaker || "none"}</div>
          <div>Screen: {props.media.errors().screen || "none"}</div>
        </div>
      </div>
    </Modal>
  );
}
