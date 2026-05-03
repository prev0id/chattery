/**
 * @typedef {Object} ChatCursor
 * @property {number} messageId
 * @property {string} timestamp
 */

/**
 * @typedef {Object} ChatMessage
 * @property {number} id
 * @property {string} text
 * @property {string} createdAt
 * @property {{id: number, username: string, avatar?: string}} sender
 */

export function normalizeChatMessage(message) {
  return {
    ...message,
    id: Number(message?.id),
    text: message?.text ?? "",
    createdAt: message?.created_at ?? message?.createdAt ?? "",
    sender: message?.sender,
  };
}

export function normalizeChatCursor(cursor) {
  if (!cursor) return null;

  return {
    messageId: Number(cursor.message_id ?? cursor.messageId),
    timestamp: cursor.timestamp,
  };
}

export function normalizeChatPage(response) {
  return {
    messages: response?.messages?.map(normalizeChatMessage) ?? [],
    cursor: normalizeChatCursor(response?.cursor),
  };
}

export function mapChatCursorToDto(chatIdKey, chatId, cursor = null) {
  if (!cursor) {
    return { [chatIdKey]: chatId };
  }

  return {
    [chatIdKey]: chatId,
    message_id: cursor.messageId,
    timestamp: cursor.timestamp,
  };
}
