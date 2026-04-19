import { For } from "solid-js";
import Modal from "./Modal";

export default function VoiceTopicSettingsModal(props) {
  return (
    <Modal id={props.id} name="Settings">
      <div class="space-y-4 mt-2">
        <div class="space-y-2">
          <label class="block font-semibold tracking-wider" for="select_camera">
            Camera device
          </label>
          <select
            id="select_camera"
            class="bg-amber-200 px-2 py-0.5 border-2 neo-shadow rounded-lg focus:outline-none focus:border-amber-500 w-full"
            value={props.media.selectedCameraId()}
            onChange={(e) => props.media.changeCamera(e.currentTarget.value)}
          >
            <For each={props.media.devices().videoInputs}>
              {(device, index) => (
                <option value={device.deviceId}>
                  {device.label || `Camera ${index() + 1}`}
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
            onChange={(e) => props.media.changeMic(e.currentTarget.value)}
          >
            <For each={props.media.devices().audioInputs}>
              {(device, index) => (
                <option value={device.deviceId}>
                  {device.label || `Microphone ${index() + 1}`}
                </option>
              )}
            </For>
          </select>
        </div>

        <div class="mt-2">
          <div class="font-semibold tracking-wider mb-2">Media errors</div>
          <div>Camera: {props.media.errors().camera || "none"}</div>
          <div>Mic: {props.media.errors().mic || "none"}</div>
          <div>Screen: {props.media.errors().screen || "none"}</div>
        </div>
      </div>
    </Modal>
  );
}
