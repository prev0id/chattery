import { createEffect, onCleanup } from "solid-js";
import { WS_CHANNEL_TYPE } from "~/lib/ws";
import { entryRoutes } from "~/shared/config/navigation";
import { routes } from "~/shared/config/routes";
import { userData } from "~/shared/stores/auth";
import { toast } from "~/shared/stores/toast";
import { appWebSocket } from "~/shared/stores/websocket";

function GlobalDmNotifications() {
  createEffect(() => {
    const unsubscribe = appWebSocket.subscribeMessage(({ channel, payload }) => {
      if (channel?.type !== WS_CHANNEL_TYPE.DM) return;
      const currentUser = userData();
      if (!currentUser || payload?.sender?.id === currentUser.id) return;

      const currentDmPath = entryRoutes.app(routes.dm.chat(channel.id));
      if (window.location.pathname === currentDmPath) return;

      toast.dmMessage(payload);
    });

    onCleanup(unsubscribe);
  });

  return null;
}

export default function AppProviders(props) {
  return (
    <>
      <GlobalDmNotifications />
      {props.children}
    </>
  );
}
