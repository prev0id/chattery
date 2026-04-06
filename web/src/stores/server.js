import { createAsync } from "@solidjs/router";
import { fetchServers } from "..//lib/api";

const servers = createAsync(() => fetchServers());
