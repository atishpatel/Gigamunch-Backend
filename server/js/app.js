function GetToken() {
    const name = 'AUTHTKN=';
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length);
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
function SetToken(cvalue) {
    const jwt = GetJWT(cvalue);
    const d = new Date(0);
    if (jwt) {
        d.setUTCSeconds(jwt.exp);
    }
    document.cookie = `AUTHTKN=${cvalue}; expires=${d.toUTCString()}; path=/`;
    if (location.hostname === 'localhost') {
        window.localStorage.setItem('AUTHTKN', cvalue);
    }
}
function GetJWT(tkn) {
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

addEventListener('UserUpdated', UpdateUser);
function IsLoggedIn() {
    const tkn = GetToken();
    if (tkn === '') {
        return false;
    }
    return true;
}
let ID = '';
let Email = '';
let FirstName = '';
let LastName = '';
let PhotoURL = '';
function UpdateUser() {
    const tkn = GetToken();
    if (!tkn) {
        return;
    }
    const jwt = GetJWT(tkn);
    if (!jwt) {
        return;
    }
    ID = jwt.id;
    Email = jwt.email;
    FirstName = jwt.first_name;
    LastName = jwt.last_name;
    PhotoURL = jwt.photo_url;
}
function IsAdmin() {
    const jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 2);
}
function HasCreditCard() {
    const jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 0);
}
function getKthBit(x, k) {
    return (((x >> k) & 1) === 1);
}
UpdateUser();


var user = Object.freeze({
	IsLoggedIn: IsLoggedIn,
	get ID () { return ID; },
	get Email () { return Email; },
	get FirstName () { return FirstName; },
	get LastName () { return LastName; },
	get PhotoURL () { return PhotoURL; },
	UpdateUser: UpdateUser,
	IsAdmin: IsAdmin,
	HasCreditCard: HasCreditCard
});

function Login(token) {
    const url = '/api/v1/Login';
    const req = {
        token,
    };
    return callFetch(url, 'POST', req).then((resp) => {
        if (resp && resp.token) {
            SetToken(resp.token);
        }
        return resp;
    });
}
function callFetch(url, method, body) {
    return fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    }).then((resp) => {
        return resp.json();
    }).catch((err) => {
        console.error('failed to callFetch', err);
        console.error('details: ', err.code, err.name, err.message, err.detail);
    });
}


var service = Object.freeze({
	Login: Login
});

APP.Service = service;
APP.User = user;
console.log('app.js loaded');
