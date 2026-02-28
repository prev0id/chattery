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
        if (tab === TabPublicChats) {
            url.searchParams.set("tab", tab);
        } else {
            url.searchParams.delete("tab");
        }
        window.history.pushState({}, "", url);
    }

    getTabFromURL() {
        const url = new URL(window.location);
        const tab = url.searchParams.get("tab");
        if (tab === TabPublicChats) {
            return TabPublicChats;
        }
        if (tab !== null && tab !== TabPrivateChats) {
            url.searchParams.delete("tab");
            window.history.replaceState({}, "", url);
        }
        return TabPrivateChats;
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
