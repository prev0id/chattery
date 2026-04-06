import { query } from "@solidjs/router";
import { createContext, useContext } from "solid-js";
import { fetchServers } from "~/lib/api";

export const ServerContext = createContext();

export function UseServerContext() {
  return useContext(ServerContext);
}

export const GetServers = query(fetchServers, "user_servers");
