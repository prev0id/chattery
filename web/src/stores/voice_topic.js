import { batch, createEffect, createSignal, onCleanup, onMount, untrack } from "solid-js";
import { createStore } from "solid-js/store";
import { WSChannelType, WSEventType } from "~/lib/ws";
import { appWebSocket } from "~/stores/websocket";

const storageKeys = {
  camera: "chattery.voice.camera_id",
  mic: "chattery.voice.mic_id",
  speaker: "chattery.voice.speaker_id",
};

function stopStream(stream) {
  if (!stream) return;
  stream.getTracks().forEach((track) => track.stop());
}

function savedDevice(key) {
  try {
    return localStorage.getItem(key) || "";
  } catch {
    return "";
  }
}

function saveDevice(key, value) {
  try {
    if (value) {
      localStorage.setItem(key, value);
    } else {
      localStorage.removeItem(key);
    }
  } catch {
    // localStorage can be unavailable in private contexts.
  }
}

function availableDeviceID(devices, preferred) {
  if (preferred && devices.some((device) => device.deviceId === preferred)) {
    return preferred;
  }
  return devices[0]?.deviceId || "";
}

function videoConstraint(deviceId) {
  return deviceId ? { deviceId: { exact: deviceId } } : true;
}

function audioConstraint(deviceId) {
  return deviceId ? { deviceId: { exact: deviceId } } : true;
}

function parsePayload(payload) {
  if (typeof payload !== "string") return payload;

  try {
    return JSON.parse(payload);
  } catch {
    return payload;
  }
}

function sameChannel(left, right) {
  return (
    left?.type === right?.type && Number(left?.id) === Number(right?.id)
  );
}

function supportsSinkID() {
  return typeof HTMLMediaElement !== "undefined" &&
    "setSinkId" in HTMLMediaElement.prototype;
}

