import { createContext, useContext } from "solid-js";

export const DmContext = createContext();

export function useDmContext() {
  return useContext(DmContext);
}
