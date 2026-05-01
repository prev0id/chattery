import { For, Show } from "solid-js";
import VoiceTopicLocalTile from "~/components/VoiceTopicLocalTile";
import VoiceTopicRemoteTile from "~/components/VoiceTopicRemoteTile";

export default function VoiceTopicGrid(props) {
  return (
    <div class="flex-1 flex gap-8 p-4 flex-wrap overflow-auto justify-center">
      <VoiceTopicLocalTile media={props.media} user={props.user} />
      <For each={props.call.participants}>
        {(participant) => (
          <VoiceTopicRemoteTile
            media={props.media}
            participant={participant}
            volume={props.volumes()[participant.id] ?? 0.8}
            onVolume={(value) => props.setVolume(participant.id, value)}
          />
        )}
      </For>
      <Show when={props.call.participants.length === 0}>
        <div class="m-auto text-lg font-semibold tracking-wider">
          Waiting for others to join
        </div>
      </Show>
    </div>
  );
}