export function createCallMedia() {
  const [devices, setDevices] = createSignal({
    videoInputs: [],
    audioInputs: [],
    audioOutputs: [],
  });

  const [selectedCameraId, setSelectedCameraId] = createSignal(
    savedDevice(storageKeys.camera),
  );
  const [selectedMicId, setSelectedMicId] = createSignal(
    savedDevice(storageKeys.mic),
  );
  const [selectedSpeakerId, setSelectedSpeakerId] = createSignal(
    savedDevice(storageKeys.speaker),
  );

  const [cameraStream, setCameraStream] = createSignal(null);
  const [micStream, setMicStream] = createSignal(null);
  const [screenStream, setScreenStream] = createSignal(null);

  const [cameraActive, setCameraActive] = createSignal(false);
  const [micActive, setMicActive] = createSignal(false);
  const [screenActive, setScreenActive] = createSignal(false);

  const [errors, setErrors] = createSignal({
    camera: "",
    mic: "",
    screen: "",
    speaker: "",
  });

  async function refreshDevices() {
    if (!navigator.mediaDevices?.enumerateDevices) return;

    const all = await navigator.mediaDevices.enumerateDevices();
    const videoInputs = all.filter((d) => d.kind === "videoinput");
    const audioInputs = all.filter((d) => d.kind === "audioinput");
    const audioOutputs = all.filter((d) => d.kind === "audiooutput");

    setDevices({
      videoInputs,
      audioInputs,
      audioOutputs,
    });

    setSelectedCameraId((current) => {
      const next = availableDeviceID(
        videoInputs,
        current || savedDevice(storageKeys.camera),
      );
      saveDevice(storageKeys.camera, next);
      return next;
    });
    setSelectedMicId((current) => {
      const next = availableDeviceID(
        audioInputs,
        current || savedDevice(storageKeys.mic),
      );
      saveDevice(storageKeys.mic, next);
      return next;
    });
    setSelectedSpeakerId((current) => {
      const next = availableDeviceID(
        audioOutputs,
        current || savedDevice(storageKeys.speaker),
      );
      saveDevice(storageKeys.speaker, next);
      return next;
    });
  }

  async function startCamera(deviceId = selectedCameraId()) {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: videoConstraint(deviceId),
        audio: false,
      });

      const previousCamera = cameraStream();
      const previousScreen = screenStream();
      batch(() => {
        setScreenStream(null);
        setScreenActive(false);
        setCameraStream(stream);
        setCameraActive(true);
      });
      stopStream(previousScreen);
      stopStream(previousCamera);
      setErrors((prev) => ({ ...prev, camera: "" }));
      return stream;
    } catch (err) {
      setErrors((prev) => ({
        ...prev,
        camera: err?.message || "Unable to start camera",
      }));
      throw err;
    }
  }

  function stopCamera() {
    const previousCamera = cameraStream();
    setCameraStream(null);
    setCameraActive(false);
    stopStream(previousCamera);
  }

  async function toggleCamera() {
    if (cameraActive()) {
      stopCamera();
      return;
    }
    await startCamera();
  }

  async function changeCamera(deviceId) {
    setSelectedCameraId(deviceId);
    saveDevice(storageKeys.camera, deviceId);
    if (cameraActive()) {
      await startCamera(deviceId);
    }
  }

  async function startMic(deviceId = selectedMicId()) {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: false,
        audio: audioConstraint(deviceId),
      });

      stopStream(micStream());
      setMicStream(stream);
      setMicActive(true);
      setErrors((prev) => ({ ...prev, mic: "" }));
      return stream;
    } catch (err) {
      setErrors((prev) => ({
        ...prev,
        mic: err?.message || "Unable to start microphone",
      }));
      throw err;
    }
  }

  function stopMic() {
    stopStream(micStream());
    setMicStream(null);
    setMicActive(false);
  }

  async function toggleMic() {
    if (micActive()) {
      stopMic();
      return;
    }
    await startMic();
  }

  async function changeMic(deviceId) {
    setSelectedMicId(deviceId);
    saveDevice(storageKeys.mic, deviceId);
    if (micActive()) {
      await startMic(deviceId);
    }
  }

  function changeSpeaker(deviceId) {
    setSelectedSpeakerId(deviceId);
    saveDevice(storageKeys.speaker, deviceId);
    setErrors((prev) => ({ ...prev, speaker: "" }));
  }

  async function applySpeaker(videoEl) {
    if (!videoEl || !supportsSinkID()) return;

    try {
      await videoEl.setSinkId(selectedSpeakerId());
      setErrors((prev) => ({ ...prev, speaker: "" }));
    } catch (err) {
      setErrors((prev) => ({
        ...prev,
        speaker: err?.message || "Unable to change output device",
      }));
    }
  }

  async function startScreenShare() {
    try {
      const stream = await navigator.mediaDevices
        .getDisplayMedia({
          video: true,
          audio: true,
        })
        .catch(() =>
          navigator.mediaDevices.getDisplayMedia({
            video: true,
            audio: false,
          }),
        );

      const previousCamera = cameraStream();
      const previousScreen = screenStream();
      batch(() => {
        setCameraStream(null);
        setCameraActive(false);
        setScreenStream(stream);
        setScreenActive(true);
      });
      stopStream(previousCamera);
      stopStream(previousScreen);
      setErrors((prev) => ({ ...prev, screen: "" }));

      const [videoTrack] = stream.getVideoTracks();
      if (videoTrack) {
        videoTrack.onended = () => {
          stopScreenShare();
        };
      }

      return stream;
    } catch (err) {
      setErrors((prev) => ({
        ...prev,
        screen: err?.message || "Unable to start screen share",
      }));
      throw err;
    }
  }

  function stopScreenShare() {
    const previousScreen = screenStream();
    setScreenStream(null);
    setScreenActive(false);
    stopStream(previousScreen);
  }

  async function toggleScreenShare() {
    if (screenActive()) {
      stopScreenShare();
      return;
    }
    await startScreenShare();
  }

  function stopAll() {
    stopCamera();
    stopMic();
    stopScreenShare();
  }

  function getPublishStream() {
    const stream = new MediaStream();

    const cam = cameraStream();
    const mic = micStream();
    const screen = screenStream();

    if (cam) cam.getVideoTracks().forEach((track) => stream.addTrack(track));
    if (mic) mic.getAudioTracks().forEach((track) => stream.addTrack(track));
    if (screen) {
      screen.getVideoTracks().forEach((track) => stream.addTrack(track));
      screen.getAudioTracks().forEach((track) => stream.addTrack(track));
    }

    return stream;
  }

  onMount(() => {
    refreshDevices();

    const handleDeviceChange = () => {
      refreshDevices();
    };

    navigator.mediaDevices?.addEventListener(
      "devicechange",
      handleDeviceChange,
    );

    onCleanup(() => {
      navigator.mediaDevices?.removeEventListener(
        "devicechange",
        handleDeviceChange,
      );
    });
  });

  onCleanup(() => {
    stopAll();
  });

  return {
    devices,
    selectedCameraId,
    selectedMicId,
    selectedSpeakerId,
    cameraStream,
    micStream,
    screenStream,
    cameraActive,
    micActive,
    screenActive,
    errors,
    supportsSinkID: supportsSinkID(),
    refreshDevices,
    startCamera,
    stopCamera,
    toggleCamera,
    changeCamera,
    startMic,
    stopMic,
    toggleMic,
    changeMic,
    changeSpeaker,
    applySpeaker,
    startScreenShare,
    stopScreenShare,
    toggleScreenShare,
    stopAll,
    getPublishStream,
  };
}

