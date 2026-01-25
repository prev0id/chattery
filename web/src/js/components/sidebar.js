import { LitElement, html } from "../lit/min.js";

export class MyElement extends LitElement {
  static properties = {
    count: { type: Number },
  };

  constructor() {
    super();
    this.count = 0;
  }

  render() {
    return html`
      <p><button @click="${this._increment}">Click Me!</button></p>
      <p>Click count: ${this.count}</p>
    `;
  }

  _increment(event) {
    this.count++;
  }
}
customElements.define("my-element", MyElement);
