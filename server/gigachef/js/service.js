'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

/* exported initService */
/* global CHEF gapi ga */

window.CHEF = window.CHEF || {};

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
      this.service = gapi.client.gigachefservice;
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
    }
    /*
     * Chef
     */

  }, {
    key: 'getGigachef',
    value: function getGigachef(callback) {
      var _this = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this.getGigachef(callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.getGigachef(request).execute(function (resp) {
        _this.logError('getGigachef', resp.err);
        callback(resp.gigachef, resp.err);
      });
    }
  }, {
    key: 'updateProfile',
    value: function updateProfile(chef, callback) {
      var _this2 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this2.updateProfile(chef, callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken(),
        gigachef: chef
      };
      this.service.updateProfile(request).execute(function (resp) {
        _this2.logError('updateProfile', resp.err);
        callback(resp.gigachef, resp.err);
        setTimeout(function () {
          _this2.refreshToken();
        }, 1);
      });
    }
    /*
     * Payout Method
     */

  }, {
    key: 'getSubMerchant',
    value: function getSubMerchant(callback) {
      var _this3 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this3.getSubMerchant(callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.getSubMerchant(request).execute(function (resp) {
        _this3.logError('getSubMerchant', resp.err);
        callback(resp.sub_merchant, resp.err);
      });
    }
  }, {
    key: 'updateSubMerchant',
    value: function updateSubMerchant(submerchant, callback) {
      var _this4 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this4.updateSubMerchant(submerchant, callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken(),
        sub_merchant: submerchant
      };
      this.service.updateSubMerchant(request).execute(function (resp) {
        _this4.logError('updateSubMerchant', resp.err);
        callback(resp.gigachef, resp.err);
        setTimeout(function () {
          _this4.refreshToken();
        }, 1);
      });
    }
    /*
     * Item
     */

  }, {
    key: 'getItems',
    value: function getItems(startLimit, endLimit, callback) {
      var _this5 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this5.getItems(startLimit, endLimit, callback);
        });
        return;
      }

      var request = {
        gigatoken: this.getToken(),
        start_limit: startLimit,
        end_limit: endLimit
      };
      this.service.getItems(request).execute(function (resp) {
        _this5.logError('getItems', resp.err);
        callback(resp.items, resp.err);
      });
    }
  }, {
    key: 'getItem',
    value: function getItem(id, callback) {
      var _this6 = this;

      // if api is not loaded, add to _callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this6.getItem(id, callback);
        });
        return;
      }

      var request = {
        id: id,
        gigatoken: this.getToken()
      };
      this.service.getItem(request).execute(function (resp) {
        _this6.logError('getItem', resp.err);
        callback(resp.item, resp.err);
      });
    }
  }, {
    key: 'saveItem',
    value: function saveItem(item, callback) {
      var _this7 = this;

      // if api is not loaded, add to callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this7.saveItem(item, callback);
        });
        return;
      }

      var request = {
        gigatoken: this.getToken(),
        item: item
      };

      this.service.saveItem(request).execute(function (resp) {
        _this7.logError('saveItem', resp.err);
        callback(resp.item, resp.err);
      });
    }
    /*
     * Post
     */

  }, {
    key: 'getPosts',
    value: function getPosts(startLimit, endLimit, callback) {
      var _this8 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this8.getPosts(startLimit, endLimit, callback);
        });
        return;
      }

      var request = {
        gigatoken: this.getToken(),
        start_limit: startLimit,
        end_limit: endLimit
      };
      this.service.getPosts(request).execute(function (resp) {
        _this8.logError('getPosts', resp.err);
        callback(resp.posts, resp.err);
      });
    }
  }, {
    key: 'publishPost',
    value: function publishPost(postReq, callback) {
      var _this9 = this;

      // if api is not loaded, add to _callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this9.postPost(postReq, callback);
        });
        return;
      }

      postReq.gigatoken = this.getToken();
      this.service.publishPost(postReq).execute(function (resp) {
        _this9.logError('publishPost', resp.err);
        callback(resp.post, resp.err);
      });
    }
  }, {
    key: 'refreshToken',
    value: function refreshToken() {
      var _this10 = this;

      // if api is not loaded, add to _callQueue
      if (!this.loaded) {
        this.callQueue.push(this.refreshToken);
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.refreshToken(request).execute(function (resp) {
        if (!_this10.logError('refreshToken', resp.err)) {
          CHEF.User.update(resp.gigatoken);
        }
      });
    }
    /*
     * utils
     */

  }, {
    key: 'getToken',
    value: function getToken() {
      return CHEF.User.token;
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

CHEF.Service = new Service();

function initService() {
  var apiRoot = void 0;
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
  gapi.client.load('gigachefservice', 'v1', function () {
    CHEF.Service.endpointLoaded();
  }, apiRoot);
}