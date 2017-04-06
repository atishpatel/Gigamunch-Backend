/* exported user */
/* global COOK */
window.COOK = window.COOK || {};

class User {
  constructor() {
    this.isLoggedIn = false;
    // get cookie
    const token = this.getTokenCookie();
    this.update(token, 0);
    if (this.token && this.token !== '') {
      this.isLoggedIn = true;
    }
  }

  update(tkn, setCookie = 1) {
    if (tkn === '') {
      return;
    }
    this.token = tkn.replace(/[+\/]/g, (m0) => {
      return m0 === '+' ? '-' : '_';
    }).replace(/=/g, '');
    const userString = tkn.split('.')[1].replace(/\s/g, '');
    const jwt = JSON.parse(window.atob(userString.replace(/[-_]/g, (m0) => {
      return m0 === '-' ? '+' : '/';
    }).replace(/[^A-Za-z0-9\+\/]/g, '')));
    for (const k in jwt) {
      if (k !== '__proto__') {
        this[k] = jwt[k];
      }
    }
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
      bubbles: true
    }));
  }

  getTokenCookie() {
    const name = 'GIGATKN=';
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') c = c.substring(1);
      if (c.indexOf(name) === 0) return c.substring(name.length, c.length);
    }
    return '';
  }

  setTokenCookie(cvalue, exptime) {
    const d = new Date(0);
    d.setUTCSeconds(exptime);
    document.cookie = `GIGATKN=${cvalue}; expires=${d.toUTCString()}; path=/`;
  }

  getKthBit(x, k) {
    return (((x >> k) & 1) === 1);
  }
}

window.COOK.User = new User();

// redirect if token is empty
if (!window.COOK.User.isLoggedIn) {
  window.location = '/becomechef?gstate=login';
}
