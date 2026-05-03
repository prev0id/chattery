/**
 * @typedef {Object} ServerTopic
 * @property {number} id
 * @property {string} name
 * @property {"text" | "voice"} type
 */

/**
 * @typedef {Object} Server
 * @property {number} id
 * @property {string} name
 * @property {string=} role
 * @property {ServerTopic[]} topics
 */

export function normalizeServerList(payload) {
  return payload?.servers ?? [];
}
