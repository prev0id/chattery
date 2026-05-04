import { useNavigate, useParams } from "@solidjs/router";
import {
  createEffect,
  createMemo,
  createResource,
  createSignal,
  For,
  onCleanup,
  Show,
} from "solid-js";
import Button from "~/shared/ui/Button";
import DmSidebarItem from "~/features/dm/components/DmSidebarItem";
import { getDms } from "~/features/dm/api";
import { DM_MESSAGES } from "~/features/dm/constants";
import { DmContext } from "~/features/dm/context";
import { createDmPreviewFromMessage } from "~/features/dm/model";
import { normalizeChatMessage } from "~/features/chat/model";
import { WS_CHANNEL_TYPE } from "~/lib/ws";
import { routes } from "~/shared/config/routes";
import { getUserErrorMessage } from "~/shared/api/errors";
import { parseRouteId } from "~/shared/lib/route";
import AppSidebar from "~/shared/ui/AppSidebar";
import { userData } from "~/shared/stores/auth";
import { appWebSocket } from "~/shared/stores/websocket";

export default function DmLayout(props) {
  const navigate = useNavigate();
  const params = useParams();
  const [loadError, setLoadError] = createSignal("");
  const loadDms = async () => {
    setLoadError("");
    try {
      return await getDms();
    } catch (error) {
      setLoadError(getUserErrorMessage(error, DM_MESSAGES.listFailed));
      return [];
    }
  };
  const [dms, { mutate: mutateDms, refetch: refetchDms }] =
    createResource(loadDms);

  const currentDmId = () => parseRouteId(params.dmId);

  const currentDm = createMemo(() =>
    dms()?.find((dm) => dm.id === currentDmId()),
  );

  const updateDmPreview = (dmId, message) => {
    const userId = userData()?.id;
    const isCurrentDm = dmId === currentDmId();
    const isOwnMessage = message?.sender?.id === userId;

    mutateDms((prev) => {
      if (!prev) return prev;

      const next = prev.map((dm) =>
        dm.id === dmId
          ? {
              ...dm,
              unread: !(isCurrentDm || isOwnMessage),
              message: createDmPreviewFromMessage(message),
            }
          : dm,
      );
      const updated = next.find((dm) => dm.id === dmId);
      if (!updated) return prev;

      return [updated, ...next.filter((dm) => dm.id !== dmId)];
    });
  };

  createEffect(() => {
    const dmId = currentDmId();
    const current = dms()?.find((dm) => dm.id === dmId);
    if (!current?.unread) return;

    mutateDms((prev) =>
      prev?.map((dm) => (dm.id === dmId ? { ...dm, unread: false } : dm)),
    );
  });

  createEffect(() => {
    const unsubscribe = appWebSocket.subscribeMessage(({ channel, payload }) => {
      if (channel?.type !== WS_CHANNEL_TYPE.DM) return;

      const dmId = Number(channel.id);
      const knownDm = dms()?.some((dm) => dm.id === dmId);
      updateDmPreview(dmId, normalizeChatMessage(payload));
      if (!knownDm) {
        refetchDms();
      }
    });

    onCleanup(unsubscribe);
  });

  return (
    <>
      <AppSidebar fallback="Loading DMs...">
        <Button
          variant="amber"
          class="mx-4"
          onClick={() => navigate(routes.dm.search())}
        >
          Search users
        </Button>
        <Show when={loadError()}>
          <p class="rounded-lg bg-red-200 px-3 py-2 text-sm font-semibold text-red-700">
            {loadError()}
          </p>
        </Show>
        <For each={dms()}>{(dm) => <DmSidebarItem dm={dm} />}</For>
      </AppSidebar>
      <main class="flex-1 flex flex-col h-full">
        <DmContext.Provider
          value={{
            currentDm,
            currentDmId,
            refetchDms,
            updateDmPreview,
          }}
        >
          {props.children}
        </DmContext.Provider>
      </main>
    </>
  );
}
