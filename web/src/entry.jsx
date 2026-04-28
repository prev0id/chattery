import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import { createEffect, onCleanup } from "solid-js";

import "./index.css";
import ServerWrapper from "./routes/server/Wrapper";
import TextTopic from "./routes/server/TextTopic";
import SelectServer from "./routes/server/Select";
import VoiceTopic from "./routes/server/VoiceTopic";
import EditServer from "./routes/server/Edit";
import ManageServer from "./routes/server/Manage";
import CreateServer from "./routes/server/Create";

import DMWrapper from "./routes/dm/Wrapper";
import SelectDM from "./routes/dm/Select";
import DM from "./routes/dm/DM";
import SearchDM from "./routes/dm/Search";
import { WSChannelType } from "./lib/ws";
import { toast } from "./stores/toast";
import { userData } from "./stores/auth";
import { appWebSocket } from "./stores/websocket";

const filters = {
  dmID: /^\d+$/,
  serverID: /^\d+$/,
  topicID: /^\d+$/,
};

function GlobalDMNotifications() {
  createEffect(() => {
    const unsubscribe = appWebSocket.subscribeMessage(({ channel, payload }) => {
      if (channel?.type !== WSChannelType.DM) return;
      const currentUser = userData();
      if (!currentUser || payload?.sender?.id === currentUser.id) return;

      toast.info(`${payload?.sender?.username ?? "DM"}: ${payload?.text ?? ""}`);
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
        <Route path="/dm" component={DMWrapper}>
          <Route path="/" component={SelectDM} />
          <Route path="/search" component={SearchDM} />
          <Route path="/:dmID" component={DM} matchFilters={filters} />
        </Route>
        <Route path="/server" component={ServerWrapper}>
          <Route path="/" component={SelectServer} />
          <Route path="/create" component={CreateServer} />
          <Route path="/manage" component={ManageServer} />
          <Route
            path="/:serverID/text/:topicID"
            component={TextTopic}
            matchFilters={filters}
          />
          <Route
            path="/:serverID/voice/:topicID"
            component={VoiceTopic}
            matchFilters={filters}
          />
          <Route
            path="/:serverID/edit"
            component={EditServer}
            matchFilters={filters}
          />
        </Route>
      </Router>
    </>
  ),
  document.getElementById("root"),
);
