// import '../node_modules/@polymer/app-layout/app-drawer-layout/app-drawer-layout.js';
// import '../node_modules/@polymer/app-layout/app-drawer/app-drawer.js';
// import '../node_modules/@polymer/app-layout/app-header-layout/app-header-layout.js';
// import '../node_modules/@polymer/app-layout/app-header/app-header.js';
// import '../node_modules/@polymer/app-layout/app-toolbar/app-toolbar.js';
// import '../node_modules/@polymer/app-route/app-location.js';
// import '../node_modules/@polymer/app-route/app-route.js';
// import '../node_modules/@polymer/iron-pages/iron-pages.js'; 
import './login-page/login-page';
import './logs-page/logs-page';

import { Element as PolymerElement } from '../node_modules/@polymer/polymer/polymer-element.js';
import template from './app-shell.html';
import { GetLog } from './service';


export class AppShell extends PolymerElement {

  // Define a string template instead of a `<template>` element.
  static get template() {
    return template;
  }

  constructor() {
    super();
    this.addEventListener('toast', this.toastEvent.bind(this));
    // TODO: Add error listener?
    if (this.page === undefined) {
      this.page = 'login';
    }
  }

  ready() {
    super.ready();
    GetLog(100).then((resp) => {
      console.log(resp);
    });
  }

  static get properties() {
    return {
      name: {
        Type: String,
      },
      page: {
        type: String,
        reflectToAttribute: true,
        observer: '_pageChanged',
      },
    };
  }

  static get observers() {
    return [
      '_routePageChanged(routeData.page)',
    ];
  }

  _pageChanged(page: String) {
    return;
    this.$.drawer.close(); // close drawer if open
    // Load page import on demand. Show 404 page if fails
    let el;
    let title: String;
    switch (page) {
      case 'login':
        el = this.$.loginPage;
        title = 'Login';
        break;
      default:
        title = page;
    }
    this.set('title', title);

    el.selected();
  }

  toastEvent(e: CustomEvent) {
    const detail = e.detail;
    let text = detail.text;
    if (text === undefined || text === null || text === '') {
      text = detail.message;
    }
    let color = '#323232';
    if (detail.error) {
      color = '#B71C1C';
    }
    if (text && text !== '') {
      this.updateStyles({
        '--paper-toast-background-color': color,
      });
      const duration = detail.duration || 3000;
      this.$.toast.show({
        text,
        duration,
      });
    }
  }
}
customElements.define('app-shell', AppShell);
