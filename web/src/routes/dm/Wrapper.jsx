import { useNavigate, useParams } from "@solidjs/router";
import { createEffect, createMemo, createResource, onCleanup } from "solid-js";
import SidebarDM from "~/components/SidebarDM";

import Toasts from "~/components/Toast";
import Sidebar from "~/components/Sidebar";
import Button from "~/components/Button";
import { DMContext } from "~/stores/dm";
import { fetchDMs } from "~/lib/api";
import { WSChannelType } from "~/lib/ws";
import { userData } from "~/stores/auth";
import { appWebSocket } from "~/stores/websocket";

export default function DMs(props) {
  const navigate = useNavigate();
  const params = useParams();
  const [dms, { mutate: mutateDMs, refetch: refetchDMs }] =
    createResource(fetchDMs);

  const currentDMID = () => parseInt(params.dmID, 10);

  const currentDM = createMemo(() =>
    dms()?.find((dm) => dm.id === currentDMID()),
  );

  const updateDMPreview = (dmID, message) => {
    const userID = userData()?.id;
    const isCurrentDM = dmID === currentDMID();
    const isOwnMessage = message?.sender?.id === userID;

    mutateDMs((prev) => {
      if (!prev) return prev;

      const next = prev.map((dm) =>
        dm.id === dmID
          ? {
              ...dm,
              unread: isCurrentDM || isOwnMessage ? false : true,
              message: {
                date: message?.created_at,
                content: message?.text ?? "",
              },
            }
          : dm,
      );
      const updated = next.find((dm) => dm.id === dmID);
      if (!updated) return prev;

      return [updated, ...next.filter((dm) => dm.id !== dmID)];
    });
  };

  createEffect(() => {
    const dmID = currentDMID();
    const current = dms()?.find((dm) => dm.id === dmID);
    if (!current?.unread) return;

    mutateDMs((prev) =>
      prev?.map((dm) =>
        dm.id === dmID ? { ...dm, unread: false } : dm,
      ),
    );
  });

  createEffect(() => {
    const unsubscribe = appWebSocket.subscribeMessage(({ channel, payload }) => {
      if (channel?.type !== WSChannelType.DM) return;

      const dmID = Number(channel.id);
      const knownDM = dms()?.some((dm) => dm.id === dmID);
      updateDMPreview(dmID, payload);
      if (!knownDM) {
        refetchDMs();
      }
    });

    onCleanup(unsubscribe);
  });

  return (
    <>
      <Sidebar fallback="Loading DMs...">
        <Button
          variant="amber"
          class="mx-4"
          onClick={() => navigate("/dm/search")}
        >
          Search users
        </Button>
        <For each={dms()}>{(dm) => <SidebarDM dm={dm} />}</For>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <DMContext.Provider
          value={{
            currentDM,
            currentDMID,
            refetchDMs,
            updateDMPreview,
          }}
        >
          {props.children}
        </DMContext.Provider>
      </main>
      <Toasts />
    </>
  );
}
