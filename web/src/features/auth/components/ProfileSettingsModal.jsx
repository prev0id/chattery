import { createSignal } from "solid-js";
import Modal from "~/shared/ui/Modal";
import ProfilePicture from "~/shared/ui/ProfilePicture";
import { userData } from "~/stores/auth";

export function ProfileSettingsModal(props) {
  const [showPassword, setShowPassword] = createSignal(false);

  let avatarFormRef;
  let avatarInputRef;

  const handleAvatarInput = () => {
    avatarFormRef?.submit();
  };

  return (
    <Modal id={props.id} name="Profile Settings">
      <div class="mt-4">
        <h2 class="text-xl font-semibold tracking-wide mb-4">Avatar</h2>

        <div class="flex items-center gap-4 mb-6">
          <ProfilePicture
            src={userData()?.avatar}
            class="size-24 object-cover shrink-0"
          />

          <div class="flex-1 space-y-3">
            <form
              ref={avatarFormRef}
              id="avatar-upload-form"
              enctype="multipart/form-data"
              class="flex items-center"
            >
              <button
                type="button"
                onClick={() => avatarInputRef?.click()}
                class="cursor-pointer px-6 py-2 text-lg font-semibold text-center border-2 neo-shadow rounded-lg bg-emerald-200 hover:bg-emerald-500 focus:outline-none focus:border-emerald-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider flex-1"
              >
                Upload new avatar
              </button>
              <input
                ref={avatarInputRef}
                type="file"
                id="avatar-file-input"
                name="avatar"
                accept="image/*"
                class="hidden"
                onChange={handleAvatarInput}
              />
            </form>

            <form
              id="avatar-delete-form"
              method="POST"
              action="/api/user/avatar/delete"
            >
              <button
                type="submit"
                class="w-full px-6 py-2 text-lg font-semibold text-center border-2 neo-shadow rounded-lg bg-red-200 hover:bg-red-500 focus:outline-none focus:border-red-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider"
              >
                Delete avatar
              </button>
            </form>
          </div>
        </div>
      </div>

      <div class="border-t-2 pt-6">
        <h2 class="text-xl font-semibold tracking-wide mb-4">Update account</h2>

        <form id="profile-update-form">
          <div class="mt-4">
            <label class="block font-semibold" for="username">
              Username
            </label>
            <input
              class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
              id="username"
              name="username"
              type="text"
              placeholder="New username (optional)"
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
              placeholder="New email (optional)"
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
              required
            />
          </div>

          <div class="input-group mt-4">
            <label class="block font-semibold" for="new-password">
              New password
              <span class="text-xs font-normal text-gray-500">
                (leave empty to keep current)
              </span>
            </label>
            <input
              class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
              id="new-password"
              name="newPassword"
              type={showPassword() ? "text" : "password"}
            />
            <button
              type="button"
              onClick={() => setShowPassword((current) => !current)}
              class="block ml-auto mt-1 px-2 border-2 neo-shadow rounded-lg bg-amber-200 hover:bg-amber-500 focus:outline-none focus:border-amber-500 hover:scale-105 transition-all duration-300 ease-in-out"
            >
              {showPassword() ? "Hide" : "Show"}
            </button>
          </div>

          <button
            type="submit"
            class="mt-8 w-full px-4 py-2 text-lg font-semibold border-2 neo-shadow rounded-lg bg-sky-200 hover:bg-sky-500 focus:outline-none focus:border-sky-500 hover:scale-105 transition-all duration-300 ease-in-out tracking-wider"
          >
            Update profile
          </button>
        </form>
      </div>
    </Modal>
  );
}
