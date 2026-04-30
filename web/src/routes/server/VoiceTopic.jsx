import { createEffect, createSignal, For, Show } from "solid-js";
import { createStore } from "solid-js/store";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import VoiceTopicMenu from "~/components/VoiceTopicMenu";
import VoiceTopicSettingsModal from "~/components/VoiceTopicSettingsModal";
import VoiceTopicStreamPreview from "~/components/VoiceTopicStreamPreview";
import { userData } from "~/stores/auth";
import { UseServerContext } from "~/stores/server";
import { createCallMedia } from "~/stores/voice_topic";

const settingsModalID = "voice_topic_call_settings";

export default function VoiceTopic() {
  const { currentServer, currentTopic } = UseServerContext();
  const media = createCallMedia();
  const [participants, setParticipants] = createStore([]);
  const [volumes, setVolumes] = createSignal({});

  let localCameraVideo;
  let localScreenVideo;

  function bindVideo(videoEl, stream) {
    if (!videoEl) return;
    if (videoEl.srcObject !== stream) {
      videoEl.srcObject = stream || null;
    }
  }

  createEffect(() => {
    bindVideo(localCameraVideo, media.cameraStream());
  });

  createEffect(() => {
    bindVideo(localScreenVideo, media.screenStream());
  });

  function addParticipant(participant) {
    setParticipants((current) => {
      const exists = current.some((item) => item.id === participant.id);
      if (exists) {
        return current.map((item) =>
          item.id === participant.id ? { ...item, ...participant } : item,
        );
      }
      return [...current, participant];
    });

    setVolumes((current) => ({
      ...current,
      [participant.id]: current[participant.id] ?? 0.8,
    }));
  }

  function removeParticipant(id) {
    setParticipants((current) => current.filter((item) => item.id !== id));
    setVolumes((current) => {
      const next = { ...current };
      delete next[id];
      return next;
    });
  }

  function setParticipantVolume(id, value) {
    setVolumes((current) => ({ ...current, [id]: value }));
  }

  function attachRemoteVideo(videoEl, stream, id) {
    if (!videoEl) return;

    if (videoEl.srcObject !== stream) {
      videoEl.srcObject = stream || null;
    }

    videoEl.volume = volumes()[id] ?? 0.8;
  }

  return (
    <>
      <Header icon={currentTopic()?.type}>
        <HeaderItem>{currentServer()?.name}</HeaderItem>
        <HeaderItem>{currentTopic()?.name}</HeaderItem>
      </Header>
      <div class="flex flex-col h-full overflow-hidden">
        <div class="flex-1 flex gap-8 p-4 flex-wrap overflow-auto justify-center">
          <Show when={media.cameraStream()}>
            <VoiceTopicStreamPreview
              ref={localCameraVideo}
              participant={userData()}
            />
          </Show>
          <Show when={media.screenStream()}>
            <VoiceTopicStreamPreview
              ref={localScreenVideo}
              participant={userData()}
            />
          </Show>
          {/* <For each={participants()}>
            {(participant) => {
              let remoteVideo;
              let cardRef;

              createEffect(() => {
                if (remoteVideo) {
                  attachRemoteVideo(
                    remoteVideo,
                    participant.stream,
                    participant.id,
                  );
                }
              });

              return (
                <article
                  ref={cardRef}
                  class="group overflow-hidden rounded-[28px] border border-white/10 bg-slate-900/70 shadow-xl shadow-black/30"
                >
                  <div class="relative aspect-video bg-black">
                    <Show
                      when={participant.stream}
                      fallback={
                        <div class="flex h-full w-full items-center justify-center text-sm text-slate-500">
                          Waiting for stream…
                        </div>
                      }
                    >
                      <video
                        ref={remoteVideo}
                        autoplay
                        playsInline
                        class="h-full w-full object-cover"
                      />
                    </Show>

                    <div class="absolute inset-x-0 bottom-0 flex items-end justify-between gap-3 bg-gradient-to-t from-black/80 to-transparent p-4">
                      <div>
                        <div class="text-sm font-medium">
                          {participant.name}
                        </div>
                        <div class="text-xs text-slate-300">
                          {participant.kind || "remote stream"}
                        </div>
                      </div>

                      <button
                        class="rounded-2xl bg-white/10 px-3 py-2 text-xs font-medium text-white transition hover:bg-white/15"
                        onClick={() => fullscreenNode(cardRef)}
                      >
                        Fullscreen
                      </button>
                    </div>
                  </div>

                  <div class="space-y-3 p-4">
                    <div class="flex items-center justify-between text-xs text-slate-400">
                      <span>Volume</span>
                      <span>
                        {Math.round((volumes()[participant.id] ?? 0.8) * 100)}%
                      </span>
                    </div>

                    <input
                      class="w-full accent-white"
                      type="range"
                      min="0"
                      max="1"
                      step="0.01"
                      value={volumes()[participant.id] ?? 0.8}
                      onInput={(e) =>
                        setParticipantVolume(
                          participant.id,
                          Number(e.currentTarget.value),
                        )
                      }
                    />
                  </div>
                </article>
              );
            }}
          </For>*/}
        </div>
      </div>
      <VoiceTopicMenu media={media} settingsModalID={settingsModalID} />
      <VoiceTopicSettingsModal media={media} id={settingsModalID} />
    </>
  );
}
