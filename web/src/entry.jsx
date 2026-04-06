import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import "./index.css";
import DMs from "./routes/dm.jsx";
import Servers from "./routes/server";
import ServerPage from "./routes/ServerWrapper";
import TextTopic from "./routes/TextTopic";
import SelectServer from "./routes/SelectServer";

const filters = {
  dmID: /^\d+$/,
  serverID: /^\d+$/,
  topicID: /^\d+$/,
};

render(
  () => (
    <Router base="/app">
      <Route path="/dm/:dmID?" component={DMs} matchFilters={filters} />
      <Route path="/server">
        <Route path="/" component={Servers} matchFilters={filters} />
        <Route
          path="/:serverID/:topicID"
          component={Servers}
          matchFilters={filters}
        />
        <Route
          path="/:serverID/edit"
          component={Servers}
          matchFilters={filters}
        />
      </Route>
      <Route path="/server2" component={ServerPage}>
        <Route path="/" component={SelectServer} />
        <Route
          path="/:serverID/text/:topicID"
          component={TextTopic}
          matchFilters={filters}
        />
      </Route>
    </Router>
  ),
  document.getElementById("root"),
);
