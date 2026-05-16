import { API_ENDPOINTS } from "~/shared/api/endpoints";
import { apiRequest } from "~/shared/api/client";

/**
 * @typedef {Object} VoiceIceServer
 * @property {string[]} urls
 * @property {string=} username
 * @property {string=} credential
 */

function mapIceServer(dto) {
  return {
    urls: dto?.urls ?? [],
    username: dto?.username || undefined,
    credential: dto?.credential || undefined,
  };
}

/**
 * Loads ICE servers for the current authenticated voice call.
 *
 * @returns {Promise<VoiceIceServer[]>}
 * @throws {import("~/shared/api/errors").ApiError}
 * @throws {import("~/shared/api/errors").AuthRequiredError}
 * @throws {import("~/shared/api/errors").NetworkError}
 */
export async function getVoiceIceServers() {
  const data = await apiRequest(API_ENDPOINTS.voice.iceServers);
  return (data.ice_servers ?? [])
    .map(mapIceServer)
    .filter((iceServer) => iceServer.urls.length > 0);
}
