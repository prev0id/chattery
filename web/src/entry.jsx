import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import "./index.css";
import ServerWrapper from "./routes/server/Wrapper";
import TextTopic from "./routes/server/TextTopic";
import SelectServer from "./routes/server/Select";
import VoiceTopic from "./routes/server/VoiceTopic";
import EditServer from "./routes/server/Edit";
import CreateServer from "./routes/server/Create";

import DMWrapper from "./routes/dm/Wrapper";
import SelectDM from "./routes/dm/Select";
import DM from "./routes/dm/DM";

const filters = {
  dmID: /^\d+$/,
  serverID: /^\d+$/,
  topicID: /^\d+$/,
};

render(
  () => (
    <Router base="/app">
      <Route path="/dm" component={DMWrapper}>
        <Route path="/" component={SelectDM} />
        <Route path="/:dmID" component={DM} matchFilters={filters} />
      </Route>
      <Route path="/server" component={ServerWrapper}>
        <Route path="/" component={SelectServer} />
        <Route path="/create" component={CreateServer} />
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
  ),
  document.getElementById("root"),
);
