import { defineConfig } from "vite";
import solid from "vite-plugin-solid";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

export default defineConfig({
  plugins: [tailwindcss(), solid()],

  server: {
    proxy: {
      "/v1": "http://localhost:8080",
    },
  },

  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },

  build: {
    target: "esnext",
    // rollupOptions: {
    //   input: {
    //     login: "login.html",
    //     signup: "signup.html",
    //     app: "app.html",
    //     app2: "app2.html",
    //   },
    // },
  },
});
