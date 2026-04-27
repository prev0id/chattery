import { defineConfig } from "vite";
import solid from "vite-plugin-solid";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

export default defineConfig({
  plugins: [tailwindcss(), solid()],

  server: {
    proxy: {
      "/v1": "http://localhost:8080",
      "/ws": {
        target: "http://localhost:8080",
        ws: true,
      },
    },
  },

  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },

  build: {
    target: "esnext",
    rollupOptions: {
      input: {
        login: "login.html",
        signup: "signup.html",
        index: "index.html",
      },
    },
  },
});
