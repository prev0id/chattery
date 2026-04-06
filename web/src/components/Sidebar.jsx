import Button from "../components/Button";
import ProfilePicture from "../components/ProfilePicture";
import { userData } from "../stores/auth";
import { ProfileSettingsModal } from "../components/ModalProfileSettings";
import { useNavigate } from "@solidjs/router";
import { createSignal, onMount, Show, Suspense } from "solid-js";
import { Loader } from "lucide-solid";

export default function Sidebar(props) {
  const navigate = useNavigate();

  return (
    <aside class="h-full w-98 flex bg-rose-50">
      <div class="w-18 border-r-3 flex flex-col gap-4 p-4">
        <Button sideways variant="amber" onClick={() => navigate("/dm")}>
          Direct
        </Button>
        <Button sideways variant="sky" onClick={() => navigate("/server")}>
          Servers
        </Button>

        <button
          class="mt-auto hover:scale-105 transition-all duration-300 ease-in-out"
          popovertarget="profile-settings-popover"
        >
          <ProfilePicture src={userData()?.avatar} />
        </button>
        <ProfileSettingsModal id="profile-settings-popover" />
      </div>
      <div class="w-80 border-r-3 flex-1 flex flex-col gap-4 p-4 bg-rose-50">
        <Suspense fallback={<LoadingSpinner text={props.fallback} />}>
          {props.children}
        </Suspense>
      </div>
    </aside>
  );
}

function LoadingSpinner(props) {
  const [show, setShow] = createSignal(false);

  onMount(() => {
    const timer = setTimeout(() => setShow(true), 300);
    return () => clearTimeout(timer);
  });

  return (
    <Show when={show()}>
      <div class="mx-auto mt-8 flex items-center gap-2">
        <Loader class="size-5 animate-spin" />
        <span class="tracking-wider text-lg font-semibold">{props.text}</span>
      </div>
    </Show>
  );
}
