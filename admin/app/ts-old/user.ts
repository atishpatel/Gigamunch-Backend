declare var COOK: any;

COOK = COOK || {};

class UserOld {
  isLoggedIn: boolean;
  token: string;
  isCook: boolean;
  isVerifiedCook: boolean;
  isAdmin: boolean;
  hasAddress: boolean;
  hasSubMerchantID: boolean;
  isOnboard: boolean;
  id: string;
  email: string;
  name: string;
  perm: number;
  photo_url: string;
  photoURL: string;

  constructor() {
    this.isLoggedIn = false;
    // get cookie
    const token = this.getTokenCookie();
    this.update(token, 0);
    if (this.token && this.token !== '') {
      this.isLoggedIn = true;
    }
  }

  update(tkn: string, setCookie = 1) {
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
  }

  getTokenCookie() {
    let name = 'GIGATKN=';
    let ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') { c = c.substring(1); }
      if (c.indexOf(name) === 0) { return c.substring(name.length, c.length).replace(/\n/g,''); }
    }
    name = 'AUTHTKN=';
    ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') { c = c.substring(1); }
      if (c.indexOf(name) === 0) { return c.substring(name.length, c.length).replace(/\n/g,''); }
    }
    return '';
  }

  setTokenCookie(cvalue: string, exptime: number) {
    const d = new Date(0);
    d.setUTCSeconds(exptime);
    document.cookie = `GIGATKN=${cvalue}; expires=${d.toUTCString()}; path=/`;
  }

  getKthBit(x: number, k: number) {
    return (((x >> k) & 1) === 1);
  }
}

COOK.User = new UserOld();

function getTokenCookie(cookieName: string) {
  const name = cookieName+'=';
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') { c = c.substring(1); }
    if (c.indexOf(name) === 0) { return c.substring(name.length, c.length); }
  }
  return '';
}
 
if (!getTokenCookie('AUTHTKN')) {
  const tmptkn = getTokenCookie('GIGATKN')
  if (tmptkn) {
    const exptime = 1522623345;
    const d = new Date(0);
    d.setUTCSeconds(exptime);
    document.cookie = `AUTHTKN=${tmptkn}; expires=${d.toUTCString()}; path=/`;
  }
}

