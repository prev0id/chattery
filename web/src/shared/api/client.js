import { AuthRequiredError, ApiError, NetworkError } from "~/shared/api/errors";

async function parseJsonResponse(response) {
  const text = await response.text();
  if (!text) return {};

  try {
    return JSON.parse(text);
  } catch {
    throw new ApiError("Invalid server response", {
      status: response.status,
      payload: { body: text },
    });
  }
}

export async function apiRequest(path, options = {}) {
  let response;
  const headers = options.body instanceof FormData
    ? { ...options.headers }
    : {
        "Content-Type": "application/json",
        ...options.headers,
      };

  try {
    response = await fetch(path, {
      ...options,
      headers,
    });
  } catch {
    throw new NetworkError();
  }

  if (response.status === 401) {
    throw new AuthRequiredError();
  }

  const data = await parseJsonResponse(response);

  if (!response.ok) {
    throw new ApiError(data.message ?? "Request failed", {
      status: response.status,
      payload: data,
    });
  }

  return data;
}
