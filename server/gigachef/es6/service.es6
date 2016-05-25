/* exported initService */
/* global CHEF gapi ga */

window.CHEF = window.CHEF || {};

class Service {
  constructor() {
    this.loaded = false;
    this.callQueue = [];
  }

  endpointLoaded() {
    this.loaded = true;
    this.service = gapi.client.gigachefservice;
    // remove functions from callQueue after calling them.
    if (this.callQueue !== undefined || this.callQueue !== null) {
      for (const fn of this.callQueue) {
        fn();
      }
      this.callQueue = [];
    }
  }
  /*
   * Item
   */
  getItems(startLimit, endLimit, callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.getItems(startLimit, endLimit, callback);
      });
      return;
    }

    const request = {
      gigatoken: this.getToken(),
      start_limit: startLimit,
      end_limit: endLimit,
    };
    this
      .service
      .getItems(request)
      .execute((resp) => {
        this.logError('getItems', resp.err);
        callback(resp.items, resp.err);
      }
      );
  }

  getItem(id, callback) {
    // if api is not loaded, add to _callQueue
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.getItem(id, callback);
      });
      return;
    }

    const request = {
      id,
      gigatoken: this.getToken(),
    };
    this
      .service
      .getItem(request)
      .execute(
        (resp) => {
          this.logError(resp.err);
          callback(resp.item, resp.err);
        }
      );
  }

  saveItem(item, callback) {
    // if api is not loaded, add to callQueue
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.saveItem(item, callback);
      });
      return;
    }

    const request = {
      gigatoken: this.getToken(),
      item,
    };

    this
      .service
      .saveItem(request)
      .execute(
        (resp) => {
          this.logError(resp.err);
          callback(resp.item, resp.err);
        }
      );
  }
  /*
   * Post
   */
  postPost(post, callback) {
    // if api is not loaded, add to _callQueue
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.postPost(post, callback);
      });
      return;
    }

    const request = {
      gigatoken: this.getToken(),
      post,
    };
    this
      .service
      .postPost(request)
      .execute(
        (resp) => {
          this.logError(resp.err);
          callback(resp.post, resp.err);
        }
      );
  }
  /*
   * utils
   */
  getToken() {
    return CHEF.User.token;
  }

  logError(fnName, err) {
    if (err !== undefined && err.code !== 0) {
      const desc = `Function: ${fnName} | Message: ${err.message} | Details: ${err.detail}`;
      console.error(desc);
      ga('send', 'exception', {
        exDescription: desc,
        exFatal: false,
      });
    }
  }
}

CHEF.Service = new Service();

function initService() {
  let apiRoot;
  switch (CHEF.Env) {
    case CHEF.DEV:
      apiRoot = 'http://localhost:8080/_ah/api';
      break;
    case CHEF.STAGE:
      apiRoot = 'https://endpoint-gigachef-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
      break;
    default:
      apiRoot = 'https://endpoint-gigachef-dot-gigamunch-omninexus.appspot.com/_ah/api';
  }
  gapi.client.load('gigachefservice', 'v1', () => { CHEF.Service.endpointLoaded(); }, apiRoot);
}
