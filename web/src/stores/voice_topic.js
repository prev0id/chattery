import { createSignal, onCleanup, onMount } from "solid-js";

function stopStream(stream) {
  if (!stream) return;
  stream.getTracks().forEach((track) => track.stop());
}

function videoConstraint(deviceId) {
  return deviceId ? { deviceId: { exact: deviceId } } : true;
}

function audioConstraint(deviceId) {
  return deviceId ? { deviceId: { exact: deviceId } } : true;
}

export function createCallMedia() {
  const [devices, setDevices] = createSignal({
    videoInputs: [],
    audioInputs: [],
  });

  const [selectedCameraId, setSelectedCameraId] = createSignal("");
  const [selectedMicId, setSelectedMicId] = createSignal("");

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
  });

  async function refreshDevices() {
    if (!navigator.mediaDevices?.enumerateDevices) return;

    const all = await navigator.mediaDevices.enumerateDevices();
    const videoInputs = all.filter((d) => d.kind === "videoinput");
    const audioInputs = all.filter((d) => d.kind === "audioinput");

    setDevices({
      videoInputs,
      audioInputs,
    });

    setSelectedCameraId((current) => current || videoInputs[0]?.deviceId || "");
    setSelectedMicId((current) => current || audioInputs[0]?.deviceId || "");
  }

  async function startCamera(deviceId = selectedCameraId()) {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: videoConstraint(deviceId),
        audio: false,
      });

      stopStream(cameraStream());
      setCameraStream(stream);
      setCameraActive(true);
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
    stopStream(cameraStream());
    setCameraStream(null);
    setCameraActive(false);
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
    if (micActive()) {
      await startMic(deviceId);
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

      stopStream(screenStream());
      setScreenStream(stream);
      setScreenActive(true);
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
    stopStream(screenStream());
    setScreenStream(null);
    setScreenActive(false);
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
    if (screen)
      screen.getVideoTracks().forEach((track) => stream.addTrack(track));

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
    cameraStream,
    micStream,
    screenStream,
    cameraActive,
    micActive,
    screenActive,
    errors,
    refreshDevices,
    startCamera,
    stopCamera,
    toggleCamera,
    changeCamera,
    startMic,
    stopMic,
    toggleMic,
    changeMic,
    startScreenShare,
    stopScreenShare,
    toggleScreenShare,
    stopAll,
    getPublishStream,
  };
}
