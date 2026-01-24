import { html } from "../common/common.js";
import "./icon.js";

class Sidebar extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  render() {
    this.innerHTML = html`
      <aside
        class="flex sticky top-0 self-start"
        x-data="{tabs:[{icon:'users', href:'#', selected: true},{icon:'groups', href:'#', selected: false},{icon:'settings', href:'#', selected: false}], groups:[]}"
      >
        <div
          class="flex flex-col items-center w-16 h-screen py-8 bg-white dark:bg-gray-900 dark:border-gray-700"
        >
          <nav class="flex flex-col items-center flex-1 space-y-8">
            <a href="/">
              <img
                class="w-auto h-10"
                src="src/assets/Logo Square.png"
                alt=""
              />
            </a>

            <template x-for="tab in tabs">
              <a
                :href="tab.href"
                :class="tab.selected ? 'text-blue-500 bg-blue-100 dark:text-blue-400 dark:bg-gray-800' : 'text-gray-500 focus:outline-nones dark:text-gray-400 dark:hover:bg-gray-800 hover:bg-gray-100'"
                class="p-1.5 inline-block rounded-lg"
              >
                <icon-element :name="tab.icon"></icon-element>
              </a>
            </template>
          </nav>

          <div class="flex flex-col items-center mt-4 space-y-8">
            <a href="#">
              <img
                class="object-cover w-8 h-8 rounded-lg"
                src="https://lipsum.app/640x480/000/fff"
                alt="avatar"
              />
            </a>
            <a
              href="#"
              class="p-1.5 inline-block text-gray-500 focus:outline-nones rounded-lg dark:text-gray-400 dark:hover:bg-gray-800 hover:bg-gray-100"
            >
              <icon-element name="logout"></icon-element>
            </a>
          </div>
        </div>

        <div
          class="h-screen px-5 py-8 overflow-y-auto bg-white border-l border-r sm:w-64 w-60 dark:bg-gray-900 dark:border-gray-700"
        >
          <div class="relative">
            <span class="absolute inset-y-0 left-0 flex items-center pl-3">
              <icon-element
                name="search"
                class="w-5 h-5 text-gray-400 flex items-center justify-center"
              ></icon-element>
            </span>

            <input
              type="text"
              class="w-full py-1.5 pl-10 pr-4 text-gray-700 bg-white border rounded-md dark:bg-gray-900 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 dark:focus:border-blue-300 focus:ring-blue-300 focus:ring-opacity-40 focus:outline-none focus:ring"
              placeholder="Search"
            />
          </div>
          <nav class="mt-4 -mx-3 space-y-6">
            <div class="space-y-2">
              <label
                class="px-3 text-xs text-gray-500 uppercase dark:text-gray-400"
                >Recent</label
              >
              <button
                class="flex items-center w-full px-5 py-2 bg-gray-100 dark:bg-gray-800 gap-x-2 focus:outline-none"
              >
                <div class="relative">
                  <img
                    class="object-cover w-8 h-8 rounded-full"
                    src="https://lipsum.app/640x480/000/fff"
                    alt=""
                  />
                </div>

                <div class="text-left">
                  <h1
                    class="text-sm font-medium text-gray-700 capitalize dark:text-white"
                  >
                    Jane Doe
                  </h1>
                  <p class="text-xs text-gray-500 dark:text-gray-400">
                    5 unread messages.
                  </p>
                </div>
              </button>

              <button
                class="flex items-center w-full px-5 py-2 gap-x-2 focus:outline-none"
              >
                <div class="relative">
                  <img
                    class="object-cover w-8 h-8 rounded-full"
                    src="https://lipsum.app/640x480/000/fff"
                    alt=""
                  />
                  <span
                    class="h-2 w-2 rounded-full bg-emerald-500 absolute right-0.5 ring-1 ring-white bottom-0"
                  ></span>
                </div>

                <div class="text-left">
                  <h1
                    class="text-sm font-medium text-gray-700 capitalize dark:text-white"
                  >
                    Username
                  </h1>
                </div>
              </button>
            </div>
          </nav>
        </div>
      </aside>
    `;
  }
}

customElements.define("sidebar-element", Sidebar);
