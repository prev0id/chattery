navigator.mediaDevices
  .getUserMedia({ audio: true, video: true })
  .then((stream) => {
    // Здесь мы можем использовать полученный поток
    console.log("Получен поток:", stream);
  })
  .catch((error) => {
    console.error("Ошибка при получении медиа-потока:", error);
  });

navigator.mediaDevices.enumerateDevices().then((devices) => {
  console.log("Список медиа-устройств:", devices);
});
