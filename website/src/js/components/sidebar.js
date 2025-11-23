import { html } from "../common/common.js";
import "./icon.js";

class Sidebar extends HTMLElement {
  constructor() {
    super();
    this.render();
  }

  render() {}
}

customElements.define("sidebar-element", Sidebar);
