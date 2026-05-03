export const entryRoutes = {
  login: () => "/login",
  signup: () => "/signup",
  app: (path = "/dm") => `/app${path}`,
};

export function redirectToLogin() {
  window.location.href = entryRoutes.login();
}

export function redirectToApp(path = "/dm") {
  window.location.href = entryRoutes.app(path);
}
