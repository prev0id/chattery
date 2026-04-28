import { useNavigate, useParams } from "@solidjs/router";
import { createEffect, createMemo, createResource } from "solid-js";
import SidebarDM from "~/components/SidebarDM";

import Toasts from "~/components/Toast";
import Sidebar from "~/components/Sidebar";
import Button from "~/components/Button";
import { DMContext } from "~/stores/dm";
import { fetchDMs } from "~/lib/api";

export default function DMs(props) {
  const navigate = useNavigate();
  const params = useParams();
  const [dms, { mutate: mutateDMs, refetch: refetchDMs }] =
    createResource(fetchDMs);

  const currentDMID = () => parseInt(params.dmID, 10);

  const currentDM = createMemo(() =>
    dms()?.find((dm) => dm.id === currentDMID()),
  );

  createEffect(() => {
    mutateDMs((prev) =>
      prev?.map((dm) =>
        dm.id === currentDMID() ? { ...dm, unread: false } : dm,
      ),
    );
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
          }}
        >
          {props.children}
        </DMContext.Provider>
      </main>
      <Toasts />
    </>
  );
}
