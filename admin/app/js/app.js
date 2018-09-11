function GetToken() {
    var name = 'AUTHTKN=';
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length).replace(/\n/g, '');
        }
    }
    if (location.hostname === 'localhost') {
        var tnk = window.localStorage.getItem('AUTHTKN');
        if (!tnk) {
            return '';
        }
        return tnk;
    }
    return '';
}
function SetToken(cvalue) {
    var jwt = GetJWT(cvalue);
    var d = new Date(0);
    if (jwt) {
        d.setUTCSeconds(jwt.exp);
    }
    document.cookie = "AUTHTKN=" + cvalue + "; expires=" + d.toUTCString() + "; path=/";
    if (location.hostname === 'localhost') {
        window.localStorage.setItem('AUTHTKN', cvalue);
    }
}
function GetJWT(tkn) {
    if (!tkn) {
        return null;
    }
    var tknConv = tkn.replace(/[+\/]/g, function (m0) {
        return m0 === '+' ? '-' : '_';
    }).replace(/=/g, '');
    var userString = tknConv.split('.')[1].replace(/\s/g, '');
    return JSON.parse(window.atob(userString.replace(/[-_]/g, function (m0) {
        return m0 === '-' ? '+' : '/';
    }).replace(/[^A-Za-z0-9\+\/]/g, '')));
}

var TokenUtil = /*#__PURE__*/Object.freeze({
    GetToken: GetToken,
    SetToken: SetToken,
    GetJWT: GetJWT
});

var baseURL = '/admin/api/v1/';
if (location.hostname === 'localhost') {
    baseURL = 'https://gigamunch-omninexus-dev.appspot.com/admin/api/v1/';
}
function GetSubscriber(email) {
    var url = baseURL + 'GetSubscriber';
    var req = {
        email: email,
    };
    return callFetch(url, 'GET', req);
}
function GetHasSubscribed(date) {
    var url = baseURL + 'GetHasSubscribed';
    var req = {
        date: date.toISOString(),
    };
    return callFetch(url, 'GET', req);
}
function GetUnpaidSublogs(limit) {
    var url = baseURL + 'GetUnpaidSublogs';
    var req = {
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
function GetSubscriberSublogs(email) {
    var url = baseURL + 'GetSubscriberSublogs';
    var req = {
        email: email,
    };
    return callFetch(url, 'GET', req);
}
function ProcessSublog(date, email) {
    var url = baseURL + 'ProcessSublog';
    var req = {
        date: date,
        email: email,
    };
    return callFetch(url, 'POST', req);
}
function GetExecutions(start, limit) {
    var url = baseURL + 'GetExecutions';
    var req = {
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
function GetExecution(id) {
    var url = baseURL + 'GetExecution';
    var req = {
        id: id,
    };
    return callFetch(url, 'GET', req);
}
function UpdateExecution(execution) {
    var url = baseURL + 'UpdateExecution';
    var req = {
        execution: execution,
    };
    return callFetch(url, 'POST', req);
}
function GetActivityForDate() {
}
function GetLogs(start, limit) {
    var url = baseURL + 'GetLogs';
    var req = {
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
function GetLog(id) {
    var url = baseURL + 'GetLog';
    var req = {
        id: id,
    };
    return callFetch(url, 'GET', req);
}
function GetLogsByEmail(start, limit, email) {
    var url = baseURL + 'GetLogsByEmail';
    var req = {
        email: email,
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
function callFetch(url, method, body) {
    var config = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'auth-token': GetToken(),
            'Access-Control-Allow-Origin': '*',
        },
    };
    var URL = url;
    if (method === 'GET') {
        URL += '?' + serializeParams(body);
    }
    else {
        config.body = JSON.stringify(body);
    }
    return fetch(URL, config)
        .then(function (resp) {
        return resp.json();
    })
        .catch(function (err) {
        console.error('failed to callFetch', err);
    });
}
function serializeParams(obj) {
    var str = [];
    var p;
    p = 0;
    for (p in obj) {
        if (obj.hasOwnProperty(p)) {
            var k = p;
            var v = obj[p];
            str.push((v !== null && typeof v === 'object') ?
                serializeParams(v) :
                encodeURIComponent(k) + '=' + encodeURIComponent(v));
        }
    }
    return str.join('&');
}

var Service = /*#__PURE__*/Object.freeze({
    GetSubscriber: GetSubscriber,
    GetHasSubscribed: GetHasSubscribed,
    GetUnpaidSublogs: GetUnpaidSublogs,
    GetSubscriberSublogs: GetSubscriberSublogs,
    ProcessSublog: ProcessSublog,
    GetExecutions: GetExecutions,
    GetExecution: GetExecution,
    UpdateExecution: UpdateExecution,
    GetActivityForDate: GetActivityForDate,
    GetLogs: GetLogs,
    GetLog: GetLog,
    GetLogsByEmail: GetLogsByEmail
});

var UserUpdated = 'UserUpdated';
function Fire(eventName, detail) {
    if (detail === void 0) { detail = {}; }
    var event = new CustomEvent(eventName, {
        detail: detail,
        bubbles: true,
        composed: true,
    });
    window.dispatchEvent(event);
}
function FireToast(t, detail) {
    var event = new CustomEvent('toast', {
        detail: detail,
        bubbles: true,
        composed: true,
    });
    t.dispatchEvent(event);
}
function FireError() {
}

var EventUtil = /*#__PURE__*/Object.freeze({
    UserUpdated: UserUpdated,
    Fire: Fire,
    FireToast: FireToast,
    FireError: FireError
});

addEventListener(UserUpdated, UpdateUser);
function IsLoggedIn() {
    var tkn = GetToken();
    if (tkn === '') {
        return false;
    }
    return true;
}
var ID = '';
var Email = '';
var FirstName = '';
var LastName = '';
var PhotoURL = '';
var Token = '';
function UpdateUser() {
    var tkn = GetToken();
    if (!tkn) {
        return;
    }
    var jwt = GetJWT(tkn);
    if (!jwt) {
        return;
    }
    ID = jwt.id;
    Email = jwt.email;
    FirstName = jwt.first_name;
    LastName = jwt.last_name;
    PhotoURL = jwt.photo_url;
    Token = tkn;
}
function IsAdmin() {
    var jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 2);
}
function HasCreditCard() {
    var jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 0);
}
function getKthBit(x, k) {
    return (((x >> k) & 1) === 1);
}
UpdateUser();

var User = /*#__PURE__*/Object.freeze({
    IsLoggedIn: IsLoggedIn,
    get ID () { return ID; },
    get Email () { return Email; },
    get FirstName () { return FirstName; },
    get LastName () { return LastName; },
    get PhotoURL () { return PhotoURL; },
    get Token () { return Token; },
    UpdateUser: UpdateUser,
    IsAdmin: IsAdmin,
    HasCreditCard: HasCreditCard
});

APP.Service = Service;
APP.User = User;
APP.Event = EventUtil;
APP.Token = TokenUtil;
