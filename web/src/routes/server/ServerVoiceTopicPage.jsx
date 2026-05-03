import { createSignal } from "solid-js";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import VoiceTopicGrid from "~/components/VoiceTopicGrid";
import VoiceTopicMenu from "~/components/VoiceTopicMenu";
import VoiceTopicSettingsModal from "~/components/VoiceTopicSettingsModal";
import VoiceTopicStatus from "~/components/VoiceTopicStatus";
import { useServerContext } from "~/features/server/context";
import { userData } from "~/stores/auth";
import { createCallMedia, createVoiceCall } from "~/stores/voice_topic";

const SETTINGS_MODAL_ID = "voice_topic_call_settings";

export default function ServerVoiceTopicPage() {
  const { currentServer, currentTopic } = useServerContext();
  const media = createCallMedia();
  const call = createVoiceCall({
    topicId: () => currentTopic()?.id,
    media,
    currentUser: userData,
  });
  const [volumes, setVolumes] = createSignal({});

  function setVolume(id, value) {
    setVolumes((current) => ({ ...current, [id]: value }));
  }

  return (
    <>
      <Header icon={currentTopic()?.type}>
        <HeaderItem>{currentServer()?.name}</HeaderItem>
        <HeaderItem>{currentTopic()?.name}</HeaderItem>
      </Header>
      <div class="flex flex-col h-full overflow-hidden">
        <VoiceTopicStatus status={call.status} error={call.error} />
        <VoiceTopicGrid
          media={media}
          call={call}
          user={userData()}
          volumes={volumes}
          setVolume={setVolume}
        />
      </div>
      <VoiceTopicMenu media={media} settingsModalID={SETTINGS_MODAL_ID} />
      <VoiceTopicSettingsModal media={media} id={SETTINGS_MODAL_ID} />
    </>
  );
}
