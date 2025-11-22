import home from "./views/home.js";
import user_chat from "./views/user_chat.js";

const routes = {
  "/": { title: "Home", render: home },
  "/chat": { title: "About", render: user_chat },
};

function router() {
  let view = routes[location.pathname];

  if (view) {
    document.title = view.title;
    app.innerHTML = view.render();
  }
}

window.addEventListener("click", (e) => {
  if (e.target.matches("[data-link]")) {
    e.preventDefault();
    history.pushState("", "", e.target.href);
    router();
  }
});

window.addEventListener("popstate", router);
window.addEventListener("DOMContentLoaded", router);
