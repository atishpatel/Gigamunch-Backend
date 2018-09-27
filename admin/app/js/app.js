var Events = {
    UserUpdated: 'user-updated',
    SignedOut: 'signed-out',
    SignedIn: 'signed-in',
};
var userLoaded = false;
function setUser(user) {
    APP.User = user;
    if (user) {
        user.getIdTokenResult(false).then(function (tokenResult) {
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
        });
    }
    var event = document.createEvent('Event');
    event.initEvent(Events.UserUpdated, true, true);
    window.dispatchEvent(event);
    console.log('user', user);
    userLoaded = true;
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
    GetUser: GetUser,
    GetToken: GetToken,
    SignOut: SignOut,
    SetupFirebase: SetupFirebase
});

var __awaiter = (undefined && undefined.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (undefined && undefined.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var _this = undefined;
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
var callFetch = function (url, method, body) { return __awaiter(_this, void 0, void 0, function () {
    var token, config, URL;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4, GetToken()];
            case 1:
                token = _a.sent();
                config = {
                    method: method,
                    headers: {
                        'Content-Type': 'application/json',
                        'auth-token': token,
                        'Access-Control-Allow-Origin': '*',
                    },
                };
                URL = url;
                if (method === 'GET') {
                    URL += '?' + serializeParams(body);
                }
                else {
                    config.body = JSON.stringify(body);
                }
                return [2, fetch(URL, config)
                        .then(function (resp) {
                        return resp.json();
                    })
                        .catch(function (err) {
                        console.error('failed to callFetch', err);
                    })];
        }
    });
}); };
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
