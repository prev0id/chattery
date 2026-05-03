import { createEffect } from "solid-js";
import VoiceTopicStreamPreview from "~/features/voice/components/VoiceTopicStreamPreview";

export default function VoiceTopicRemoteTile(props) {
  let videoEl;

  createEffect(() => {
    props.media.selectedSpeakerId();
    if (videoEl) {
      props.media.applySpeaker(videoEl);
    }
  });

  return (
    <VoiceTopicStreamPreview
      stream={props.participant.stream}
      hasVideo={props.participant.hasVideo}
      participant={props.participant.user}
      label={props.participant.label}
      volume={props.volume}
      onVolume={props.onVolume}
      onVideoReady={(el) => {
        videoEl = el;
        props.media.applySpeaker(el);
      }}
      emptyText={props.participant.label || "Remote participant"}
    />
  );
}
