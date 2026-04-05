import { createResource } from "solid-js";
import { toast } from "./toast";

async function fetchUserData() {
  try {
    const res = await fetch("/v1/user/me");
    if (res.status === 401) {
      window.location.href = "/login";
      return null;
    }
    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      toast.error(data.message ?? "Failed to load user data");
      return null;
    }
    const data = await res.json();
    return data.me;
  } catch (err) {
    toast.error("Network error – please check your connection");
    return null;
  }
}

export const [userData, { refetch: refetchUserData }] =
  createResource(fetchUserData);
