export class ApiError extends Error {
  constructor(message, { status = 0, payload = null } = {}) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

export class AuthRequiredError extends ApiError {
  constructor(message = "Authentication is required") {
    super(message, { status: 401 });
    this.name = "AuthRequiredError";
  }
}

export class NetworkError extends Error {
  constructor(message = "Network error - please check your connection") {
    super(message);
    this.name = "NetworkError";
  }
}

export function getUserErrorMessage(error, fallback = "Request failed") {
  if (error instanceof NetworkError) {
    return error.message || fallback;
  }

  if (error instanceof AuthRequiredError) {
    return "Please sign in to continue";
  }

  if (error instanceof ApiError) {
    return fallback;
  }

  return fallback;
}
