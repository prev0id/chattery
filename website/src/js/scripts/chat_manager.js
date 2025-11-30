import { MESSAGE_TYPES } from "./consts.js";

class ChatManager {
  constructor(websocketManager) {
    this.wsManager = websocketManager;
    this.setupSubscriptions();
  }

  setupSubscriptions() {
    this.unsubscribeFunctions = [
      this.wsManager.subscribe(
        MESSAGE_TYPES.ENUM_TYPE_JOIN_CHAT,
        this.handleJoinChat.bind(this),
      ),
      this.wsManager.subscribe(
        MESSAGE_TYPES.ENUM_TYPE_LEAVE_CHAT,
        this.handleLeaveChat.bind(this),
      ),
    ];
  }

  joinChat(chatId) {
    const message = {
      type: MESSAGE_TYPES.ENUM_TYPE_JOIN_CHAT,
      join_chat: { chat_id: chatId },
    };

    const success = this.wsManager.sendMessage(message);
    if (success) {
      console.log("Join chat message sent:", chatId);
    }
    return success;
  }

  leaveChat(chatId) {
    const message = {
      type: MESSAGE_TYPES.ENUM_TYPE_LEAVE_CHAT,
      leave_chat: { chat_id: chatId },
    };

    const success = this.wsManager.sendMessage(message);
    if (success) {
      console.log("Leave chat message sent:", chatId);
    }
    return success;
  }

  handleJoinChat(message) {
    if (message.join_chat) {
      const { chat_id } = message.join_chat;
      console.log(`User joined chat: ${chat_id}`);
      this.onUserJoined?.(chat_id);
    }
  }

  handleLeaveChat(message) {
    if (message.leave_chat) {
      const { chat_id } = message.leave_chat;
      console.log(`User left chat: ${chat_id}`);
      this.onUserLeft?.(chat_id);
    }
  }

  onUserJoined(callback) {
    this.onUserJoined = callback;
  }

  onUserLeft(callback) {
    this.onUserLeft = callback;
  }

  destroy() {
    this.unsubscribeFunctions.forEach((unsub) => unsub());
  }
}

export default ChatManager;
