import { useNavigate } from "@solidjs/router";
import { Loader, LogOut, Settings } from "lucide-solid";
import { createSignal, onMount, Show, Suspense } from "solid-js";
import Button from "~/shared/ui/Button";
import { logoutUser } from "~/features/auth/api";
import { ProfileSettingsModal } from "~/features/auth/components/ProfileSettingsModal";
import { AUTH_MESSAGES } from "~/features/auth/constants";
import ProfilePicture from "~/shared/ui/ProfilePicture";
import { redirectToLogin } from "~/shared/config/navigation";
import { routes } from "~/shared/config/routes";
import { userData } from "~/shared/stores/auth";
import { getUserErrorMessage } from "~/shared/api/errors";
import { toast } from "~/shared/stores/toast";

export default function AppSidebar(props) {
  const navigate = useNavigate();
  const [isProfileMenuOpen, setIsProfileMenuOpen] = createSignal(false);
  const [isProfileSettingsOpen, setIsProfileSettingsOpen] = createSignal(false);
  const [isLoggingOut, setIsLoggingOut] = createSignal(false);

  const closeProfileMenu = () => setIsProfileMenuOpen(false);

  const handleProfileMenuFocusOut = (event) => {
    if (event.currentTarget.contains(event.relatedTarget)) return;
    closeProfileMenu();
  };

  const handleLogout = async () => {
    if (isLoggingOut()) return;

    setIsLoggingOut(true);
    try {
      await logoutUser();
      redirectToLogin();
    } catch (error) {
      toast.error(getUserErrorMessage(error, AUTH_MESSAGES.logoutFailed));
      setIsLoggingOut(false);
    }
  };

  const handleOpenProfileSettings = () => {
    setIsProfileSettingsOpen(true);
    closeProfileMenu();
  };

  return (
    <aside class="h-full w-98 flex bg-rose-50">
      <div class="w-18 border-r-3 flex flex-col gap-4 p-4">
        <Button sideways variant="amber" onClick={() => navigate(routes.dm.list())}>
          Direct
        </Button>
        <Button
          sideways
          variant="sky"
          onClick={() => navigate(routes.server.list())}
        >
          Servers
        </Button>

        <div
          class="relative mt-auto"
          onMouseEnter={() => setIsProfileMenuOpen(true)}
          onMouseLeave={closeProfileMenu}
          onFocusIn={() => setIsProfileMenuOpen(true)}
          onFocusOut={handleProfileMenuFocusOut}
        >
          <button
            class="hover:scale-105 transition-all duration-300 ease-in-out"
            type="button"
            aria-label="Open profile settings"
            aria-expanded={isProfileMenuOpen()}
            onClick={() => setIsProfileMenuOpen((current) => !current)}
          >
            <ProfilePicture src={userData()?.avatar} />
          </button>

          <Show when={isProfileMenuOpen()}>
            <>
              <div class="absolute bottom-0 left-full z-10 h-full w-3" />
              <div class="absolute bottom-0 left-full z-20 ml-3 w-36 border-2 neo-shadow rounded-lg bg-white p-0.5">
                <button
                  class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm tracking-wide hover:bg-sky-100 focus:outline-none focus:bg-sky-100"
                  type="button"
                  onClick={handleOpenProfileSettings}
                >
                  <Settings class="size-5" />
                  <span>Settings</span>
                </button>
                <hr class="my-0.5 border-t-2" />
                <button
                  class="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm tracking-wide text-rose-700 hover:bg-rose-100 focus:outline-none focus:bg-rose-100 disabled:opacity-60"
                  type="button"
                  disabled={isLoggingOut()}
                  onClick={handleLogout}
                >
                  <LogOut class="size-5" />
                  <span>Log out</span>
                </button>
              </div>
            </>
          </Show>
        </div>
        <ProfileSettingsModal
          id="profile-settings-popover"
          open={isProfileSettingsOpen()}
          onClose={() => setIsProfileSettingsOpen(false)}
        />
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
