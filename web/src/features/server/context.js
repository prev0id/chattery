import { createContext, useContext } from "solid-js";

export const ServerContext = createContext();

export function useServerContext() {
  return useContext(ServerContext);
}
