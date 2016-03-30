/* global Polymer User */

((function main(document) {
  // Grab a reference to our auto-binding template
  // and give it some initial binding values
  // Learn more about auto-binding templates at http://goo.gl/Dx1u2g
  const app = document.querySelector('#app');

  // Sets app default base URL
  app.baseURL = '/gigachef';
  // Gets user
  app.user = new User();
  // Set up toolbar
  app.toolbar = {};
  //   title: 'Gigamunch',
  //   subtitle: 'Light your inner cook!',
  //   icon: 'menu',
  // };
  // The subpath is for the view, edit, etc subpaths
  app.subpath = '';

  // app.addEventListener('update-toolbar', (e) => {
  //   if (e.detail.title !== undefined && e.detail.title === '') {
  //     app.toolbar.title = e.detail.title;
  //   }
  //   if (e.detail.subtitle !== undefined && e.detail.subtitle === '') {
  //     app.toolbar.subtitle = e.detail.subtitle;
  //   }
  //   if (e.detail.icon !== undefined && e.detail.icon !== '') {
  //     app.toolbar.icon = e.detail.icon;
  //   } else {
  //     app.toolbar.icon = 'menu';
  //   }
  // });

  // Listen for template bound event to know when bindings
  // have resolved and content has been stamped to the page
  app.addEventListener('dom-change', () => {
    // console.log('Our app is ready to rock!');
  });

  window.addEventListener('WebComponentsReady', () => {
    // imports are loaded and elements have been registered
  });

  // Main area's paper-scroll-header-panel custom condensing transformation of
  // the appName in the middle-container and the bottom title in the bottom-container.
  // The appName is moved to top and shrunk on condensing. The bottom sub title
  // is shrunk to nothing on condensing.
  window.addEventListener('paper-header-transform', (e) => {
    const appName = Polymer.dom(document).querySelector('#mainToolbar .app-name');
    const middleContainer = Polymer.dom(document).querySelector('#mainToolbar .middle-container');
    const bottomContainer = Polymer.dom(document).querySelector('#mainToolbar .bottom-container');
    const detail = e.detail;
    const heightDiff = detail.height - detail.condensedHeight;
    const yRatio = Math.min(1, detail.y / heightDiff);
    // appName max size when condensed. The smaller the number the smaller the condensed size.
    const maxMiddleScale = 0.65;
    const auxHeight = heightDiff - detail.y;
    const auxScale = heightDiff / (1 - maxMiddleScale);
    const scaleMiddle = Math.max(maxMiddleScale, auxHeight / auxScale + maxMiddleScale);
    const scaleBottom = 1 - yRatio;

    // Move/translate middleContainer
    Polymer.Base.transform(`translate3d(0,${yRatio * 100}%,0)`, middleContainer);

    // Scale bottomContainer and bottom sub title to nothing and back
    Polymer.Base.transform(`scale(${scaleBottom}) translateZ(0)`, bottomContainer);

    // Scale middleContainer appName
    Polymer.Base.transform(`scale(${scaleMiddle}) translateZ(0)`, appName);
  });

  // Scroll page to top and expand header
  app.scrollPageToTop = function scrollPageToTop() {
    app
      .$
      .headerPanelMain
      .scrollToTop(true);
  };

  app.closeDrawer = function closeDrawer() {
    app
      .$
      .paperDrawerPanel
      .closeDrawer();
  };

  // app.openDrawer = function openDrawer() {
  //   app
  //     .$
  //     .paperDrawerPanel
  //     .openDrawer();
  // };
})(document));
