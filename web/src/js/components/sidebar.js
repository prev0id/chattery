import { LitElement, html, css } from "../lit/min.js";

const tabPrivateChats = "private";
const tabGlobalChats = "global";

export class Sidebar extends LitElement {
    static styles = css`
        *,
        *::before,
        *::after {
            box-sizing: border-box;
        }
        *:not(dialog) {
            margin: 0;
        }
        aside {
            display: flex;
            height: 100%;
        }
        .menu-bar {
            display: flex;
            flex-direction: column;
            align-items: center;
            width: 3rem;
            height: 100%;
            border-right: 1px solid var(--color-border);
            padding-top: 1rem;
        }
        .chat-bar {
            overflow: auto;
            width: 15rem;
            height: 100%;
            border-right: 1px solid var(--color-border);
            background-color: var(--color-muted-soft);
        }
        .logo {
            height: 2.9rem;
            width: 2.9rem;
        }
        .tab-button {
            border: none;
            padding: 0;
            background: none;
            cursor: pointer;
            padding: 0.3rem;
            border-radius: var(--radius-md);
            color: var(--color-muted);
        }
        .tab-button:hover {
            background-color: var(--color-muted);
        }
        .tab-image {
            height: 2rem;
            width: 2rem;
        }
        .last {
            margin-top: auto;
        }
    `;

    static properties = {
        tab: {},
        settings_modal_open: {},
    };

    constructor() {
        super();
        this.settings_modal_open = false;
    }

    render() {
        return html`
            <aside>
                <div class="menu-bar">
                    <img class="logo" src="/src/assets/logo.svg" />
                    <button class="tab-button">
                        <img
                            class="tab-image"
                            src="/src/assets/icons/users.svg"
                        />
                    </button>
                    <button class="tab-button">
                        <img
                            class="tab-image"
                            src="/src/assets/icons/globe.svg"
                        />
                    </button>
                    <button class="tab-button last">
                        <img
                            class="tab-image"
                            src="/src/assets/icons/settings.svg"
                        />
                    </button>
                    <button class="tab-button">
                        <img
                            class="tab-image"
                            src="/src/assets/icons/settings.svg"
                        />
                    </button>
                </div>
                <div class="chat-bar"></div>
            </aside>
        `;
    }
}
customElements.define("chattery-sidebar", Sidebar);
