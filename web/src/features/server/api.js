import { API_ENDPOINTS } from "~/shared/api/endpoints";
import { apiRequest } from "~/shared/api/client";
import { normalizeServerList } from "~/features/server/model";

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

export function joinServer(serverId) {
  return apiRequest(API_ENDPOINTS.server.join, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId }),
  });
}

export function leaveServer(serverId) {
  return apiRequest(API_ENDPOINTS.server.leave, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId }),
  });
}

export function createServer(name) {
  return apiRequest(API_ENDPOINTS.server.create, {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export function updateServer(serverId, name) {
  return apiRequest(API_ENDPOINTS.server.update, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId, name }),
  });
}

export function deleteServer(serverId) {
  return apiRequest(API_ENDPOINTS.server.delete, {
    method: "DELETE",
    body: JSON.stringify({ server_id: serverId }),
  });
}

export function createTopic(serverId, name, type) {
  return apiRequest(API_ENDPOINTS.server.createTopic, {
    method: "POST",
    body: JSON.stringify({ server_id: serverId, name, type }),
  });
}

export function updateTopic(topicId, name) {
  return apiRequest(API_ENDPOINTS.server.updateTopic, {
    method: "POST",
    body: JSON.stringify({ topic_id: topicId, name }),
  });
}

export function deleteTopic(topicId) {
  return apiRequest(API_ENDPOINTS.server.deleteTopic, {
    method: "DELETE",
    body: JSON.stringify({ topic_id: topicId }),
  });
}

export function getTopicMessages(topicId, cursor = null) {
  const body = {
    cursor: cursor
      ? {
          topic_id: topicId,
          message_id: cursor.message_id,
          timestamp: cursor.timestamp,
        }
      : { topic_id: topicId },
  };

  return apiRequest(API_ENDPOINTS.server.topicMessages, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function sendTopicMessage(topicId, text) {
  return apiRequest(API_ENDPOINTS.server.sendTopicMessage, {
    method: "POST",
    body: JSON.stringify({ topic_id: topicId, text }),
  });
}
