import { LitElement, html, css } from "../lit/min.js";
import "./icon.js";
import "./sidebar-menu.js";
import { TabPrivateChats, TabPublicChats } from "./sidebar-menu.js";

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
        .content-bar {
            overflow: auto;
            width: 15rem;
            height: 100%;
            border-right: 1px solid var(--color-border);
            background-color: var(--color-muted-soft);
        }
        .tab-button {
            border: none;
            padding: 0;
            background: none;
            cursor: pointer;
            padding: 0.3rem;
            border-radius: var(--radius-lg);
            color: var(--color-muted);
            margin-top: 1rem;
        }
        .tab-button:hover {
            background-color: var(--color-muted);
        }
        .last {
            margin-top: auto;
        }
    `;

    static properties = {
        tab: { type: String },
    };

    constructor() {
        super();
        this.tab = this.getTabFromURL();
        this.show_settings = false;
        this.addEventListener("sidebar-tab-change", this.handleTabChange);
    }

    handleTabChange(event) {
        const tab = event.detail.tab;
        this.setTabToURL(tab);
        this.tab = tab;
    }

    setTabToURL(tab) {
        const url = new URL(window.location);
        url.searchParams.set("tab", tab);
        window.history.pushState({}, "", url);
    }

    getTabFromURL() {
        const url = new URL(window.location);
        let tab = url.searchParams.get("tab");
        if (tab !== TabPublicChats && tab !== TabPrivateChats) {
            console.log(tab);
            tab = TabPrivateChats;
            this.setTabToURL(tab);
        }
        return tab;
    }

    render() {
        return html`
            <aside>
                <sidebar-menu tab=${this.tab}></sidebar-menu>
            </aside>
        `;
    }
}
customElements.define("chattery-sidebar", Sidebar);
