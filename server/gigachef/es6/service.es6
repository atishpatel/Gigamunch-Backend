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
    this.refreshToken();
  }
  /*
   * Chef
   */
  getGigachef(callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.getGigachef(callback);
      });
      return;
    }
    const request = {
      gigatoken: this.getToken(),
    };
    this
      .service
      .getGigachef(request)
      .execute(
        (resp) => {
          this.logError('getGigachef', resp.err);
          callback(resp.gigachef, resp.err);
        }
      );
  }

  updateProfile(chef, callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.updateProfile(chef, callback);
      });
      return;
    }
    const request = {
      gigatoken: this.getToken(),
      gigachef: chef,
    };
    this
      .service
      .updateProfile(request)
      .execute(
        (resp) => {
          this.logError('updateProfile', resp.err);
          callback(resp.gigachef, resp.err);
          setTimeout(() => { this.refreshToken(); }, 1);
        }
      );
  }
  /*
   * Payout Method
   */
  getSubMerchant(callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.getSubMerchant(callback);
      });
      return;
    }
    const request = {
      gigatoken: this.getToken(),
    };
    this
      .service
      .getSubMerchant(request)
      .execute(
        (resp) => {
          this.logError('getSubMerchant', resp.err);
          callback(resp.sub_merchant, resp.err);
        }
      );
  }

  updateSubMerchant(submerchant, callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.updateSubMerchant(submerchant, callback);
      });
      return;
    }
    const request = {
      gigatoken: this.getToken(),
      sub_merchant: submerchant,
    };
    this
      .service
      .updateSubMerchant(request)
      .execute(
        (resp) => {
          this.logError('updateSubMerchant', resp.err);
          callback(resp.gigachef, resp.err);
          setTimeout(() => { this.refreshToken(); }, 1);
        }
      );
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
      .execute(
        (resp) => {
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
          this.logError('getItem', resp.err);
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
          this.logError('saveItem', resp.err);
          callback(resp.item, resp.err);
        }
      );
  }
  /*
   * Post
   */
  getPosts(startLimit, endLimit, callback) {
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.getPosts(startLimit, endLimit, callback);
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
     .getPosts(request)
     .execute(
       (resp) => {
         this.logError('getPosts', resp.err);
         callback(resp.posts, resp.err);
       }
     );
  }

  publishPost(postReq, callback) {
    // if api is not loaded, add to _callQueue
    if (!this.loaded) {
      this.callQueue.push(() => {
        this.postPost(postReq, callback);
      });
      return;
    }

    postReq.gigatoken = this.getToken();
    this
      .service
      .publishPost(postReq)
      .execute(
        (resp) => {
          this.logError('publishPost', resp.err);
          callback(resp.post, resp.err);
        }
      );
  }

  refreshToken() {
    // if api is not loaded, add to _callQueue
    if (!this.loaded) {
      this.callQueue.push(this.refreshToken);
      return;
    }
    const request = {
      gigatoken: this.getToken(),
    };
    this
      .service
      .refreshToken(request)
      .execute(
        (resp) => {
          if (!this.logError('refreshToken', resp.err)) {
            CHEF.User.update(resp.gigatoken);
          }
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
    if (err !== undefined && (err.code === undefined || err.code !== 0)) {
      const desc = `Function: ${fnName} | Message: ${err.message} | Details: ${err.detail}`;
      console.error(desc);
      ga('send', 'exception', {
        exDescription: desc,
        exFatal: false,
      });
      if (err.code !== undefined && err.code === 452) { // code signout
        window.location = '/signout';
      }
      return true;
    }
    return false;
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
