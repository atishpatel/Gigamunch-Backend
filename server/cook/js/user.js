"use strict";
COOK = COOK || {};
var User = (function () {
    function User() {
        this.isLoggedIn = false;
        // get cookie
        var token = this.getTokenCookie();
        this.update(token, 0);
        if (this.token && this.token !== '') {
            this.isLoggedIn = true;
        }
    }
    User.prototype.update = function (tkn, setCookie) {
        if (setCookie === void 0) { setCookie = 1; }
        if (tkn === '') {
            return;
        }
        this.token = tkn.replace(/[+\/]/g, function (m0) {
            return m0 === '+' ? '-' : '_';
        }).replace(/=/g, '');
        var userString = tkn.split('.')[1].replace(/\s/g, '');
        var jwt = JSON.parse(window.atob(userString.replace(/[-_]/g, function (m0) {
            return m0 === '-' ? '+' : '/';
        }).replace(/[^A-Za-z0-9\+\/]/g, '')));
        this.id = jwt.id;
        this.email = jwt.email;
        this.name = jwt.name;
        this.perm = jwt.perm;
        this.photo_url = jwt.photo_url;
        // set permissions
        this.isCook = this.getKthBit(jwt.perm, 0);
        this.isVerifiedCook = this.getKthBit(jwt.perm, 1);
        this.isAdmin = this.getKthBit(jwt.perm, 2);
        this.hasAddress = this.getKthBit(jwt.perm, 4);
        this.hasSubMerchantID = this.getKthBit(jwt.perm, 5);
        this.isOnboard = this.getKthBit(jwt.perm, 6);
        // update coookie to new token
        if (setCookie) {
            this.setTokenCookie(tkn, jwt.exp);
        }
        document.dispatchEvent(new Event('userUpdated', {
            bubbles: true,
        }));
    };
    User.prototype.getTokenCookie = function () {
        var name = 'GIGATKN=';
        var ca = document.cookie.split(';');
        for (var i = 0; i < ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0) === ' ') {
                c = c.substring(1);
            }
            if (c.indexOf(name) === 0) {
                return c.substring(name.length, c.length);
            }
        }
        return '';
    };
    User.prototype.setTokenCookie = function (cvalue, exptime) {
        var d = new Date(0);
        d.setUTCSeconds(exptime);
        document.cookie = "GIGATKN=" + cvalue + "; expires=" + d.toUTCString() + "; path=/";
    };
    User.prototype.getKthBit = function (x, k) {
        return (((x >> k) & 1) === 1);
    };
    return User;
}());
COOK.User = new User();
// redirect if token is empty
if (!COOK.User.isLoggedIn) {
    window.location.href = '/becomechef?gstate=login';
}
