'use strict';

/* global Polymer User */

(function main(document) {
  // Grab a reference to our auto-binding template
  // and give it some initial binding values
  // Learn more about auto-binding templates at http://goo.gl/Dx1u2g
  var app = document.querySelector('#app');

  // Sets app default base URL
  app.baseURL = '/gigachef';
  // Gets user
  app.user = new User();
  // Set up toolbar
  app.title = 'Gigamunch';
  app.subtitle = 'Light your inner cook!';
  app.icon = 'menu';
  // The subpath is for the view, edit, etc subpaths
  app.subpath = '';

  // Listen for template bound event to know when bindings
  // have resolved and content has been stamped to the page
  app.addEventListener('dom-change', function () {
    // console.log('Our app is ready to rock!');
  });

  // imports are loaded and elements have been registered
  window.addEventListener('WebComponentsReady', function () {
    if (app.user.token === undefined || app.user.token === '') {
      window.location = '/login?mode=select';
    }
    app.service = app.$.service;
  });

  // Main area's paper-scroll-header-panel custom condensing transformation of
  // the appName in the middle-container and the bottom title in the bottom-container.
  // The appName is moved to top and shrunk on condensing. The bottom sub title
  // is shrunk to nothing on condensing.
  window.addEventListener('paper-header-transform', function (e) {
    var appName = Polymer.dom(document).querySelector('#mainToolbar .app-name');
    var middleContainer = Polymer.dom(document).querySelector('#mainToolbar .middle-container');
    var bottomContainer = Polymer.dom(document).querySelector('#mainToolbar .bottom-container');
    var detail = e.detail;
    var heightDiff = detail.height - detail.condensedHeight;
    var yRatio = Math.min(1, detail.y / heightDiff);
    // appName max size when condensed. The smaller the number the smaller the condensed size.
    var maxMiddleScale = 0.65;
    var auxHeight = heightDiff - detail.y;
    var auxScale = heightDiff / (1 - maxMiddleScale);
    var scaleMiddle = Math.max(maxMiddleScale, auxHeight / auxScale + maxMiddleScale);
    var scaleBottom = 1 - yRatio;

    // Move/translate middleContainer
    Polymer.Base.transform('translate3d(0,' + yRatio * 100 + '%,0)', middleContainer);

    // Scale bottomContainer and bottom sub title to nothing and back
    Polymer.Base.transform('scale(' + scaleBottom + ') translateZ(0)', bottomContainer);

    // Scale middleContainer appName
    Polymer.Base.transform('scale(' + scaleMiddle + ') translateZ(0)', appName);
  });

  // Scroll page to top and expand header
  app.scrollPageToTop = function scrollPageToTop() {
    app.$.headerPanelMain.scrollToTop(true);
  };

  app.closeDrawer = function closeDrawer() {
    app.$.paperDrawerPanel.closeDrawer();
  };

  app.openDrawer = function openDrawer() {
    app.$.paperDrawerPanel.openDrawer();
  };

  app.toast = function toast(text) {
    app.$.toast.text = text;
    app.$.toast.show();
  };

  app.updateToolbar = function updateToolbar(e) {
    if (e.detail.title !== undefined && e.detail.title === '') {
      app.title = e.detail.title;
    }
    if (e.detail.subtitle !== undefined && e.detail.subtitle === '') {
      app.subtitle = e.detail.subtitle;
    }
    // TODO: add on click call iconCallback
    if (e.detail.icon !== undefined && e.detail.icon !== '') {
      app.icon = e.detail.icon;
    } else {
      app.icon = 'menu';
    }
  };
})(document);