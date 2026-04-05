import { defineConfig } from "vite";
import solid from "vite-plugin-solid";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss(), solid()],

  server: {
    proxy: {
      "/v1": "http://localhost:8080",
    },
  },

  build: {
    rollupOptions: {
      input: {
        login: "login.html",
        signup: "signup.html",
        app: "app.html",
      },
    },
  },
});
