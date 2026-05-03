/**
 * @typedef {Object} DmUser
 * @property {number} id
 * @property {string} username
 * @property {string=} avatar
 */

/**
 * @typedef {Object} DmPreviewMessage
 * @property {string=} date
 * @property {string} content
 */

/**
 * @typedef {Object} Dm
 * @property {number} id
 * @property {DmUser} user
 * @property {boolean} unread
 * @property {DmPreviewMessage=} message
 */

export function normalizeDmList(payload) {
  return payload?.dms ?? [];
}

export function normalizeUserSearchResults(payload) {
  return payload?.users ?? [];
}

export function createDmPreviewFromMessage(message) {
  return {
    date: message?.createdAt,
    content: message?.text ?? "",
  };
}
