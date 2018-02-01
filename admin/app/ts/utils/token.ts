export function GetToken(): string {
  const name = 'AUTHTKN=';
  const ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) === 0) {
      return c.substring(name.length, c.length).replace(/\n/g,'');
    }
  }
  if (location.hostname === 'localhost') {
    const tnk = window.localStorage.getItem('AUTHTKN');
    if (!tnk) {
      return '';
    }
    return tnk;
  }
  return '';
}

export function SetToken(cvalue: string) {
  const jwt = GetJWT(cvalue);
  const d = new Date(0);
  d.setUTCSeconds(jwt.exp);
  document.cookie = `AUTHTKN=${cvalue}; expires=${d.toUTCString()}; path=/`;
  if (location.hostname === 'localhost') {
    window.localStorage.setItem('AUTHTKN', cvalue);
  }
}

interface JWT {
  id: string;
  email: string;
  last_name: string;
  first_name: string;
  photo_url: string;
  perm: number;
  exp: number;
}

export function GetJWT(tkn: string): JWT | null {
  if (!tkn) {
    return null;
  }
  const tknConv = tkn.replace(/[+\/]/g, (m0) => {
    return m0 === '+' ? '-' : '_';
  }).replace(/=/g, '');
  const userString = tknConv.split('.')[1].replace(/\s/g, '');
  return JSON.parse(window.atob(userString.replace(/[-_]/g, (m0) => {
    return m0 === '-' ? '+' : '/';
  }).replace(/[^A-Za-z0-9\+\/]/g, '')));
}
