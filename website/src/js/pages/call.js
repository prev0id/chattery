import { html } from "../common/common.js";

class Call extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  render() {
    this.innerHTML = html`
      <h1 class="text-2xl font-bold mb-4">Call</h1>
      <p>content</p>

      <div class="flex justify-center items-center gap-3 sm:gap-x-5">
        <div
          class="flex bg-white border divide-x rounded-lg dark:bg-gray-900 dark:border-gray-700 dark:divide-gray-700"
        >
          <div x-data="{ tooltip_open: false }" class="relative">
            <button
              @mouseenter="tooltip_open=true"
              @mouseleave="tooltip_open=false"
              class="rounded-l-lg px-4 py-2 font-medium text-gray-600 sm:px-6 dark:hover:bg-gray-800 dark:text-gray-300 hover:bg-gray-100"
            >
              <icon-element name="micro_on"></icon-element>
              <icon-element name="micro_off" class="hidden"></icon-element>
            </button>
            <p
              x-show="tooltip_open"
              class="absolute px-5 py-3 text-center text-gray-600 truncate -translate-x-1/2 bg-white rounded-lg shadow-lg -top-14 left-1/2 dark:shadow-none shadow-gray-200 dark:bg-gray-800 dark:text-white"
            >
              Turn On Microphone
            </p>
          </div>
          <div x-data="{ tooltip_open: false }" class="relative">
            <button
              @mouseenter="tooltip_open=true"
              @mouseleave="tooltip_open=false"
              class="px-4 py-2 font-medium text-gray-600 sm:px-6 dark:hover:bg-gray-800 dark:text-gray-300 hover:bg-gray-100"
            >
              <icon-element name="screenshare_on"></icon-element>
              <icon-element
                name="screenshare_off"
                class="hidden"
              ></icon-element>
              <p
                x-show="tooltip_open"
                class="absolute px-5 py-3 text-center text-gray-600 truncate -translate-x-1/2 bg-white rounded-lg shadow-lg -top-14 left-1/2 dark:shadow-none shadow-gray-200 dark:bg-gray-800 dark:text-white"
              >
                Turn On Screenshare
              </p>
            </button>
          </div>
          <div x-data="{ tooltip_open: false }" class="relative">
            <button
              @mouseenter="tooltip_open=true"
              @mouseleave="tooltip_open=false"
              class="rounded-r-lg px-4 py-2 font-medium text-gray-600 sm:px-6 dark:hover:bg-gray-800 dark:text-gray-300 hover:bg-gray-100"
            >
              <icon-element name="webcam_on"></icon-element>
              <icon-element name="webcam_off" class="hidden"></icon-element>
            </button>
            <p
              x-show="tooltip_open"
              class="absolute px-5 py-3 text-center text-gray-600 truncate -translate-x-1/2 bg-white rounded-lg shadow-lg -top-14 left-1/2 dark:shadow-none shadow-gray-200 dark:bg-gray-800 dark:text-white"
            >
              Turn On Webcam
            </p>
          </div>
        </div>
        <div x-data="{ tooltip_open: false }" class="relative">
          <button
            @mouseenter="tooltip_open=true"
            @mouseleave="tooltip_open=false"
            class="text-blue-600 hover:bg-blue-100 rounded-lg border px-4 py-2"
          >
            <icon-element name="call_settings"></icon-element>
          </button>
          <p
            x-show="tooltip_open"
            class="absolute px-5 py-3 text-center text-gray-600 truncate -translate-x-1/2 bg-white rounded-lg shadow-lg -top-14 left-1/2 dark:shadow-none shadow-gray-200 dark:bg-gray-800 dark:text-white"
          >
            Settings
          </p>
        </div>
        <div x-data="{ tooltip_open: false }" class="relative">
          <button
            @mouseenter="tooltip_open=true"
            @mouseleave="tooltip_open=false"
            class="text-red-600 hover:bg-red-100 rounded-lg border px-4 py-2"
          >
            <icon-element name="leave_call"></icon-element>
          </button>
          <p
            x-show="tooltip_open"
            class="absolute px-5 py-3 text-center text-gray-600 truncate -translate-x-1/2 bg-white rounded-lg shadow-lg -top-14 left-1/2 dark:shadow-none shadow-gray-200 dark:bg-gray-800 dark:text-white"
          >
            Leave Call
          </p>
        </div>
      </div>
    `;
  }
}

customElements.define("call-page", Call);
