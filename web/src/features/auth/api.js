import { API_ENDPOINTS } from "~/shared/api/endpoints";
import { apiRequest } from "~/shared/api/client";

export function loginUser({ login, password }) {
  return apiRequest(API_ENDPOINTS.user.login, {
    method: "POST",
    body: JSON.stringify({ login, password }),
  });
}

export function createUser({ username, login, password }) {
  return apiRequest(API_ENDPOINTS.user.create, {
    method: "POST",
    body: JSON.stringify({ username, login, password }),
  });
}

export async function getCurrentUser() {
  const data = await apiRequest(API_ENDPOINTS.user.me);
  return data.me;
}
