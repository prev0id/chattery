import { useParams, createAsync } from "@solidjs/router";
import { createMemo } from "solid-js";
import SidebarDM from "~/components/SidebarDM";

import Toasts from "~/components/Toast";
import Sidebar from "~/components/Sidebar";
import Button from "~/components/Button";
import { DMContext, GetDMs } from "~/stores/dm";

export default function DMs(props) {
  const params = useParams();
  const dms = createAsync(() => GetDMs());

  const currentDMID = () => parseInt(params.dmID, 10);

  const currentDM = createMemo(() =>
    dms()?.find((dm) => dm.id === currentDMID()),
  );

  return (
    <>
      <Sidebar fallback="Loading DMs...">
        <Button variant="amber" class="mx-4">
          Search users
        </Button>
        <For each={dms()}>{(dm) => <SidebarDM dm={dm} />}</For>
      </Sidebar>
      <main class="flex-1 flex flex-col h-full">
        <DMContext.Provider
          value={{
            currentDM,
            currentDMID,
          }}
        >
          {props.children}
        </DMContext.Provider>
      </main>
      <Toasts />
    </>
  );
}
