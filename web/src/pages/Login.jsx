import { createSignal } from "solid-js";
import Toasts from "../components/Toast";
import { toast } from "../stores/toast";
import FormTextInput from "../components/FormTextInput";
import Button from "../components/Button";

export default function Login() {
  const [login, setLogin] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [showPassword, setShowPassword] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setIsLoading(true);

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
        window.location.href = "/app/dm";
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
