import VoiceTopicStreamPreview from "~/features/voice/components/VoiceTopicStreamPreview";

export default function VoiceTopicLocalTile(props) {
  const cameraStream = () => props.media.cameraStream();
  const screenStream = () => props.media.screenStream();
  const fallbackStream = () => cameraStream() || screenStream();

  return (
    <VoiceTopicStreamPreview
      stream={fallbackStream()}
      hasVideo={Boolean(fallbackStream()?.getVideoTracks().length)}
      muted
      hideVolume
      participant={props.user}
      label={props.user?.username || "You"}
      emptyText="You are in call"
    />
  );
}
