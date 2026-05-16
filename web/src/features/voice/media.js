import { batch, createSignal, onCleanup, onMount } from "solid-js";

const storageKeys = {
  camera: "chattery.voice.camera_id",
  mic: "chattery.voice.mic_id",
  speaker: "chattery.voice.speaker_id",
};

const MEDIA_ERROR_MESSAGES = {
  unsupported: "Media devices are not supported by this browser.",
  displayUnsupported: "Screen sharing is not supported by this browser.",
  camera:
    "Unable to start camera. Check camera permissions and device availability.",
  mic: "Unable to start microphone. Check microphone permissions and device availability.",
  screen:
    "Unable to start screen share. Check browser permissions and try again.",
  speaker: "Unable to change output device.",
  deviceLabelsPermission:
    "Allow microphone or camera access to load real device names.",
  deviceLabelsUnavailable:
    "Connect a microphone or camera to load real device names.",
  deviceLabels:
    "Unable to load real device names. Check browser permissions and device availability.",
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
    // allow unavailable localStorage
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

function supportsSinkID() {
  return (
    typeof HTMLMediaElement !== "undefined" &&
    "setSinkId" in HTMLMediaElement.prototype
  );
}

function createMediaError(message) {
  const error = new Error(message);
  error.name = "MediaDeviceError";
  return error;
}

function hiddenDeviceLabel(kind, index) {
  if (kind === "videoinput") {
    return `Camera ${index + 1} (name hidden until access is granted)`;
  }
  if (kind === "audioinput") {
    return `Microphone ${index + 1} (name hidden until access is granted)`;
  }
  return `Speaker ${index + 1} (name hidden until access is granted)`;
}

function hasHiddenLabels(deviceList) {
  return deviceList.some((device) => !device.label);
}

function labelRequestError(err) {
  if (err?.name === "MediaDeviceError") {
    return err.message;
  }
  if (
    err?.name === "NotAllowedError" ||
    err?.name === "PermissionDeniedError" ||
    err?.name === "SecurityError"
  ) {
    return MEDIA_ERROR_MESSAGES.deviceLabelsPermission;
  }
  if (
    err?.name === "NotFoundError" ||
    err?.name === "DevicesNotFoundError" ||
    err?.name === "OverconstrainedError"
  ) {
    return MEDIA_ERROR_MESSAGES.deviceLabelsUnavailable;
  }
  return MEDIA_ERROR_MESSAGES.deviceLabels;
}

function assertUserMediaSupported() {
  if (!navigator.mediaDevices?.getUserMedia) {
    throw createMediaError(MEDIA_ERROR_MESSAGES.unsupported);
  }
}

function assertDisplayMediaSupported() {
  if (!navigator.mediaDevices?.getDisplayMedia) {
    throw createMediaError(MEDIA_ERROR_MESSAGES.displayUnsupported);
  }
}

/**
 * Creates local camera, microphone, screen-share and device-selection state.
 *
 * @returns {Object}
 */
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
  const [deviceLabelsPending, setDeviceLabelsPending] = createSignal(false);
  const [deviceLabelsError, setDeviceLabelsError] = createSignal("");
  let deviceLabelsRequest = null;

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

  async function requestTemporaryDeviceAccess(constraints) {
    const stream = await navigator.mediaDevices.getUserMedia(constraints);
    stopStream(stream);
  }

  async function ensureDeviceLabels() {
    if (deviceLabelsRequest) {
      return deviceLabelsRequest;
    }

    deviceLabelsRequest = (async () => {
      setDeviceLabelsPending(true);
      setDeviceLabelsError("");

      try {
        assertUserMediaSupported();
        await refreshDevices();

        const currentDevices = devices();
        const needsVideo = hasHiddenLabels(currentDevices.videoInputs);
        const needsAudio =
          hasHiddenLabels(currentDevices.audioInputs) ||
          hasHiddenLabels(currentDevices.audioOutputs);

        if (!needsVideo && !needsAudio) {
          return;
        }

        try {
          await requestTemporaryDeviceAccess({
            video: needsVideo,
            audio: needsAudio,
          });
        } catch (err) {
          const shouldRetrySeparately =
            needsVideo &&
            needsAudio &&
            (err?.name === "NotFoundError" ||
              err?.name === "DevicesNotFoundError" ||
              err?.name === "OverconstrainedError");

          if (!shouldRetrySeparately) {
            throw err;
          }

          let recovered = false;

          if (needsAudio) {
            try {
              await requestTemporaryDeviceAccess({
                video: false,
                audio: true,
              });
              recovered = true;
            } catch (audioErr) {
              if (
                audioErr?.name !== "NotFoundError" &&
                audioErr?.name !== "DevicesNotFoundError" &&
                audioErr?.name !== "OverconstrainedError"
              ) {
                throw audioErr;
              }
            }
          }

          if (needsVideo) {
            try {
              await requestTemporaryDeviceAccess({
                video: true,
                audio: false,
              });
              recovered = true;
            } catch (videoErr) {
              if (
                videoErr?.name !== "NotFoundError" &&
                videoErr?.name !== "DevicesNotFoundError" &&
                videoErr?.name !== "OverconstrainedError"
              ) {
                throw videoErr;
              }
            }
          }

          if (!recovered) {
            throw err;
          }
        }

        await refreshDevices();
      } catch (err) {
        setDeviceLabelsError(labelRequestError(err));
      } finally {
        setDeviceLabelsPending(false);
        deviceLabelsRequest = null;
      }
    })();

    return deviceLabelsRequest;
  }

  async function startCamera(deviceId = selectedCameraId()) {
    try {
      assertUserMediaSupported();
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
      await refreshDevices();
      setErrors((prev) => ({ ...prev, camera: "" }));
      setDeviceLabelsError("");
      return stream;
    } catch (err) {
      setErrors((prev) => ({
        ...prev,
        camera:
          err?.name === "MediaDeviceError"
            ? err.message
            : MEDIA_ERROR_MESSAGES.camera,
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
      assertUserMediaSupported();
      const stream = await navigator.mediaDevices.getUserMedia({
        video: false,
        audio: audioConstraint(deviceId),
      });

      stopStream(micStream());
      setMicStream(stream);
      setMicActive(true);
      await refreshDevices();
      setErrors((prev) => ({ ...prev, mic: "" }));
      setDeviceLabelsError("");
      return stream;
    } catch (err) {
      setErrors((prev) => ({
        ...prev,
        mic:
          err?.name === "MediaDeviceError"
            ? err.message
            : MEDIA_ERROR_MESSAGES.mic,
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
        speaker: MEDIA_ERROR_MESSAGES.speaker,
      }));
    }
  }

  async function startScreenShare() {
    try {
      assertDisplayMediaSupported();
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
        screen:
          err?.name === "MediaDeviceError"
            ? err.message
            : MEDIA_ERROR_MESSAGES.screen,
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
    deviceLabelsPending,
    deviceLabelsError,
    getHiddenDeviceLabel: hiddenDeviceLabel,
    supportsSinkID: supportsSinkID(),
    refreshDevices,
    ensureDeviceLabels,
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
