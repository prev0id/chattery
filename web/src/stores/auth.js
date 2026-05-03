import { createResource } from "solid-js";
import { getCurrentUser } from "~/features/auth/api";
import { AUTH_MESSAGES } from "~/features/auth/constants";
import { AuthRequiredError, getUserErrorMessage } from "~/shared/api/errors";
import { toast } from "./toast";

async function fetchUserData() {
  try {
    return await getCurrentUser();
  } catch (error) {
    if (error instanceof AuthRequiredError) {
      window.location.href = "/login";
      return null;
    }
    toast.error(getUserErrorMessage(error, AUTH_MESSAGES.meFailed));
    return null;
  }
}

export const [userData, { refetch: refetchUserData }] =
  createResource(fetchUserData);
