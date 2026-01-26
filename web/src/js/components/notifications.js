import { LitElement, html, css } from "../lit/min.js";

class NotificationCenter extends LitElement {
    static styles = css`
        .container {
            position: fixed;
            top: 2rem;
            right: 2rem;
            z-index: 1000;
            display: flex;
            flex-direction: column;
            gap: 1rem;
            max-width: 25rem;
        }
        .notification {
            padding: 0.5rem 1rem;
            border-radius: var(--radius-md);
            color: var(--color-background);
            display: flex;
            justify-content: space-between;
            align-items: center;
            border: 1px solid var(--color-border);
        }
        .info {
            background-color: var(--color-info);
        }
        .error {
            background-color: var(--color-danger);
        }
        .close {
            cursor: pointer;
            margin-left: 1rem;
            font-weight: semibold;
        }
        .status {
            font-weight: semibold;
            margin-right: 0.5rem;
        }
    `;

    constructor() {
        super();
        this.notifications = [];
        this.nextId = 0;
    }

    connectedCallback() {
        super.connectedCallback();
        window.addEventListener(
            "show-notification",
            this.handleNotification.bind(this),
        );
    }

    disconnectedCallback() {
        super.disconnectedCallback();
        window.removeEventListener(
            "show-notification",
            this.handleNotification.bind(this),
        );
    }

    handleNotification(event) {
        const { type = "info", message, status } = event.detail || {};
        const id = this.nextId++;
        const timeoutId = setTimeout(() => this.removeNotification(id), 5000);

        this.notifications = [
            ...this.notifications,
            { id, type, message, status, timeoutId },
        ];
        this.requestUpdate();
    }

    removeNotification(id) {
        const notif = this.notifications.find((n) => n.id === id);
        if (notif) clearTimeout(notif.timeoutId);
        this.notifications = this.notifications.filter((n) => n.id !== id);
        this.requestUpdate();
    }

    render() {
        return html`
            <div class="container">
                ${this.notifications.map(
                    (notif) => html`
                        <div class="notification ${notif.type}">
                            ${notif.status
                                ? html`<span class="status"
                                      >${notif.status}:</span
                                  >`
                                : ""}
                            <span>${notif.message}</span>
                            <span
                                class="close"
                                @click="${() =>
                                    this.removeNotification(notif.id)}"
                                >x</span
                            >
                        </div>
                    `,
                )}
            </div>
        `;
    }
}

customElements.define("notification-center", NotificationCenter);
