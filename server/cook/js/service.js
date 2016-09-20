'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

/* exported initService */
/* global COOK gapi ga initService */

window.COOK = window.COOK || {};

var Service = function () {
  function Service() {
    _classCallCheck(this, Service);

    this.loaded = false;
    this.callQueue = [];
  }

  _createClass(Service, [{
    key: 'endpointLoaded',
    value: function endpointLoaded() {
      this.loaded = true;
      this.service = gapi.client.cookservice;
      // remove functions from callQueue after calling them.
      if (this.callQueue !== undefined || this.callQueue !== null) {
        var _iteratorNormalCompletion = true;
        var _didIteratorError = false;
        var _iteratorError = undefined;

        try {
          for (var _iterator = this.callQueue[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
            var fn = _step.value;

            fn();
          }
        } catch (err) {
          _didIteratorError = true;
          _iteratorError = err;
        } finally {
          try {
            if (!_iteratorNormalCompletion && _iterator.return) {
              _iterator.return();
            }
          } finally {
            if (_didIteratorError) {
              throw _iteratorError;
            }
          }
        }

        this.callQueue = [];
      }
      this.refreshToken();
      ga('send', {
        hitType: 'timing',
        timingCategory: 'endpoint',
        timingVar: 'load',
        timingValue: window.performance.now()
      });
    }
    /*
     * Cook
     */

  }, {
    key: 'getCook',
    value: function getCook(callback) {
      var _this = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this.getCook(callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.getCook(request).execute(function (resp) {
        _this.logError('getCook', resp.err);
        callback(resp.cook, resp.err);
      });
    }
  }, {
    key: 'updateCook',
    value: function updateCook(cook, callback) {
      var _this2 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this2.updateCook(cook, callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken(),
        cook: cook
      };
      this.service.updateCook(request).execute(function (resp) {
        _this2.logError('updateCook', resp.err);
        callback(resp.cook, resp.err);
        setTimeout(function () {
          _this2.refreshToken();
        }, 1);
      });
    }
    /*
     * Payout Method
     */
    // getSubMerchant(callback) {
    //   if (!this.loaded) {
    //     this.callQueue.push(() => {
    //       this.getSubMerchant(callback);
    //     });
    //     return;
    //   }
    //   const request = {
    //     gigatoken: this.getToken(),
    //   };
    //   this
    //     .service
    //     .getSubMerchant(request)
    //     .execute(
    //       (resp) => {
    //         this.logError('getSubMerchant', resp.err);
    //         callback(resp.sub_merchant, resp.err);
    //       }
    //     );
    // }

    // updateSubMerchant(submerchant, callback) {
    //   if (!this.loaded) {
    //     this.callQueue.push(() => {
    //       this.updateSubMerchant(submerchant, callback);
    //     });
    //     return;
    //   }
    //   const request = {
    //     gigatoken: this.getToken(),
    //     sub_merchant: submerchant,
    //   };
    //   this
    //     .service
    //     .updateSubMerchant(request)
    //     .execute(
    //       (resp) => {
    //         this.logError('updateSubMerchant', resp.err);
    //         callback(resp.gigachef, resp.err);
    //         setTimeout(() => { this.refreshToken(); }, 1);
    //       }
    //     );
    // }
    /*
     * Menu
     */

  }, {
    key: 'getMenus',
    value: function getMenus(callback) {
      var _this3 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this3.getMenus(callback);
        });
        return;
      }

      var request = {
        gigatoken: this.getToken()
      };
      this.service.getMenus(request).execute(function (resp) {
        _this3.logError('getMenus', resp.err);
        callback(resp.menus, resp.err);
      });
    }
    /*
     * Item
     */

  }, {
    key: 'getItem',
    value: function getItem(id, callback) {
      var _this4 = this;

      // if api is not loaded, add to callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this4.getItem(id, callback);
        });
        return;
      }

      var request = {
        id: id,
        gigatoken: this.getToken()
      };
      this.service.getItem(request).execute(function (resp) {
        _this4.logError('getItem', resp.err);
        callback(resp.item, resp.err);
      });
    }
  }, {
    key: 'saveItem',
    value: function saveItem(item, callback) {
      var _this5 = this;

      // if api is not loaded, add to callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this5.saveItem(item, callback);
        });
        return;
      }

      var request = {
        gigatoken: this.getToken(),
        item: item
      };

      this.service.saveItem(request).execute(function (resp) {
        _this5.logError('saveItem', resp.err);
        callback(resp.item, resp.err);
      });
    }
    /*
     * Post
     */
    // getPosts(startLimit, endLimit, callback) {
    //   if (!this.loaded) {
    //     this.callQueue.push(() => {
    //       this.getPosts(startLimit, endLimit, callback);
    //     });
    //     return;
    //   }

    //   const request = {
    //     gigatoken: this.getToken(),
    //     start_limit: startLimit,
    //     end_limit: endLimit,
    //   };
    //   this
    //    .service
    //    .getPosts(request)
    //    .execute(
    //      (resp) => {
    //        this.logError('getPosts', resp.err);
    //        callback(resp.posts, resp.err);
    //      }
    //    );
    // }

    // publishPost(postReq, callback) {
    //   // if api is not loaded, add to _callQueue
    //   if (!this.loaded) {
    //     this.callQueue.push(() => {
    //       this.postPost(postReq, callback);
    //     });
    //     return;
    //   }

    //   postReq.gigatoken = this.getToken();
    //   this
    //     .service
    //     .publishPost(postReq)
    //     .execute(
    //       (resp) => {
    //         this.logError('publishPost', resp.err);
    //         callback(resp.post, resp.err);
    //       }
    //     );
    // }

  }, {
    key: 'refreshToken',
    value: function refreshToken() {
      var _this6 = this;

      // if api is not loaded, add to _callQueue
      if (!this.loaded) {
        this.callQueue.push(this.refreshToken);
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.refreshToken(request).execute(function (resp) {
        if (!_this6.logError('refreshToken', resp.err)) {
          COOK.User.update(resp.gigatoken);
        }
      });
    }
    /*
     * utils
     */

  }, {
    key: 'getToken',
    value: function getToken() {
      return COOK.User.token;
    }
  }, {
    key: 'logError',
    value: function logError(fnName, err) {
      if (err !== undefined && (err.code === undefined || err.code !== 0)) {
        var desc = 'Function: ' + fnName + ' | Message: ' + err.message + ' | Details: ' + err.detail;
        console.error(desc);
        ga('send', 'exception', {
          exDescription: desc,
          exFatal: false
        });
        if (err.code !== undefined && err.code === 452) {
          // code signout
          window.location = '/signout';
        }
        return true;
      }
      return false;
    }
  }]);

  return Service;
}();

window.COOK.Service = new Service();

function initService() {
  var apiRoot = void 0;
  switch (COOK.Env) {
    case COOK.DEV:
      apiRoot = 'http://localhost:8080/_ah/api';
      break;
    case COOK.STAGE:
      apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
      break;
    default:
      apiRoot = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/api';
  }
  gapi.client.load('cookservice', 'v1', function () {
    COOK.Service.endpointLoaded();
  }, apiRoot);
  ga('send', {
    hitType: 'timing',
    timingCategory: 'endpoint',
    timingVar: 'init',
    timingValue: window.performance.now()
  });
}