'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

/* exported user */
/* global COOK */
window.COOK = window.COOK || {};

var User = function () {
  function User() {
    _classCallCheck(this, User);

    this.isLoggedIn = false;
    // get cookie
    var token = this.getTokenCookie();
    this.update(token, 0);
    if (this.token !== undefined && this.token !== '') {
      this.isLoggedIn = true;
    }
  }

  _createClass(User, [{
    key: 'update',
    value: function update(token) {
      var setCookie = arguments.length <= 1 || arguments[1] === undefined ? 1 : arguments[1];

      if (token === '') {
        return;
      }
      this.token = token;
      var userString = token.split('.')[1];
      var jwt = JSON.parse(window.atob(userString));
      for (var k in jwt) {
        if (k !== '__proto__') {
          this[k] = jwt[k];
        }
      }
      // set permissions
      this.isCook = this.getKthBit(jwt.perm, 0);
      this.isVerifiedCook = this.getKthBit(jwt.perm, 1);
      this.hasAddress = this.getKthBit(jwt.perm, 4);
      this.hasSubMerchantID = this.getKthBit(jwt.perm, 5);
      this.isOnboard = this.getKthBit(jwt.perm, 6);
      // update coookie to new token
      if (setCookie) {
        this.setTokenCookie(token, jwt.exp);
      }
      document.dispatchEvent(new Event('userUpdated', {
        bubbles: true
      }));
    }
  }, {
    key: 'getTokenCookie',
    value: function getTokenCookie() {
      var name = 'GIGATKN=';
      var ca = document.cookie.split(';');
      for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === ' ') {
          c = c.substring(1);
        }if (c.indexOf(name) === 0) return c.substring(name.length, c.length);
      }
      return '';
    }
  }, {
    key: 'setTokenCookie',
    value: function setTokenCookie(cvalue, exptime) {
      var d = new Date(0);
      d.setUTCSeconds(exptime);
      document.cookie = 'GIGATKN=' + cvalue + '; expires=' + d.toUTCString() + '; path=/';
    }
  }, {
    key: 'getKthBit',
    value: function getKthBit(x, k) {
      return (x >> k & 1) === 1;
    }
  }]);

  return User;
}();

window.COOK.User = new User();

// redirect if token is empty
if (!window.COOK.User.isLoggedIn) {
  window.location = '/becomechef?mode=select';
}