export function createVoiceCall(props) {
  const [status, setStatus] = createSignal("idle");
  const [error, setError] = createSignal("");
  const [participants, setParticipants] = createStore([]);

  const senders = new Map();
  const iceCandidateQueue = [];
  let pc = null;
  let currentChannel = null;
  let makingOffer = false;
  let pendingNegotiation = false;
  let pendingServerOffer = null;
  let iceFlushTimer = null;
  let negotiationQueue = Promise.resolve();

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
    const message = err?.message || fallback;
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
    setParticipants((current) => current.filter((item) => item.id !== Number(id)));
  }

  function setVoiceState(payload) {
    const next = (payload?.participants || [])
      .map((participant) => ({
        id: Number(participant.user_id),
        user: participant.sender,
        label: participant.sender?.username || `User #${participant.user_id}`,
      }))
      .filter((participant) => participant.id && participant.id !== currentUser()?.id);

    setParticipants((current) =>
      next.map((participant) => {
        const existing = current.find((item) => item.id === participant.id);
        return existing ? { ...existing, ...participant } : participant;
      }),
    );
  }

  function flushICECandidates() {
    iceFlushTimer = null;
    if (!iceCandidateQueue.length) return;

    const candidates = iceCandidateQueue.splice(0, iceCandidateQueue.length);
    sendSignal(WSEventType.VoiceICECandidates, { candidates });
  }

  function enqueueICECandidate(candidate) {
    iceCandidateQueue.push(candidate.toJSON());
    if (!iceFlushTimer) {
      iceFlushTimer = setTimeout(flushICECandidates, 50);
    }
  }

  function addRemoteTrack(event) {
    const [eventStream] = event.streams;
    const stream = eventStream || new MediaStream([event.track]);
    const userID = Number(stream.id);
    if (!userID || currentUser()?.id === userID) return;
    const isVideo = event.track.kind === "video";

    upsertParticipant({
      id: userID,
      stream,
      hasVideo: isVideo || stream.getVideoTracks().some((track) => !track.muted && track.readyState === "live"),
    });

    stream.getVideoTracks().forEach((track) => {
      track.onmute = () => {
        upsertParticipant({ id: userID, stream, hasVideo: false });
      };
      track.onunmute = () => {
        upsertParticipant({ id: userID, stream, hasVideo: true });
      };
    });

    event.track.onended = () => {
      upsertParticipant({ id: userID, stream, hasVideo: false });
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
          sendSignal(WSEventType.VoiceOffer, {
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
        existing.generateKeyFrame().catch(() => {});
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

  function createPeerConnection() {
    const connection = new RTCPeerConnection();

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
    sendSignal(WSEventType.VoiceAnswer, {
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
      if (event.type === WSEventType.VoiceAnswer) {
        await handleVoiceAnswer(event.payload);
      } else if (event.type === WSEventType.VoiceOffer) {
        await handleVoiceOffer(event.payload);
      } else if (event.type === WSEventType.VoiceICECandidate) {
        await handleVoiceICECandidate(event.payload);
      } else if (event.type === WSEventType.VoiceICECandidates) {
        await handleVoiceICECandidates(event.payload);
      } else if (event.type === WSEventType.VoiceState) {
        setVoiceState(parsePayload(event.payload));
      } else if (event.type === WSEventType.VoiceJoined) {
        const payload = parsePayload(event.payload);
        upsertParticipant({
          id: payload?.user_id,
          user: payload?.sender,
          label: payload?.sender?.username || `User #${payload?.user_id}`,
        });
      } else if (event.type === WSEventType.VoiceLeft) {
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

    currentChannel = { type: WSChannelType.VoiceTopic, id };
    setStatus("connecting");
    setError("");
    setParticipants([]);

    pc = createPeerConnection();
    appWebSocket.join(currentChannel);
    await syncLocalTracks(false);
    await negotiate();
  }

  function stop() {
    if (currentChannel) {
      appWebSocket.leave(currentChannel);
    }
    senders.clear();
    iceCandidateQueue.splice(0, iceCandidateQueue.length);
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
