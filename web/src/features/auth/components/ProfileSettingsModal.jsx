import { createEffect, createSignal, onCleanup, Show } from "solid-js";
import {
  deleteCurrentUserAvatar,
  updateCurrentUser,
  uploadCurrentUserAvatar,
} from "~/features/auth/api";
import { AUTH_MESSAGES } from "~/features/auth/constants";
import Modal from "~/shared/ui/Modal";
import ProfilePicture from "~/shared/ui/ProfilePicture";
import { getUserErrorMessage } from "~/shared/api/errors";
import {
  bumpAvatarVersion,
  currentUserAvatar,
  refetchUserData,
  userData,
} from "~/shared/stores/auth";
import { toast } from "~/shared/stores/toast";

export function ProfileSettingsModal(props) {
  const [username, setUsername] = createSignal("");
  const [login, setLogin] = createSignal("");
  const [currentPassword, setCurrentPassword] = createSignal("");
  const [newPassword, setNewPassword] = createSignal("");
  const [showPassword, setShowPassword] = createSignal(false);
  const [avatarFile, setAvatarFile] = createSignal(null);
  const [avatarPreview, setAvatarPreview] = createSignal("");
  const [isAvatarPending, setIsAvatarPending] = createSignal(false);
  const [isProfilePending, setIsProfilePending] = createSignal(false);
  const [formError, setFormError] = createSignal("");

  let avatarInputRef;

  createEffect(() => {
    const user = userData();
    if (!user) return;

    setUsername(user.username ?? "");
    setLogin(user.email ?? "");
  });

  onCleanup(() => {
    if (avatarPreview()) {
      URL.revokeObjectURL(avatarPreview());
    }
  });

  const resetAvatarDraft = () => {
    if (avatarPreview()) {
      URL.revokeObjectURL(avatarPreview());
    }
    setAvatarFile(null);
    setAvatarPreview("");
  };

  const handleAvatarInput = (event) => {
    const file = event.currentTarget.files?.[0];
    if (!file || isAvatarPending()) return;

    resetAvatarDraft();
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
    event.currentTarget.value = "";
  };

  const handleSaveAvatar = async () => {
    const file = avatarFile();
    if (!file || isAvatarPending()) return;

    setIsAvatarPending(true);
    try {
      await uploadCurrentUserAvatar(file);
      bumpAvatarVersion();
      await refetchUserData();
      resetAvatarDraft();
      toast.success(AUTH_MESSAGES.avatarUploadSucceeded);
    } catch (error) {
      toast.error(getUserErrorMessage(error, AUTH_MESSAGES.avatarUploadFailed));
    } finally {
      setIsAvatarPending(false);
    }
  };

  const handleDeleteAvatar = async () => {
    if (isAvatarPending()) return;

    setIsAvatarPending(true);
    try {
      await deleteCurrentUserAvatar();
      await refetchUserData();
      bumpAvatarVersion();
      resetAvatarDraft();
      toast.success(AUTH_MESSAGES.avatarDeleteSucceeded);
    } catch (error) {
      toast.error(getUserErrorMessage(error, AUTH_MESSAGES.avatarDeleteFailed));
    } finally {
      setIsAvatarPending(false);
    }
  };

  const handleProfileSubmit = async (event) => {
    event.preventDefault();
    if (isProfilePending()) return;

    setFormError("");
    setIsProfilePending(true);
    try {
      await updateCurrentUser({
        username: username().trim(),
        login: login().trim(),
        currentPassword: currentPassword(),
        password: newPassword(),
      });
      await refetchUserData();
      setCurrentPassword("");
      setNewPassword("");
      toast.success(AUTH_MESSAGES.updateSucceeded);
    } catch (error) {
      const message = getUserErrorMessage(error, AUTH_MESSAGES.updateFailed);
      setFormError(message);
      toast.error(message);
    } finally {
      setIsProfilePending(false);
    }
  };

  return (
    <Modal
      id={props.id}
      name="Profile Settings"
      open={props.open}
      onClose={props.onClose}
    >
      <div class="mt-4">
        <h2 class="text-xl font-semibold tracking-wide mb-4">Avatar</h2>

        <div class="flex items-center gap-4 mb-6">
          <ProfilePicture
            src={avatarPreview() || currentUserAvatar()}
            class="size-24 object-cover shrink-0"
          />

          <div class="flex-1 space-y-3">
            <Show
              when={avatarFile()}
              fallback={
                <>
                  <button
                    type="button"
                    disabled={isAvatarPending()}
                    onClick={() => avatarInputRef?.click()}
                    class="cursor-pointer px-6 py-2 text-lg font-semibold text-center border-2 neo-shadow rounded-lg bg-emerald-200 hover:bg-emerald-500 focus:outline-none focus:border-emerald-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider w-full disabled:opacity-60"
                  >
                    Upload new avatar
                  </button>

                  <button
                    type="button"
                    disabled={isAvatarPending()}
                    onClick={handleDeleteAvatar}
                    class="w-full px-6 py-2 text-lg font-semibold text-center border-2 neo-shadow rounded-lg bg-red-200 hover:bg-red-500 focus:outline-none focus:border-red-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider disabled:opacity-60"
                  >
                    Delete avatar
                  </button>
                </>
              }
            >
              <button
                type="button"
                disabled={isAvatarPending()}
                onClick={handleSaveAvatar}
                class="cursor-pointer px-6 py-2 text-lg font-semibold text-center border-2 neo-shadow rounded-lg bg-emerald-200 hover:bg-emerald-500 focus:outline-none focus:border-emerald-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider w-full disabled:opacity-60"
              >
                Save avatar
              </button>
            </Show>
            <input
              ref={avatarInputRef}
              type="file"
              name="avatar"
              accept="image/*"
              class="hidden"
              onChange={handleAvatarInput}
            />
          </div>
        </div>
      </div>

      <div class="border-t-2 pt-6">
        <h2 class="text-xl font-semibold tracking-wide mb-4">Update account</h2>

        <form id="profile-update-form" onSubmit={handleProfileSubmit}>
          <div class="mt-4">
            <label class="block font-semibold" for="username">
              Username
            </label>
            <input
              class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
              id="username"
              name="username"
              type="text"
              value={username()}
              onInput={(event) => setUsername(event.currentTarget.value)}
              required
            />
          </div>

          <div class="mt-4">
            <label class="block font-semibold" for="login">
              Login / Email
            </label>
            <input
              class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
              id="login"
              name="login"
              type="email"
              value={login()}
              onInput={(event) => setLogin(event.currentTarget.value)}
              required
            />
          </div>

          <div class="mt-4">
            <label class="block font-semibold" for="current-password">
              Current password
            </label>
            <input
              class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
              id="current-password"
              name="currentPassword"
              type="password"
              value={currentPassword()}
              onInput={(event) => setCurrentPassword(event.currentTarget.value)}
            />
          </div>

          <div class="input-group mt-4">
            <label class="block font-semibold" for="new-password">
              New password
            </label>
            <div class="flex gap-2">
              <input
                class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
                id="new-password"
                name="newPassword"
                type={showPassword() ? "text" : "password"}
                value={newPassword()}
                onInput={(event) => setNewPassword(event.currentTarget.value)}
              />
              <button
                type="button"
                onClick={() => setShowPassword((current) => !current)}
                class="block ml-auto px-2 border-2 neo-shadow rounded-lg bg-amber-200 hover:bg-amber-500 focus:outline-none focus:border-amber-500 hover:scale-105 transition-all duration-300 ease-in-out"
                aria-label={showPassword() ? "Hide password" : "Show password"}
              >
                {showPassword() ? "Hide" : "Show"}
              </button>
            </div>
          </div>

          <Show when={formError()}>
            <p class="mt-4 text-sm font-semibold text-rose-700">
              {formError()}
            </p>
          </Show>

          <button
            type="submit"
            disabled={isProfilePending()}
            class="mt-8 w-full px-4 py-2 text-lg font-semibold border-2 neo-shadow rounded-lg bg-sky-200 hover:bg-sky-500 focus:outline-none focus:border-sky-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider disabled:opacity-60"
          >
            Update profile
          </button>
        </form>
      </div>
    </Modal>
  );
}
