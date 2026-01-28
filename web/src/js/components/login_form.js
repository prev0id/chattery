import { LitElement, html, css } from "../lit/min.js";

class LoginForm extends LitElement {
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
            display: block;
            max-width: 400px;
            padding: 2rem;
            border: 1px solid var(--color-border);
            border-radius: var(--radius-lg);
            background-color: var(--color-muted-soft);
        }
        form {
            display: flex;
            flex-direction: column;
            gap: 1.25rem;
            width: 100%;
        }
        .input-group {
            position: relative;
            display: block;
            min-width: 250px;
        }
        input {
            width: 100%;
            padding: 0.5rem;
            border: 1px solid var(--color-border);
            border-radius: var(--radius-md);
            outline: none;
        }
        input:focus {
            border-color: var(--color-info);
        }
        label {
            display: block;
            margin-bottom: 0.25rem;
        }
        button {
            background-color: var(--color-button);
            color: var(--color-background);
            border: none;
            border-radius: var(--radius-md);
            cursor: pointer;
        }
        button:hover {
            opacity: 0.8;
        }
        button[type="submit"] {
            padding: 0.5rem;
        }
        .toggle-password {
            position: absolute;
            right: 0.5rem;
            top: 72%;
            transform: translateY(-50%);
            padding: 0.25rem;
        }
        .signin-link {
            text-align: center;
            margin-top: 1rem;
            font-size: 0.75rem;
        }
        .signin-link a {
            color: var(--color-info);
            text-decoration: none;
        }
        .signin-link a:hover {
            text-decoration: underline;
        }
        .greetings {
            gap: 0;
        }
        .greetings-comment {
            font-size: 0.75rem;
            text-align: center;
        }
        .greetings-header {
            font-size: 1.1rem;
            text-align: center;
            font-weight: 600;
        }
    `;

    constructor() {
        super();
        this.showPassword = false;
    }

    async _handleSubmit(event) {
        event.preventDefault();
        const form = event.target;
        if (!form.checkValidity()) {
            return;
        }
        const login = form.email.value;
        const password = form.password.value;
        const data = { login, password };

        try {
            const response = await fetch("/v1/user/login", {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            });
            if (!response.ok) {
                const data = await response.json();
                window.dispatchEvent(
                    new CustomEvent("show-notification", {
                        detail: {
                            type: "error",
                            message: data.message,
                            status: response.status,
                        },
                    }),
                );
                return;
            }
            window.location.replace("/app");
        } catch (error) {
            window.dispatchEvent(
                new CustomEvent("show-notification", {
                    detail: { type: "error", message: error },
                }),
            );
        }
    }

    _togglePassword() {
        this.showPassword = !this.showPassword;
        this.requestUpdate();
    }

    render() {
        return html`
            <form @submit="${this._handleSubmit}">
                <div class="greetings">
                    <p class="greetings-header">Welcome Back</p>
                    <p class="greetings-comment">
                        Enter your credentials to access your account.
                    </p>
                </div>
                <div>
                    <label for="email">Email</label>
                    <input id="email" name="login" type="email" required />
                </div>
                <div class="input-group">
                    <label for="password">Password</label>
                    <input
                        id="password"
                        name="password"
                        type="${this.showPassword ? "text" : "password"}"
                        required
                    />
                    <button
                        type="button"
                        class="toggle-password"
                        @click="${this._togglePassword}"
                    >
                        ${this.showPassword ? "Hide" : "Show"}
                    </button>
                </div>
                <button type="submit">Login</button>
            </form>
            <div class="signin-link">
                Don't have an account? <a href="/signup">Sign Up</a>
            </div>
        `;
    }
}

customElements.define("login-form", LoginForm);
