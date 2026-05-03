import { Route } from "@solidjs/router";
import DmChatPage from "~/routes/dm/DmChatPage";
import DmLayout from "~/routes/dm/DmLayout";
import DmSearchPage from "~/routes/dm/DmSearchPage";
import DmSelectPage from "~/routes/dm/DmSelectPage";
import ServerCreatePage from "~/routes/server/ServerCreatePage";
import ServerEditPage from "~/routes/server/ServerEditPage";
import ServerLayout from "~/routes/server/ServerLayout";
import ServerManagePage from "~/routes/server/ServerManagePage";
import ServerSelectPage from "~/routes/server/ServerSelectPage";
import ServerTextTopicPage from "~/routes/server/ServerTextTopicPage";
import ServerVoiceTopicPage from "~/routes/server/ServerVoiceTopicPage";

export const ROUTE_PARAM_FILTERS = {
  dmId: /^\d+$/,
  serverId: /^\d+$/,
  topicId: /^\d+$/,
};

export default function AppRoutes() {
  return (
    <>
      <Route path="/dm" component={DmLayout}>
        <Route path="/" component={DmSelectPage} />
        <Route path="/search" component={DmSearchPage} />
        <Route
          path="/:dmId"
          component={DmChatPage}
          matchFilters={ROUTE_PARAM_FILTERS}
        />
      </Route>
      <Route path="/server" component={ServerLayout}>
        <Route path="/" component={ServerSelectPage} />
        <Route path="/create" component={ServerCreatePage} />
        <Route path="/manage" component={ServerManagePage} />
        <Route
          path="/:serverId/text/:topicId"
          component={ServerTextTopicPage}
          matchFilters={ROUTE_PARAM_FILTERS}
        />
        <Route
          path="/:serverId/voice/:topicId"
          component={ServerVoiceTopicPage}
          matchFilters={ROUTE_PARAM_FILTERS}
        />
        <Route
          path="/:serverId/edit"
          component={ServerEditPage}
          matchFilters={ROUTE_PARAM_FILTERS}
        />
      </Route>
    </>
  );
}
