import { API_ENDPOINTS } from "~/shared/api/endpoints";
import { apiRequest } from "~/shared/api/client";

/**
 * Authenticates a user.
 *
 * @param {{login: string, password: string}} credentials
 * @returns {Promise<Object>}
 * @throws {import("~/shared/api/errors").ApiError}
 * @throws {import("~/shared/api/errors").NetworkError}
 */
export function loginUser({ login, password }) {
  return apiRequest(API_ENDPOINTS.user.login, {
    method: "POST",
    body: JSON.stringify({ login, password }),
  });
}

/**
 * Creates a user account.
 *
 * @param {{username: string, login: string, password: string}} values
 * @returns {Promise<Object>}
 * @throws {import("~/shared/api/errors").ApiError}
 * @throws {import("~/shared/api/errors").NetworkError}
 */
export function createUser({ username, login, password }) {
  return apiRequest(API_ENDPOINTS.user.create, {
    method: "POST",
    body: JSON.stringify({ username, login, password }),
  });
}

export function logoutUser() {
  return fetch(API_ENDPOINTS.user.logout, {
    method: "POST",
    redirect: "manual",
  });
}

export function updateCurrentUser({ username, login, currentPassword, password }) {
  return apiRequest(API_ENDPOINTS.user.update, {
    method: "PUT",
    body: JSON.stringify({ username, login, currentPassword, password }),
  });
}

export function uploadCurrentUserAvatar(file) {
  const formData = new FormData();
  formData.append("image", file);

  return apiRequest(API_ENDPOINTS.image.upload, {
    method: "POST",
    body: formData,
  });
}

export function deleteCurrentUserAvatar() {
  return apiRequest(API_ENDPOINTS.image.delete, {
    method: "DELETE",
  });
}

/**
 * Loads the current authenticated user.
 *
 * @returns {Promise<Object>}
 * @throws {import("~/shared/api/errors").ApiError}
 * @throws {import("~/shared/api/errors").AuthRequiredError}
 * @throws {import("~/shared/api/errors").NetworkError}
 */
export async function getCurrentUser() {
  const data = await apiRequest(API_ENDPOINTS.user.me);
  return data.me;
}
