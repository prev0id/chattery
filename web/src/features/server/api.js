import { API_ENDPOINTS } from "~/shared/api/endpoints";
import { apiRequest } from "~/shared/api/client";
import { normalizeServerList } from "~/features/server/model";
import { mapChatCursorToDto } from "~/features/chat/model";

/**
 * Loads servers joined by the current user.
 *
 * @returns {Promise<import("~/features/server/model").Server[]>}
 */
export async function getServers() {
  const data = await apiRequest(API_ENDPOINTS.server.list);
  return normalizeServerList(data);
}

/**
 * Searches servers that current user can join.
 *
 * @param {string} query
 * @returns {Promise<import("~/features/server/model").Server[]>}
 */
export async function searchServers(query) {
  const params = new URLSearchParams({ query });
  const data = await apiRequest(`${API_ENDPOINTS.server.search}?${params}`);
  return normalizeServerList(data);
}

/**
 * Joins a server.
 *
 * @param {number} serverId
 * @returns {Promise<Object>}
 */
export function joinServer(serverId) {
  return apiRequest(API_ENDPOINTS.server.join, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId }),
  });
}

/**
 * Leaves a server.
 *
 * @param {number} serverId
 * @returns {Promise<Object>}
 */
export function leaveServer(serverId) {
  return apiRequest(API_ENDPOINTS.server.leave, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId }),
  });
}

/**
 * Creates a server.
 *
 * @param {string} name
 * @returns {Promise<{id: number}>}
 */
export function createServer(name) {
  return apiRequest(API_ENDPOINTS.server.create, {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

/**
 * Updates a server.
 *
 * @param {number} serverId
 * @param {string} name
 * @returns {Promise<Object>}
 */
export function updateServer(serverId, name) {
  return apiRequest(API_ENDPOINTS.server.update, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId, name }),
  });
}

/**
 * Deletes a server.
 *
 * @param {number} serverId
 * @returns {Promise<Object>}
 */
export function deleteServer(serverId) {
  return apiRequest(API_ENDPOINTS.server.delete, {
    method: "DELETE",
    body: JSON.stringify({ server_id: serverId }),
  });
}

/**
 * Creates a topic.
 *
 * @param {number} serverId
 * @param {string} name
 * @param {"text" | "voice"} type
 * @returns {Promise<Object>}
 */
export function createTopic(serverId, name, type) {
  return apiRequest(API_ENDPOINTS.server.createTopic, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId, name, type }),
  });
}

/**
 * Updates a topic.
 *
 * @param {number} topicId
 * @param {string} name
 * @returns {Promise<Object>}
 */
export function updateTopic(topicId, name) {
  return apiRequest(API_ENDPOINTS.server.updateTopic, {
    method: "POST",
    body: JSON.stringify({ topic_id: topicId, name }),
  });
}

/**
 * Deletes a topic.
 *
 * @param {number} topicId
 * @returns {Promise<Object>}
 */
export function deleteTopic(topicId) {
  return apiRequest(API_ENDPOINTS.server.deleteTopic, {
    method: "DELETE",
    body: JSON.stringify({ topic_id: topicId }),
  });
}

/**
 * Loads topic messages.
 *
 * @param {number} topicId
 * @param {import("~/features/chat/model").ChatCursor=} cursor
 * @returns {Promise<Object>}
 */
export function getTopicMessages(topicId, cursor = null) {
  const body = {
    cursor: mapChatCursorToDto("topic_id", topicId, cursor),
  };

  return apiRequest(API_ENDPOINTS.server.topicMessages, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

/**
 * Sends a topic message.
 *
 * @param {number} topicId
 * @param {string} text
 * @returns {Promise<Object>}
 */
export function sendTopicMessage(topicId, text) {
  return apiRequest(API_ENDPOINTS.server.sendTopicMessage, {
    method: "POST",
    body: JSON.stringify({ topic_id: topicId, text }),
  });
}
