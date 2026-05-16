import { createEffect, createSignal, onCleanup, untrack } from "solid-js";
import { createStore } from "solid-js/store";
import { getVoiceIceServers } from "~/features/voice/api";
import { WS_CHANNEL_TYPE, WS_EVENT_TYPE } from "~/lib/ws";
import { appWebSocket } from "~/shared/stores/websocket";

const ICE_BATCH_DELAY_MS = 50;

function parsePayload(payload) {
  if (typeof payload !== "string") return payload;

  try {
    return JSON.parse(payload);
  } catch {
    return payload;
  }
}

function sameChannel(left, right) {
  return left?.type === right?.type && Number(left?.id) === Number(right?.id);
}

function mergeTracks(...trackGroups) {
  const stream = new MediaStream();
  trackGroups
    .flat()
    .filter(Boolean)
    .forEach((track) => {
      if (!stream.getTracks().some((item) => item.id === track.id)) {
        stream.addTrack(track);
      }
    });
  return stream;
}

/**
 * Creates voice-topic call signaling and peer-connection state.
 *
 * @param {{topicId: import("solid-js").Accessor<number>, media: ReturnType<import("~/features/voice/media").createCallMedia>, currentUser: Function}} props
 * @returns {Object}
 */
export function createVoiceCall(props) {
  const [status, setStatus] = createSignal("idle");
  const [error, setError] = createSignal("");
  const [participants, setParticipants] = createStore([]);

  const senders = new Map();
  const iceCandidateQueue = [];
  const pendingRemoteStreams = [];
  let pc = null;
  let currentChannel = null;
  let makingOffer = false;
  let pendingNegotiation = false;
  let pendingServerOffer = null;
  let iceFlushTimer = null;
  let negotiationQueue = Promise.resolve();
  let startGeneration = 0;

  function topicId() {
    const value = props.topicId?.();
    return value ? Number(value) : 0;
  }

  function currentUser() {
    return typeof props.currentUser === "function"
      ? props.currentUser()
      : props.currentUser;
  }

  function sendSignal(type, payload) {
    if (!currentChannel) return;
    appWebSocket.sendEvent(type, currentChannel, payload);
  }

  function setCallError(err, fallback) {
    const message = err?.userMessage || fallback;
    setError(message);
    setStatus("error");
  }

  function upsertParticipant(participant) {
    const id = Number(participant.id);
    if (!id || currentUser()?.id === id) return;

    setParticipants((current) => {
      if (current.some((item) => item.id === id)) {
        return current.map((item) =>
          item.id === id ? { ...item, ...participant, id } : item,
        );
      }
      return [...current, { ...participant, id }];
    });
  }

  function removeParticipant(id) {
    setParticipants((current) =>
      current.filter((item) => item.id !== Number(id)),
    );
  }

  function takePendingRemoteStream(userId) {
    const index = pendingRemoteStreams.findIndex(
      (item) => Number(item.userId) === Number(userId),
    );
    if (index >= 0) {
      return pendingRemoteStreams.splice(index, 1)[0];
    }

    return pendingRemoteStreams.shift() ?? null;
  }

  function attachRemoteTrackToParticipant(userId, stream, track, hasVideo) {
    if (!userId || (!stream && !track)) return;

    setParticipants((current) => {
      const existing = current.find((item) => item.id === Number(userId));
      const nextStream = mergeTracks(
        existing?.stream?.getTracks() ?? [],
        stream?.getTracks() ?? [],
        [track],
      );

      return current.map((item) =>
        item.id === Number(userId)
          ? {
              ...item,
              stream: nextStream,
              hasVideo:
                hasVideo ||
                item.hasVideo ||
                nextStream
                  .getVideoTracks()
                  .some((videoTrack) => videoTrack.readyState === "live"),
            }
          : item,
      );
    });
  }

  function findRemoteUserIdForStream(stream) {
    const streamUserId = Number(stream?.id);
    if (
      streamUserId &&
      streamUserId !== currentUser()?.id &&
      participants.some((participant) => participant.id === streamUserId)
    ) {
      return streamUserId;
    }

    return participants.find((participant) => !participant.stream)?.id;
  }

  function setVoiceState(payload) {
    const next = (payload?.participants || [])
      .map((participant) => ({
        id: Number(participant.user_id),
        user: participant.sender,
        label: participant.sender?.username || `User #${participant.user_id}`,
      }))
      .filter(
        (participant) => participant.id && participant.id !== currentUser()?.id,
      );

    setParticipants((current) =>
      next.map((participant) => {
        const existing = current.find((item) => item.id === participant.id);
        const pendingStream = existing?.stream
          ? null
          : takePendingRemoteStream(participant.id);
        return {
          ...participant,
          ...existing,
          stream:
            existing?.stream ??
            (pendingStream
              ? mergeTracks(pendingStream.stream?.getTracks() ?? [], [
                  pendingStream.track,
                ])
              : undefined),
          hasVideo: existing?.hasVideo ?? pendingStream?.hasVideo,
        };
      }),
    );
    setStatus("connected");
    setError("");
  }

  function flushICECandidates() {
    iceFlushTimer = null;
    if (!iceCandidateQueue.length) return;

    const candidates = iceCandidateQueue.splice(0, iceCandidateQueue.length);
    sendSignal(WS_EVENT_TYPE.VoiceICECandidates, { candidates });
  }

  function enqueueICECandidate(candidate) {
    iceCandidateQueue.push(candidate.toJSON());
    if (!iceFlushTimer) {
      iceFlushTimer = setTimeout(flushICECandidates, ICE_BATCH_DELAY_MS);
    }
  }

  function addRemoteTrack(event) {
    const [eventStream] = event.streams;
    const stream = eventStream || new MediaStream([event.track]);
    const isVideo = event.track.kind === "video";
    const hasVideo =
      isVideo ||
      stream
        .getVideoTracks()
        .some((track) => !track.muted && track.readyState === "live");
    const streamUserId = Number(stream?.id);
    const userId = findRemoteUserIdForStream(stream);

    if (userId && currentUser()?.id !== userId) {
      attachRemoteTrackToParticipant(userId, stream, event.track, hasVideo);
    } else {
      pendingRemoteStreams.push({
        userId: streamUserId,
        stream,
        track: event.track,
        hasVideo,
      });
    }

    stream.getVideoTracks().forEach((track) => {
      track.onmute = () => {
        if (userId) upsertParticipant({ id: userId, stream, hasVideo: false });
      };
      track.onunmute = () => {
        if (userId) upsertParticipant({ id: userId, stream, hasVideo: true });
      };
    });

    event.track.onended = () => {
      if (userId) upsertParticipant({ id: userId, stream, hasVideo: false });
    };
  }

  function negotiate() {
    negotiationQueue = negotiationQueue
      .then(async () => {
        if (!pc || pc.signalingState !== "stable") {
          pendingNegotiation = true;
          return;
        }

        try {
          makingOffer = true;
          const offer = await pc.createOffer();
          if (!pc || pc.signalingState !== "stable") {
            pendingNegotiation = true;
            return;
          }
          await pc.setLocalDescription(offer);
          sendSignal(WS_EVENT_TYPE.VoiceOffer, {
            type: offer.type,
            sdp: offer.sdp,
          });
        } catch (err) {
          if (pc?.signalingState === "stable") {
            setCallError(err, "Unable to negotiate voice connection");
          } else {
            pendingNegotiation = true;
          }
        } finally {
          makingOffer = false;
        }
      })
      .catch((err) => {
        setCallError(err, "Unable to negotiate voice connection");
      });

    return negotiationQueue;
  }

  async function syncSender(key, track, stream, shouldNegotiate = true) {
    if (!pc) return;

    const existing = senders.get(key);
    if (track && existing) {
      await existing.replaceTrack(track);
      if (key === "video" && typeof existing.generateKeyFrame === "function") {
        existing.generateKeyFrame().catch(() => {
          // Keyframe generation is best-effort and unsupported by some browsers.
        });
      }
      return;
    }

    if (track && !existing) {
      senders.set(key, pc.addTrack(track, stream));
      if (shouldNegotiate) {
        await negotiate();
      }
      return;
    }

    if (!track && existing) {
      pc.removeTrack(existing);
      senders.delete(key);
      if (shouldNegotiate) {
        await negotiate();
      }
    }
  }

  async function syncLocalTracks(shouldNegotiate = true) {
    if (!pc) return;

    try {
      const camera = props.media.cameraStream();
      const mic = props.media.micStream();
      const screen = props.media.screenStream();
      const video = screen || camera;

      await syncSender(
        "video",
        video?.getVideoTracks()[0],
        video,
        shouldNegotiate,
      );
      await syncSender("mic", mic?.getAudioTracks()[0], mic, shouldNegotiate);
      await syncSender(
        "screen-audio",
        screen?.getAudioTracks()[0],
        screen,
        shouldNegotiate,
      );
    } catch (err) {
      setCallError(err, "Unable to update local media tracks");
    }
  }

  function createPeerConnection(iceServers) {
    const connection = new RTCPeerConnection({ iceServers });

    connection.addTransceiver("audio", { direction: "recvonly" });
    connection.addTransceiver("video", { direction: "recvonly" });

    connection.onicecandidate = (event) => {
      if (!event.candidate) return;
      enqueueICECandidate(event.candidate);
    };

    connection.ontrack = addRemoteTrack;
    connection.onsignalingstatechange = () => {
      if (connection.signalingState === "stable" && pendingServerOffer) {
        const offer = pendingServerOffer;
        pendingServerOffer = null;
        handleVoiceOffer(offer).catch((err) => {
          setCallError(err, "Unable to handle voice offer");
        });
      }
    };

    connection.onconnectionstatechange = () => {
      const state = connection.connectionState;
      if (state === "connected") {
        setStatus("connected");
        setError("");
      } else if (state === "connecting") {
        setStatus("connecting");
      } else if (["failed", "closed", "disconnected"].includes(state)) {
        setStatus(state);
      }
    };

    return connection;
  }

  async function handleVoiceAnswer(payload) {
    if (!pc) return;
    const desc = parsePayload(payload);
    await pc.setRemoteDescription(desc);

    if (pendingNegotiation) {
      pendingNegotiation = false;
      await negotiate();
    }
  }

  async function handleVoiceOffer(payload) {
    if (!pc) return;
    const desc = parsePayload(payload);
    const readyForOffer = !makingOffer && pc.signalingState === "stable";
    if (!readyForOffer) {
      if (pc.signalingState === "have-local-offer") {
        try {
          await pc.setLocalDescription({ type: "rollback" });
          makingOffer = false;
          pendingNegotiation = false;
        } catch {
          pendingServerOffer = payload;
          return;
        }
      } else {
        pendingServerOffer = payload;
        return;
      }
    }

    await pc.setRemoteDescription(desc);
    const answer = await pc.createAnswer();
    await pc.setLocalDescription(answer);
    sendSignal(WS_EVENT_TYPE.VoiceAnswer, {
      type: answer.type,
      sdp: answer.sdp,
    });
  }

  async function handleVoiceICECandidate(payload) {
    if (!pc) return;
    const candidate = parsePayload(payload);
    if (!candidate?.candidate) return;
    await pc.addIceCandidate(candidate);
  }

  async function handleVoiceICECandidates(payload) {
    if (!pc) return;
    const parsed = parsePayload(payload);
    const candidates = parsed?.candidates || [];

    for (const candidate of candidates) {
      if (!candidate?.candidate) continue;
      await pc.addIceCandidate(candidate);
    }
  }

  async function handleVoiceEvent(event) {
    if (!currentChannel || !sameChannel(event.channel, currentChannel)) return;

    try {
      if (event.type === WS_EVENT_TYPE.VoiceAnswer) {
        await handleVoiceAnswer(event.payload);
      } else if (event.type === WS_EVENT_TYPE.VoiceOffer) {
        await handleVoiceOffer(event.payload);
      } else if (event.type === WS_EVENT_TYPE.VoiceICECandidate) {
        await handleVoiceICECandidate(event.payload);
      } else if (event.type === WS_EVENT_TYPE.VoiceICECandidates) {
        await handleVoiceICECandidates(event.payload);
      } else if (event.type === WS_EVENT_TYPE.VoiceState) {
        setVoiceState(parsePayload(event.payload));
      } else if (event.type === WS_EVENT_TYPE.VoiceJoined) {
        const payload = parsePayload(event.payload);
        upsertParticipant({
          id: payload?.user_id,
          user: payload?.sender,
          label: payload?.sender?.username || `User #${payload?.user_id}`,
        });
        setStatus("connected");
        setError("");
      } else if (event.type === WS_EVENT_TYPE.VoiceLeft) {
        const payload = parsePayload(event.payload);
        removeParticipant(payload?.user_id);
      }
    } catch (err) {
      setCallError(err, "Unable to handle voice event");
    }
  }

  async function start() {
    const id = topicId();
    if (!id) return;

    stop();
    const generation = startGeneration;
    const channel = { type: WS_CHANNEL_TYPE.VoiceTopic, id };

    setStatus("connecting");
    setError("");
    setParticipants([]);

    let iceServers;
    try {
      iceServers = await getVoiceIceServers();
    } catch (err) {
      if (generation === startGeneration) {
        setCallError(err, "Unable to load voice connection settings");
      }
      return;
    }
    if (generation !== startGeneration) {
      return;
    }

    currentChannel = channel;
    pc = createPeerConnection(iceServers);
    appWebSocket.join(currentChannel);
    await syncLocalTracks(false);
    if (
      generation !== startGeneration ||
      !sameChannel(currentChannel, channel)
    ) {
      return;
    }
    await negotiate();
  }

  function stop() {
    startGeneration += 1;
    if (currentChannel) {
      appWebSocket.leave(currentChannel);
    }
    senders.clear();
    iceCandidateQueue.splice(0, iceCandidateQueue.length);
    pendingRemoteStreams.splice(0, pendingRemoteStreams.length);
    if (iceFlushTimer) {
      clearTimeout(iceFlushTimer);
      iceFlushTimer = null;
    }
    pendingNegotiation = false;
    pendingServerOffer = null;
    makingOffer = false;
    negotiationQueue = Promise.resolve();
    currentChannel = null;
    setParticipants([]);

    if (pc) {
      pc.onicecandidate = null;
      pc.ontrack = null;
      pc.onconnectionstatechange = null;
      pc.close();
      pc = null;
    }
    setStatus("idle");
  }

  const unsubscribeEvent = appWebSocket.subscribeEvent(handleVoiceEvent);
  const unsubscribeError = appWebSocket.subscribeError((payload) => {
    setCallError(payload, payload?.message || "Voice websocket error");
  });

  createEffect(() => {
    const id = topicId();
    if (!id) return;
    untrack(() => {
      start();
    });
  });

  createEffect(() => {
    props.media.cameraStream();
    props.media.micStream();
    props.media.screenStream();
    if (untrack(() => pc)) {
      syncLocalTracks();
    }
  });

  onCleanup(() => {
    unsubscribeEvent();
    unsubscribeError();
    stop();
  });

  return {
    status,
    error,
    participants,
    currentUser,
    stop,
  };
}
