import { query } from "@solidjs/router";
import { createContext, useContext } from "solid-js";
import { fetchDMs } from "~/lib/api";

export const DMContext = createContext();

export function UseDMContext() {
  return useContext(DMContext);
}

export const GetDMs = query(fetchDMs, "user_dms");
