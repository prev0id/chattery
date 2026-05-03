import { Router } from "@solidjs/router";
import AppProviders from "~/app/providers";
import AppRoutes from "~/app/routes";

export default function App() {
  return (
    <AppProviders>
      <Router base="/app">
        <AppRoutes />
      </Router>
    </AppProviders>
  );
}
