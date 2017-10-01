import { Element as PolymerElement } from '../../node_modules/@polymer/polymer/polymer-element.js';
import { GetLogs } from '../service';

export class LogsPage extends PolymerElement {

  static get template() {
    return `
    <h1>{{name}}</h1>
    `;
  }

  constructor() {
    super();
    this.name = 'Logs Page';
  }

  ready() {
    super.ready();
  }

  selected() {
    console.log('logs selected');
    GetLogs(0, 1000).then((resp) => {
      console.log(resp);
    });
  }

  static get properties() {
    return {
      name: {
        Type: String,
      },
    };
  }
}

customElements.define('logs-page', LogsPage);
