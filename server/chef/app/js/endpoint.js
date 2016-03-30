'use strict';

/* global gapi */
/* exported initEndpoints */

function _registerEndpointListeners() {
  var app = document.querySelector('#app');

  app.addEventListener('getApplication', function (e) {
    console.log('getApplication in endpoint file', e);
  });

  // window.addEventListener('saveApplication', () => {
  //
  // });
}

function initEndpoints() {
  var ROOT = undefined;
  if (window.location.host === 'localhost:8080') {
    ROOT = 'localhost:8081';
  } else if (window.location.host === 'gigamunch-omninexus-dev.appspot.com') {
    ROOT = 'https://endpoint-gigamuncher-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
  } else {
    ROOT = 'https://endpoint-gigamuncher-dot-gigamunch-omninexus.appspot.com/_ah/api';
  }
  gapi.client.load('gigachefservice', 'v1', function () {
    _registerEndpointListeners();
  }, ROOT);
}