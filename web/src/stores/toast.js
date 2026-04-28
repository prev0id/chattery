import { createStore } from "solid-js/store";
import { For, createUniqueId } from "solid-js";

const [toasts, setToasts] = createStore([]);

const removeToast = (id) => {
  setToasts((prev) => prev.filter((t) => t.id !== id));
};

const addToast = (variant, message, data = null) => {
  const id = createUniqueId();
  const newToast = { id, variant, message, data };

  setToasts((prev) => [newToast, ...prev]);

  setTimeout(() => {
    removeToast(id);
  }, 10000);
};

export const toast = {
  info: (message) => addToast("info", message),
  warning: (message) => addToast("warning", message),
  error: (message) => addToast("error", message),
  success: (message) => addToast("success", message),
  dmMessage: (message) => addToast("dm-message", message?.text, message),
};

export { toasts, removeToast };
