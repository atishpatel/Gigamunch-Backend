var Events = {
    UserUpdated: 'user-updated',
    SignedOut: 'signed-out',
    SignedIn: 'signed-in',
};
var userLoaded = false;
function fireUserUpdated() {
    var event = document.createEvent('Event');
    event.initEvent(Events.UserUpdated, true, true);
    window.dispatchEvent(event);
}
function setUser(user) {
    APP.User = user;
    if (user) {
        user.getIdTokenResult(false)
            .then(function (tokenResult) {
            var adminClaim = tokenResult.claims['admin'];
            if (adminClaim) {
                user.Admin = true;
            }
            else {
                user.Admin = false;
            }
            user.IsAdmin = function () {
                return user.Admin;
            };
            APP.User = user;
            console.log('user', user);
            fireUserUpdated();
            userLoaded = true;
        });
        return;
    }
    APP.User = user;
    console.log('user', user);
    fireUserUpdated();
    userLoaded = true;
}
function IsAdmin() {
    return GetUser().then(function (user) {
        if (!user) {
            return false;
        }
        return user.getIdTokenResult(false).then(function (tokenResult) {
            var adminClaim = tokenResult.claims['admin'];
            if (adminClaim) {
                return true;
            }
            return false;
        });
    });
}
function GetUser() {
    return new Promise(function (resolve, reject) {
        if (userLoaded) {
            resolve(APP.User);
        }
        var unsubscribe = firebase.auth().onAuthStateChanged(function (user) {
            if (!APP.User) {
                setUser(user);
            }
            unsubscribe();
            resolve(APP.User);
        }, reject);
    });
}
function GetToken() {
    return GetUser().then(function (user) {
        if (!user) {
            return '';
        }
        return user.getIdToken(false);
    });
}
function SignOut() {
    firebase.auth().signOut();
}
function SetupFirebase() {
    var config;
    if (APP.IsProd) {
        config = {
            apiKey: 'AIzaSyC-1vqT4YIKXVmrGkaoVSj1BJnm48NxlT0',
            authDomain: 'gigamunch-omninexus.firebaseapp.com',
            databaseURL: 'https://gigamunch-omninexus.firebaseio.com',
            projectId: 'gigamunch-omninexus',
            storageBucket: 'gigamunch-omninexus.appspot.com',
            messagingSenderId: '837147123677',
        };
    }
    else {
        config = {
            apiKey: 'AIzaSyBHPe4B4k72ljnBXszda6AJGGg21YqEJ4g',
            authDomain: 'gigamunch-omninexus-dev.firebaseapp.com',
            databaseURL: 'https://gigamunch-omninexus-dev.firebaseio.com',
            projectId: 'gigamunch-omninexus-dev',
            storageBucket: 'gigamunch-omninexus-dev.appspot.com',
            messagingSenderId: '108585202286',
        };
    }
    firebase.initializeApp(config);
}
SetupFirebase();
firebase.auth().onAuthStateChanged(function (user) {
    var eventName;
    if (!user) {
        eventName = Events.SignedOut;
    }
    else {
        eventName = Events.SignedIn;
        setUser(user);
    }
    var event = document.createEvent('Event');
    event.initEvent(eventName, true, true);
    window.dispatchEvent(event);
});

var Auth = /*#__PURE__*/Object.freeze({
    Events: Events,
    IsAdmin: IsAdmin,
    GetUser: GetUser,
    GetToken: GetToken,
    SignOut: SignOut,
    SetupFirebase: SetupFirebase
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
function GetExecution(idOrDate) {
    var url = baseURL + 'GetExecution';
    var req = {
        idOrDate: idOrDate,
    };
    return callFetch(url, 'GET', req);
}
function UpdateExecution(mode, execution) {
    var url = baseURL + 'UpdateExecution';
    var req = {
        mode: mode,
        execution: execution,
    };
    return callFetch(url, 'POST', req);
}
function GetActivityForDate() { }
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
function GetLogsByEmail(start, limit, id) {
    var url = baseURL + 'GetLogsByEmail';
    var req = {
        id: id,
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
function GetLogsByExecution(execution_id) {
    var url = baseURL + 'GetLogsByExecution';
    var req = {
        execution_id: execution_id,
    };
    return callFetch(url, 'GET', req);
}
function callFetch(url, method, body) {
    return GetToken().then(function (token) {
        var config = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'auth-token': token,
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
            try {
                return resp.json();
            }
            catch (err) {
                return {
                    error: {
                        code: resp.status,
                        message: 'Unknown server error',
                    }
                };
            }
        })
            .catch(function (err) {
            console.error('failed to callFetch', err);
        });
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
            str.push(v !== null && typeof v === 'object' ? serializeParams(v) : encodeURIComponent(k) + '=' + encodeURIComponent(v));
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
    GetLogsByEmail: GetLogsByEmail,
    GetLogsByExecution: GetLogsByExecution
});

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
    Fire: Fire,
    FireToast: FireToast,
    FireError: FireError
});

APP.Auth = Auth;
APP.Service = Service;
APP.Event = EventUtil;
console.log('app.js loaded');
function GetURLParmas() {
    var vars = {};
    window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function (m, key, value) {
        vars[key] = value;
        return value;
    });
    return vars;
}
GetURLParmas();
