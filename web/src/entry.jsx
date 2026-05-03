import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import { createEffect, onCleanup } from "solid-js";

import "./index.css";
import ServerLayout from "./routes/server/ServerLayout";
import ServerTextTopicPage from "./routes/server/ServerTextTopicPage";
import ServerSelectPage from "./routes/server/ServerSelectPage";
import ServerVoiceTopicPage from "./routes/server/ServerVoiceTopicPage";
import ServerEditPage from "./routes/server/ServerEditPage";
import ServerManagePage from "./routes/server/ServerManagePage";
import ServerCreatePage from "./routes/server/ServerCreatePage";

import DmLayout from "./routes/dm/DmLayout";
import DmSelectPage from "./routes/dm/DmSelectPage";
import DmChatPage from "./routes/dm/DmChatPage";
import DmSearchPage from "./routes/dm/DmSearchPage";
import { WSChannelType } from "./lib/ws";
import { toast } from "./stores/toast";
import { userData } from "./stores/auth";
import { appWebSocket } from "./stores/websocket";

const filters = {
  dmId: /^\d+$/,
  serverId: /^\d+$/,
  topicId: /^\d+$/,
};

function GlobalDMNotifications() {
  createEffect(() => {
    const unsubscribe = appWebSocket.subscribeMessage(({ channel, payload }) => {
      if (channel?.type !== WSChannelType.DM) return;
      const currentUser = userData();
      if (!currentUser || payload?.sender?.id === currentUser.id) return;

      const currentDMPath = `/app/dm/${channel.id}`;
      if (window.location.pathname === currentDMPath) return;

      toast.dmMessage(payload);
    });

    onCleanup(unsubscribe);
  });

  return null;
}

render(
  () => (
    <>
      <GlobalDMNotifications />
      <Router base="/app">
        <Route path="/dm" component={DmLayout}>
          <Route path="/" component={DmSelectPage} />
          <Route path="/search" component={DmSearchPage} />
          <Route path="/:dmId" component={DmChatPage} matchFilters={filters} />
        </Route>
        <Route path="/server" component={ServerLayout}>
          <Route path="/" component={ServerSelectPage} />
          <Route path="/create" component={ServerCreatePage} />
          <Route path="/manage" component={ServerManagePage} />
          <Route
            path="/:serverId/text/:topicId"
            component={ServerTextTopicPage}
            matchFilters={filters}
          />
          <Route
            path="/:serverId/voice/:topicId"
            component={ServerVoiceTopicPage}
            matchFilters={filters}
          />
          <Route
            path="/:serverId/edit"
            component={ServerEditPage}
            matchFilters={filters}
          />
        </Route>
      </Router>
    </>
  ),
  document.getElementById("root"),
);
