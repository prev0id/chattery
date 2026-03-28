import { createSignal } from "solid-js";
import Toast from "../components/Toast";
import FormTextInput from "../components/FormTextInput";
import Button from "../components/Button";

export default function Login() {
  const [login, setLogin] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [showPassword, setShowPassword] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal(null); // { status?: number, message?: string }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const res = await fetch("/v1/user/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          login: login(),
          password: password(),
        }),
      });

      if (res.ok) {
        window.location.href = "/app";
        return;
      }

      let message = "Something went wrong";
      try {
        const data = await res.json();
        message = data.message || message;
      } catch (_) {}

      setError({
        status: res.status,
        message,
      });
    } catch (err) {
      setError({
        status: null,
        message: "Network error – please check your connection",
      });
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
          <FormTextInput
            label="Email"
            type="email"
            value={login}
            onInput={(event) => setLogin(event.currentTarget.value)}
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
          {isLoading() ? "Logging in..." : "Login"}
        </Button>

        <div>
          Don't have an account?{" "}
          <a href="/signup" class="text-amber-500 hover:text-sky-500">
            Sign Up
          </a>
        </div>
      </form>

      {error() && (
        <Toast
          status={error().status}
          message={error().message}
          onClose={() => setError(null)}
        />
      )}
    </>
  );
}
