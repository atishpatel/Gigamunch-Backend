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
     * Onboard
     */

  }, {
    key: 'finishOnboarding',
    value: function finishOnboarding(cook, submerchant, callback) {
      var _this = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this.finishOnboarding(cook, submerchant, callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken(),
        cook: cook,
        sub_merchant: submerchant
      };
      this.service.finishOnboarding(request).execute(function (resp) {
        _this.logError('finishOnboarding', resp.err);
        COOK.User.update(resp.gigatoken);
        callback(resp.err);
      });
    }

    /*
     * Cook
     */

  }, {
    key: 'getCook',
    value: function getCook(callback) {
      var _this2 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this2.getCook(callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.getCook(request).execute(function (resp) {
        _this2.logError('getCook', resp.err);
        callback(resp.cook, resp.err);
      });
    }
  }, {
    key: 'updateCook',
    value: function updateCook(cook, callback) {
      var _this3 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this3.updateCook(cook, callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken(),
        cook: cook
      };
      this.service.updateCook(request).execute(function (resp) {
        _this3.logError('updateCook', resp.err);
        callback(resp.cook, resp.err);
        setTimeout(function () {
          _this3.refreshToken();
        }, 1);
      });
    }

    /*
     * Payout Method
     */

  }, {
    key: 'getSubMerchant',
    value: function getSubMerchant(callback) {
      var _this4 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this4.getSubMerchant(callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.getSubMerchant(request).execute(function (resp) {
        _this4.logError('getSubMerchant', resp.err);
        callback(resp.sub_merchant, resp.err);
      });
    }
  }, {
    key: 'updateSubMerchant',
    value: function updateSubMerchant(submerchant, callback) {
      var _this5 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this5.updateSubMerchant(submerchant, callback);
        });
        return;
      }
      var request = {
        gigatoken: this.getToken(),
        sub_merchant: submerchant
      };
      this.service.updateSubMerchant(request).execute(function (resp) {
        _this5.logError('updateSubMerchant', resp.err);
        callback(resp.cook, resp.err);
        setTimeout(function () {
          _this5.refreshToken();
        }, 1);
      });
    }

    /*
     * Menu
     */

  }, {
    key: 'getMenus',
    value: function getMenus(callback) {
      var _this6 = this;

      if (!this.loaded) {
        this.callQueue.push(function () {
          _this6.getMenus(callback);
        });
        return;
      }

      var request = {
        gigatoken: this.getToken()
      };
      this.service.getMenus(request).execute(function (resp) {
        _this6.logError('getMenus', resp.err);
        callback(resp.menus, resp.err);
      });
    }

    /*
     * Item
     */

  }, {
    key: 'getItem',
    value: function getItem(id, callback) {
      var _this7 = this;

      // if api is not loaded, add to callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this7.getItem(id, callback);
        });
        return;
      }

      var request = {
        id: id,
        gigatoken: this.getToken()
      };
      this.service.getItem(request).execute(function (resp) {
        _this7.logError('getItem', resp.err);
        callback(resp.item, resp.err);
      });
    }
  }, {
    key: 'saveItem',
    value: function saveItem(item, callback) {
      var _this8 = this;

      // if api is not loaded, add to callQueue
      if (!this.loaded) {
        this.callQueue.push(function () {
          _this8.saveItem(item, callback);
        });
        return;
      }
      item.min_servings = Number(item.min_servings);
      item.max_servings = Number(item.max_servings);
      item.cook_price_per_serving = Number(item.cook_price_per_serving);
      var request = {
        gigatoken: this.getToken(),
        item: item
      };

      this.service.saveItem(request).execute(function (resp) {
        _this8.logError('saveItem', resp.err);
        callback(resp.item, resp.err);
      });
    }
  }, {
    key: 'refreshToken',
    value: function refreshToken() {
      var _this9 = this;

      // if api is not loaded, add to _callQueue
      if (!this.loaded) {
        this.callQueue.push(this.refreshToken);
        return;
      }
      var request = {
        gigatoken: this.getToken()
      };
      this.service.refreshToken(request).execute(function (resp) {
        if (!_this9.logError('refreshToken', resp.err)) {
          COOK.User.update(resp.gigatoken);
        }
      });
    }

    /*
     * Utils
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
  var apiRoot = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/api';
  if (COOK.isDev) {
    apiRoot = 'http://localhost:8080/_ah/api';
  } else if (COOK.isStage) {
    apiRoot = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/api';
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