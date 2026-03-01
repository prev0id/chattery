import { LitElement, html, css } from "../lit/min.js";
import { TabPrivateChats, TabPublicChats } from "./sidebar-menu.js";

export class SidebarContent extends LitElement {
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
      width: 100%;
      height: 100%;
      overflow: hidden;
    }
    .header {
      padding: 1rem;
      border-bottom: 1px solid var(--color-border);
      font-weight: 600;
      font-size: 0.9rem;
    }
    .chat-list {
      flex: 1;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
    }
    .chat-item {
      display: flex;
      flex-direction: column;
      padding: 0.75rem 1rem;
      border-bottom: 1px solid var(--color-muted);
      cursor: pointer;
      gap: 0.2rem;
    }
    .chat-item:hover {
      background-color: var(--color-muted-soft);
    }
    .chat-item.selected {
      background-color: var(--color-info-soft);
    }
    .chat-name {
      font-size: 0.875rem;
      font-weight: 500;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .chat-preview {
      font-size: 0.75rem;
      color: #666;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .chat-preview-time {
      font-size: 0.7rem;
      color: #999;
      align-self: flex-end;
      flex-shrink: 0;
    }
    .chat-preview-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 0.5rem;
    }
    .loading {
      padding: 1.5rem 1rem;
      font-size: 0.8rem;
      color: #999;
      text-align: center;
    }
    .empty {
      padding: 1.5rem 1rem;
      font-size: 0.8rem;
      color: #999;
      text-align: center;
    }
  `;

  static properties = {
    tab: { type: String, attribute: true },
    selectedChatId: { type: Number },
    _chats: { type: Array, state: true },
    _loading: { type: Boolean, state: true },
  };

  constructor() {
    super();
    this.tab = TabPrivateChats;
    this.selectedChatId = null;
    this._chats = [];
    this._loading = false;
    this._abortController = null;
  }

  updated(changedProperties) {
    if (changedProperties.has("tab")) {
      this._fetchChats();
    }
  }

  connectedCallback() {
    super.connectedCallback();
    this._fetchChats();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    if (this._abortController) {
      this._abortController.abort();
    }
  }

  async _fetchChats() {
    if (this._abortController) {
      this._abortController.abort();
    }
    this._abortController = new AbortController();

    this._loading = true;
    this._chats = [];

    const endpoint =
      this.tab === TabPublicChats
        ? "/v1/chat/list/public"
        : "/v1/chat/list/private";

    try {
      const response = await fetch(endpoint, {
        signal: this._abortController.signal,
      });

      if (!response.ok) {
        const data = await response.json().catch(() => ({}));
        window.dispatchEvent(
          new CustomEvent("show-notification", {
            detail: {
              type: "error",
              message: data.message || "Failed to load chats",
              status: response.status,
            },
          }),
        );
        return;
      }

      const data = await response.json();
      this._chats = data.chats || [];
    } catch (error) {
      if (error.name === "AbortError") return;
      window.dispatchEvent(
        new CustomEvent("show-notification", {
          detail: {
            type: "error",
            message: error.message || "Failed to load chats",
          },
        }),
      );
    } finally {
      this._loading = false;
    }
  }

  _selectChat(chat) {
    this.selectedChatId = chat.id;
    this.dispatchEvent(
      new CustomEvent("chat-selected", {
        detail: { chat },
        bubbles: true,
        composed: true,
      }),
    );
  }

  _formatTime(isoString) {
    if (!isoString) return "";
    const date = new Date(isoString);
    const now = new Date();
    const isToday = date.toDateString() === now.toDateString();
    if (isToday) {
      return date.toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    }
    return date.toLocaleDateString([], { month: "short", day: "numeric" });
  }

  _renderChat(chat) {
    const hasPreview = chat.last_message && chat.last_message.text;
    const isSelected = this.selectedChatId === chat.id;

    return html`
      <div
        class="chat-item ${isSelected ? "selected" : ""}"
        @click=${() => this._selectChat(chat)}
      >
        <div class="chat-name">${chat.name}</div>
        ${hasPreview
          ? html`
              <div class="chat-preview-row">
                <span class="chat-preview">${chat.last_message.text}</span>
                <span class="chat-preview-time"
                  >${this._formatTime(chat.last_message.created_at)}</span
                >
              </div>
            `
          : ""}
      </div>
    `;
  }

  render() {
    const title =
      this.tab === TabPublicChats ? "Public Chats" : "Private Chats";

    return html`
      <div class="header">${title}</div>
      <div class="chat-list">
        ${this._loading
          ? html`<div class="loading">Loading...</div>`
          : this._chats.length === 0
            ? html`<div class="empty">No chats yet</div>`
            : this._chats.map((chat) => this._renderChat(chat))}
      </div>
    `;
  }
}

customElements.define("sidebar-content", SidebarContent);
