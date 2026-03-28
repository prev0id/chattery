import { createSignal } from "solid-js";
import Toasts from "../components/Toast";
import { toast } from "../stores/toast";
import FormTextInput from "../components/FormTextInput";
import Button from "../components/Button";

export default function Signup() {
  const [username, setUsername] = createSignal("");
  const [login, setLogin] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [repeatPassword, setRepeatPassword] = createSignal("");
  const [showPassword, setShowPassword] = createSignal(false);
  const [showRepeatPassword, setShowRepeatPassword] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setIsLoading(true);

    if (password() !== repeatPassword()) {
      toast.error("Passwords do not match");
      setIsLoading(false);
      return;
    }

    try {
      const res = await fetch("/v1/user/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: username(),
          login: login(),
          password: password(),
        }),
      });

      if (res.ok) {
        window.location.href = "/app";
        return;
      }

      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Something went wrong");
    } catch (err) {
      toast.error("Network error – please check your connection");
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
          <FormTextInput
            label="Username"
            type="text"
            value={username}
            onInput={(e) => setUsername(e.currentTarget.value)}
            required
          />
        </div>

        <div class="mt-4">
          <FormTextInput
            label="Email"
            type="email"
            value={login}
            onInput={(e) => setLogin(e.currentTarget.value)}
            required
          />
        </div>

        <div class="mt-4">
          <FormTextInput
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
          <FormTextInput
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
