import { createSignal } from "solid-js";
import Header from "~/shared/ui/Header";
import HeaderItem from "~/shared/ui/HeaderItem";
import VoiceTopicGrid from "~/features/voice/components/VoiceTopicGrid";
import VoiceTopicMenu from "~/features/voice/components/VoiceTopicMenu";
import VoiceTopicSettingsModal from "~/features/voice/components/VoiceTopicSettingsModal";
import VoiceTopicStatus from "~/features/voice/components/VoiceTopicStatus";
import { useServerContext } from "~/features/server/context";
import { userData } from "~/shared/stores/auth";
import { createVoiceCall } from "~/features/voice/call";
import { createCallMedia } from "~/features/voice/media";

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
  const [isSettingsOpen, setIsSettingsOpen] = createSignal(false);

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
      <VoiceTopicMenu
        media={media}
        onOpenSettings={() => setIsSettingsOpen(true)}
      />
      <VoiceTopicSettingsModal
        media={media}
        id={SETTINGS_MODAL_ID}
        open={isSettingsOpen()}
        onClose={() => setIsSettingsOpen(false)}
      />
    </>
  );
}
