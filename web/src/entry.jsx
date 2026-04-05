import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";

import "./index.css";
import DMs from "./routes/dm.jsx";
import Servers from "./routes/server";

const filters = {
  dmID: /^\d+$/,
  serverID: /^\d+$/,
  topicID: /^\d+$/,
  topicType: /^(text|voice)$/,
};

render(
  () => (
    <Router base="/app">
      <Route path="/dm/:dmID?" component={DMs} matchFilters={filters} />
      <Route path="/server">
        <Route path="/" component={Servers} matchFilters={filters} />
        <Route
          path="/:serverID/:topicType/:topicID"
          component={Servers}
          matchFilters={filters}
        />
      </Route>
    </Router>
  ),
  document.getElementById("root"),
);
