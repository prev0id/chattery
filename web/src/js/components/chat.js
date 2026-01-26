import { LitElement, html } from "../lit/min.js";

export class Chat extends LitElement {
    static properties = {};

    constructor() {
        super();
    }

    render() {
        return html` <main>main</main> `;
    }
}
customElements.define("chattery-chat", Chat);
