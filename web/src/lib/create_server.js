import { action, redirect } from "@solidjs/router";

async function createServer(name) {
  try {
    const res = await fetch("/v1/server/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      return { error: data.message ?? "Failed to update server" };
    }

    return await res.json();
  } catch {
    return { error: "Network error – please check your connection" };
  }
}

export const createServerAction = action(async (formData) => {
  const name = formData.get("name");

  const result = await createServer(name);
  if (result?.error) {
    return { ok: false, error: result.error };
  }

  return redirect(`/server/${result.id}/edit`);
}, "create_server_action");
