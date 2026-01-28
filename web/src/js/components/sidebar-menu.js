import { LitElement, html, css } from "../lit/min.js";
import "./icon.js";

export const TabPrivateChats = "private";
export const TabPublicChats = "global";

export class SidebarMenu extends LitElement {
    static styles = css`
        *,
        *::before,
        *::after {
            box-sizing: border-box;
        }
        *:not(dialog) {
            margin: 0;
        }
        :host {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            width: 3rem;
            height: 100%;
            border-right: 1px solid var(--color-border);
            padding-top: 1rem;
            padding-bottom: 1rem;
            gap: 1rem;
        }
        .profile-image {
            height: 2.5rem;
            width: 2.5rem;
            border-radius: var(--radius-lg);
            margin-top: auto;
        }
        .logo {
            height: 2.9rem;
            width: 2.9rem;
        }
        .icon {
            height: 2rem;
            width: 2rem;
        }
        .icon-button {
            background: none;
            border: 1px solid var(--color-background);
            border-radius: var(--radius-md);
            cursor: pointer;
            padding: 0.2rem;
        }
        .icon-button:hover {
            background-color: var(--color-muted);
            border: 1px solid var(--color-border);
        }
        .selected {
            background-color: var(--color-muted-soft);
            border: 1px solid var(--color-border);
        }
    `;

    static properties = {
        tab: { type: String, attribute: true },
        show_sidebar: { type: Boolean },
    };

    constructor() {
        super();
        this.tab = TabPrivateChats;
    }

    setTab(tab) {
        if (this.tab == tab) {
            return;
        }
        let tabChange = new CustomEvent("sidebar-tab-change", {
            detail: { tab: tab },
            bubbles: true,
            composed: true,
        });
        this.dispatchEvent(tabChange);
    }

    render() {
        return html`
            <img class="logo" src="/src/assets/logo.svg" />
            <button
                class="icon-button ${this.tab === TabPrivateChats
                    ? "selected"
                    : ""}"
                @click=${() => this.setTab(TabPrivateChats)}
            >
                <chattery-icon
                    name="private_chat"
                    class="icon"
                    size="2rem"
                    color="var(--color-button)"
                ></chattery-icon>
            </button>
            <button
                class="icon-button ${this.tab === TabPublicChats
                    ? "selected"
                    : ""}"
                @click=${() => this.setTab(TabPublicChats)}
            >
                <chattery-icon
                    name="public_chat"
                    class="icon"
                    size="2rem"
                    color="var(--color-button)"
                ></chattery-icon>
            </button>
            <img
                class="profile-image"
                src="https://upload.wikimedia.org/wikipedia/commons/thumb/5/5f/Gurudongmar_Lake_Sikkim%2C_India_%28edit%29.jpg/960px-Gurudongmar_Lake_Sikkim%2C_India_%28edit%29.jpg"
            />
            <button class="icon-button">
                <chattery-icon
                    class="icon"
                    name="settings"
                    size="2rem"
                    color="var(--color-button)"
                ></chattery-icon>
            </button>
        `;
    }
}
customElements.define("sidebar-menu", SidebarMenu);
