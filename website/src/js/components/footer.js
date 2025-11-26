import { html } from "../common/common.js";
import "./icon.js";

class Footer extends HTMLElement {
  constructor() {
    super();
    this.render();
  }

  render() {
    this.innerHTML = html`
      <footer class="p-5 border-t bg-white dark:bg-gray-900">
        <div
          class="container flex flex-col items-center justify-between mx-auto space-y-4 sm:space-y-0 sm:flex-row"
        >
          <a href="/">
            <img class="w-auto h-12" src="./src/assets/Fullname.png" alt="" />
          </a>

          <p class="text-sm text-gray-600 dark:text-gray-300">
            Developed by
            <a
              href="https://deev.pro"
              target="_blank"
              class="underline dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400"
            >
              Deev Semyon.
            </a>
            ©MIT Licensed.
          </p>

          <div class="flex -mx-2">
            <a
              href="https://github.com/prev0id/chattery"
              target="_blank"
              class="mx-2 text-gray-600 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400"
              aria-label="Github"
            >
              <icon-element name="github"></icon-element>
            </a>
            <a
              href="https://github.com/prev0id/chattery/blob/master/LICENSE"
              target="_blank"
              class="mx-2 text-gray-600 dark:text-gray-300 hover:text-blue-500 dark:hover:text-blue-400"
              aria-label="Copyright"
            >
              <icon-element name="copyright"></icon-element>
            </a>
          </div>
        </div>
      </footer>
    `;
  }
}

customElements.define("footer-element", Footer);
