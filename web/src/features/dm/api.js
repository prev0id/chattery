import { API_ENDPOINTS } from "~/shared/api/endpoints";
import { apiRequest } from "~/shared/api/client";
import {
  normalizeDmList,
  normalizeUserSearchResults,
} from "~/features/dm/model";

/**
 * Loads direct messages available to the current user.
 *
 * @returns {Promise<import("~/features/dm/model").Dm[]>}
 */
export async function getDms() {
  const data = await apiRequest(API_ENDPOINTS.dm.list);
  return normalizeDmList(data);
}

/**
 * Searches users that can be added to a direct message.
 *
 * @param {string} query
 * @returns {Promise<import("~/features/dm/model").DmUser[]>}
 */
export async function searchDmUsers(query) {
  const params = new URLSearchParams({ query });
  const data = await apiRequest(`${API_ENDPOINTS.dm.searchUsers}?${params}`);
  return normalizeUserSearchResults(data);
}

/**
 * Creates a direct message with another user.
 *
 * @param {number} participantId
 * @returns {Promise<{id: number}>}
 */
export function createDm(participantId) {
  return apiRequest(API_ENDPOINTS.dm.create, {
    method: "POST",
    body: JSON.stringify({ participant_id: participantId }),
  });
}

/**
 * Loads a page of direct message messages.
 *
 * @param {number} dmId
 * @param {{message_id: number, timestamp: string}=} cursor
 * @returns {Promise<{messages: Array, cursor?: Object}>}
 */
export function getDmMessages(dmId, cursor = null) {
  const body = {
    cursor: cursor
      ? {
          dm_id: dmId,
          message_id: cursor.message_id,
          timestamp: cursor.timestamp,
        }
      : { dm_id: dmId },
  };

  return apiRequest(API_ENDPOINTS.dm.messages, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

/**
 * Sends a message to a direct message.
 *
 * @param {number} dmId
 * @param {string} text
 * @returns {Promise<Object>}
 */
export function sendDmMessage(dmId, text) {
  return apiRequest(API_ENDPOINTS.dm.sendMessage, {
    method: "POST",
    body: JSON.stringify({ dm_id: dmId, text }),
  });
}

/**
 * Marks a direct message as read up to a message id.
 *
 * @param {number} dmId
 * @param {number} messageId
 * @returns {Promise<Object>}
 */
export function markDmRead(dmId, messageId) {
  return apiRequest(API_ENDPOINTS.dm.markRead, {
    method: "POST",
    body: JSON.stringify({ dm_id: dmId, message_id: messageId }),
  });
}
