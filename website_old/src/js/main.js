import WebSocketManager from "./scripts/websocket_manager.js";
import ChatManager from "./scripts/chat_manager.js";

const wsManager = new WebSocketManager(`ws://${window.location.host}/ws`, {
  reconnectInterval: 3000,
  maxReconnectAttempts: 5,
});

// Set up connection event handlers
wsManager.onConnected(() => {
  console.log("Connected to server");
  document.body.style.backgroundColor = "#e8f5e8";
});

wsManager.onDisconnected(() => {
  console.log("❌ Disconnected from server");
  document.body.style.backgroundColor = "#f5e8e8";
});

wsManager.onError((error) => {
  console.error("WebSocket error:", error);
});

// Create chat manager
const chatManager = new ChatManager(wsManager);

// Set up chat event handlers
chatManager.onUserJoined((chatId) => {
  console.log(`User joined chat: ${chatId}`);
  addMessageToUI(`User joined chat: ${chatId}`);
});

chatManager.onUserLeft((chatId) => {
  console.log(`User left chat: ${chatId}`);
  addMessageToUI(`User left chat: ${chatId}`);
});

function joinChat() {
  const chatId = document.getElementById("chatId").value || "default-chat";
  const success = chatManager.joinChat(chatId);

  if (success) {
    addMessageToUI(`You joined chat: ${chatId}`);
  } else {
    addMessageToUI("Failed to join chat - not connected");
  }
}

function leaveChat() {
  const chatId = document.getElementById("chatId").value || "default-chat";
  const success = chatManager.leaveChat(chatId);

  if (success) {
    addMessageToUI(`You left chat: ${chatId}`);
  } else {
    addMessageToUI("Failed to leave chat - not connected");
  }
}

function addMessageToUI(message) {
  const messagesDiv = document.getElementById("messages");
  const messageElement = document.createElement("div");
  messageElement.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
  messagesDiv.appendChild(messageElement);
  messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// Expose functions to global scope for HTML buttons
window.joinChat = joinChat;
window.leaveChat = leaveChat;

// Clean up when page is closed
window.addEventListener("beforeunload", () => {
  chatManager.destroy();
  wsManager.disconnect();
});
