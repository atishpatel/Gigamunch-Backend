'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

/* exported user */

var User = function () {
  function User() {
    _classCallCheck(this, User);

    // get cookie
    var token = this._getTokenCookie();
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
      this.isChef = this._getKthBit(jwt.perm, 0);
      this.isVerifiedChef = this._getKthBit(jwt.perm, 1);
      // update coookie to new token
      if (setCookie) {
        this._setTokenCookie(token, jwt.exp);
      }
    }
  }, {
    key: '_getTokenCookie',
    value: function _getTokenCookie() {
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
    key: '_setTokenCookie',
    value: function _setTokenCookie(cvalue, exptime) {
      var d = new Date();
      d.setTime(exptime);
      var expires = '\'expires=\'' + d.toUTCString();
      document.cookie = 'GIGATKN=' + cvalue + '; ' + expires;
    }
  }, {
    key: '_getKthBit',
    value: function _getKthBit(x, k) {
      return (x >> k & 1) === 1;
    }
  }]);

  return User;
}();

window.CHEF = window.CHEF || {};
window.CHEF.User = new User();