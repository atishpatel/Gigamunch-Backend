/* exported user */
/* global CHEF */
window.CHEF = window.CHEF || {};

class User {
  constructor() {
    this.isLoggedIn = false;
    // get cookie
    const token = this.getTokenCookie();
    this.update(token, 0);
    if (this.token !== undefined && this.token !== '') {
      this.isLoggedIn = true;
    }
  }

  update(token, setCookie = 1) {
    if (token === '') {
      return;
    }
    this.token = token;
    const userString = token.split('.')[1];
    const jwt = JSON.parse(window.atob(userString));
    for (const k in jwt) {
      if (k !== '__proto__') {
        this[k] = jwt[k];
      }
    }
    // set permissions
    this.isChef = this.getKthBit(jwt.perm, 0);
    this.isVerifiedChef = this.getKthBit(jwt.perm, 1);
    this.hasAddress = this.getKthBit(jwt.perm, 4);
    this.hasSubMerchantID = this.getKthBit(jwt.perm, 5);
    // update coookie to new token
    if (setCookie) {
      this.setTokenCookie(token, jwt.exp);
    }
    document.dispatchEvent(new Event('userUpdated', { bubbles: true }));
  }

  getTokenCookie() {
    const name = 'GIGATKN=';
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ')c = c.substring(1);
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

CHEF.User = new User();

// redirect if token is empty
if (!CHEF.User.isLoggedIn) {
  window.location = '/login?mode=select';
}
