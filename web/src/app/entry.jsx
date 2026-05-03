import { render } from "solid-js/web";
import App from "~/app/App";
import "~/styles/index.css";

export function mountApp(root = document.getElementById("root")) {
  render(() => <App />, root);
}
