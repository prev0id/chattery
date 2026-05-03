import { createSignal } from "solid-js";
import Toasts from "~/shared/ui/Toasts";
import { toast } from "../stores/toast";
import TextField from "~/shared/ui/TextField";
import Button from "~/shared/ui/Button";
import { loginUser } from "~/features/auth/api";
import { AUTH_MESSAGES } from "~/features/auth/constants";
import { getUserErrorMessage } from "~/shared/api/errors";

export default function Login() {
  const [login, setLogin] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [showPassword, setShowPassword] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setIsLoading(true);

    try {
      await loginUser({
        login: login(),
        password: password(),
      });
      window.location.href = "/app/dm";
    } catch (error) {
      toast.error(getUserErrorMessage(error, AUTH_MESSAGES.loginFailed));
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
          <h1 class="text-2xl font-semibold tracking-wider">Welcome Back</h1>
          <p>Enter your credentials to access your account.</p>
        </div>

        <div class="mt-4">
          <TextField
            label="Email"
            type="email"
            value={login}
            onInput={(event) => setLogin(event.currentTarget.value)}
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
            autocomplete="off"
          />
          <Button
            type="button"
            variant="amber"
            onClick={() => setShowPassword((prev) => !prev)}
            class={"block ml-auto mt-1"}
            smallText
          >
            {showPassword() ? "Hide" : "Show"}
          </Button>
        </div>

        <Button
          type="submit"
          variant="sky"
          disabled={isLoading()}
          class={"mb-4"}
        >
          Login
        </Button>

        <div>
          Don't have an account?{" "}
          <a href="/signup" class="text-amber-600 hover:text-sky-600">
            Sign Up
          </a>
        </div>
      </form>
      <Toasts />
    </>
  );
}
