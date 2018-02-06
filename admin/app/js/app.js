const UserUpdated = 'UserUpdated';
function Fire(eventName, detail = {}) {
    const event = new CustomEvent(eventName, {
        detail,
        bubbles: true,
        composed: true,
    });
    window.dispatchEvent(event);
}
function FireToast(t, detail) {
    const event = new CustomEvent('toast', {
        detail,
        bubbles: true,
        composed: true,
    });
    t.dispatchEvent(event);
}
function FireError() {
}


var EventUtil = Object.freeze({
	UserUpdated: UserUpdated,
	Fire: Fire,
	FireToast: FireToast,
	FireError: FireError
});

function GetToken() {
    const name = 'AUTHTKN=';
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length).replace(/\n/g, '');
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
    d.setUTCSeconds(jwt.exp);
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

const baseURL = '/admin/api/v1/';
function GetUnpaidSublogs(limit) {
    const url = baseURL + 'GetUnpaidSublogs';
    const req = {
        limit,
    };
    return callFetch(url, 'GET', req);
}
function ProcessSublog(date, email) {
    const url = baseURL + 'ProcessSublog';
    const req = {
        date,
        email,
    };
    return callFetch(url, 'POST', req);
}
function Login(token) {
    const url = baseURL + 'Login';
    const req = {
        token,
    };
    return callFetch(url, 'POST', req).then((resp) => {
        if (resp && resp.token) {
            SetToken(resp.token);
            Fire(UserUpdated);
        }
        return resp;
    });
}
function Refresh(token) {
    const url = baseURL + 'Refresh';
    const req = {
        token,
    };
    return callFetch(url, 'POST', req).then((resp) => {
        if (resp && resp.token) {
            SetToken(resp.token);
            Fire(UserUpdated);
        }
        return resp;
    });
}
function GetActivityForDate() {
}
function GetLogs(start, limit) {
    const url = baseURL + 'GetLogs';
    const req = {
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}
function GetLog(id) {
    const url = baseURL + 'GetLog';
    const req = {
        id,
    };
    return callFetch(url, 'GET', req);
}
function callFetch(url, method, body) {
    const config = {
        method,
        headers: {
            'Content-Type': 'application/json',
            'auth-token': GetToken(),
        },
    };
    let URL = url;
    if (method === 'GET') {
        URL += '?' + serializeParams(body);
    }
    else {
        config.body = JSON.stringify(body);
    }
    return fetch(URL, config).then((resp) => {
        return resp.json();
    }).catch((err) => {
        console.error('failed to callFetch', err);
    });
}
function serializeParams(obj) {
    const str = [];
    let p;
    p = 0;
    for (p in obj) {
        if (obj.hasOwnProperty(p)) {
            const k = p;
            const v = obj[p];
            str.push((v !== null && typeof v === 'object') ?
                serializeParams(v) :
                encodeURIComponent(k) + '=' + encodeURIComponent(v));
        }
    }
    return str.join('&');
}


var Service = Object.freeze({
	GetUnpaidSublogs: GetUnpaidSublogs,
	ProcessSublog: ProcessSublog,
	Login: Login,
	Refresh: Refresh,
	GetActivityForDate: GetActivityForDate,
	GetLogs: GetLogs,
	GetLog: GetLog
});

addEventListener(UserUpdated, UpdateUser);
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


var User = Object.freeze({
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

APP.Service = Service;
APP.User = User;
APP.Event = EventUtil;
