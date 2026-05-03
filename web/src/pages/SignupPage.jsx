import { createSignal } from "solid-js";
import Toasts from "~/shared/ui/Toasts";
import TextField from "~/shared/ui/TextField";
import Button from "~/shared/ui/Button";
import { createUser } from "~/features/auth/api";
import { AUTH_MESSAGES } from "~/features/auth/constants";
import { getUserErrorMessage } from "~/shared/api/errors";
import { redirectToApp } from "~/shared/config/navigation";

export default function SignupPage() {
  const [username, setUsername] = createSignal("");
  const [login, setLogin] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [repeatPassword, setRepeatPassword] = createSignal("");
  const [showPassword, setShowPassword] = createSignal(false);
  const [showRepeatPassword, setShowRepeatPassword] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(false);
  const [formError, setFormError] = createSignal("");

  const handleSubmit = async (event) => {
    event.preventDefault();
    setFormError("");
    setIsLoading(true);

    if (password() !== repeatPassword()) {
      setFormError(AUTH_MESSAGES.passwordMismatch);
      setIsLoading(false);
      return;
    }

    try {
      await createUser({
        username: username(),
        login: login(),
        password: password(),
      });
      redirectToApp();
    } catch (error) {
      setFormError(getUserErrorMessage(error, AUTH_MESSAGES.signupFailed));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <form
        onSubmit={handleSubmit}
        class="border-3 w-md neo-shadow-lg rounded-xl m-auto p-4 bg-white"
      >
        <div>
          <h1 class="text-2xl font-semibold tracking-wider">Create Account</h1>
          <p>Get started by creating your account.</p>
        </div>

        <div class="mt-4">
          <TextField
            label="Username"
            type="text"
            value={username}
            onInput={(e) => setUsername(e.currentTarget.value)}
            required
          />
        </div>

        <div class="mt-4">
          <TextField
            label="Email"
            type="email"
            value={login}
            onInput={(e) => setLogin(e.currentTarget.value)}
            required
          />
        </div>

        <div class="mt-4">
          <TextField
            label="Password"
            type={showPassword() ? "text" : "password"}
            value={password}
            onInput={(e) => setPassword(e.currentTarget.value)}
            required
          />

          <Button
            type="button"
            variant="amber"
            onClick={() => setShowPassword((prev) => !prev)}
            class="block ml-auto mt-1"
            smallText
          >
            {showPassword() ? "Hide" : "Show"}
          </Button>
        </div>

        <div>
          <TextField
            label="Repeat Password"
            type={showRepeatPassword() ? "text" : "password"}
            value={repeatPassword}
            onInput={(e) => setRepeatPassword(e.currentTarget.value)}
            required
          />
          <Button
            type="button"
            variant="amber"
            onClick={() => setShowRepeatPassword((prev) => !prev)}
            class="block ml-auto mt-1"
            smallText
          >
            {showRepeatPassword() ? "Hide" : "Show"}
          </Button>
        </div>

        <Button type="submit" variant="sky" disabled={isLoading()} class="mb-4">
          {isLoading() ? "Creating..." : "Sign Up"}
        </Button>

        {formError() && (
          <p class="mb-4 rounded-lg bg-red-200 px-4 py-1 text-red-700">
            {formError()}
          </p>
        )}

        <div>
          Already have an account?{" "}
          <a href="/login" class="text-amber-600 hover:text-sky-600">
            Login
          </a>
        </div>
      </form>
      <Toasts />
    </>
  );
}
