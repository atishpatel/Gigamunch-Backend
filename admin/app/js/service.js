var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
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
var _this = this;
import { GetToken } from './auth';
var baseURL = '/admin/api/v1/';
if (location.hostname === 'localhost') {
    baseURL = 'https://gigamunch-omninexus-dev.appspot.com/admin/api/v1/';
}
export function GetSubscriber(email) {
    var url = baseURL + 'GetSubscriber';
    var req = {
        email: email,
    };
    return callFetch(url, 'GET', req);
}
export function GetHasSubscribed(date) {
    var url = baseURL + 'GetHasSubscribed';
    var req = {
        date: date.toISOString(),
    };
    return callFetch(url, 'GET', req);
}
export function GetUnpaidSublogs(limit) {
    var url = baseURL + 'GetUnpaidSublogs';
    var req = {
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetSubscriberSublogs(email) {
    var url = baseURL + 'GetSubscriberSublogs';
    var req = {
        email: email,
    };
    return callFetch(url, 'GET', req);
}
export function ProcessSublog(date, email) {
    var url = baseURL + 'ProcessSublog';
    var req = {
        date: date,
        email: email,
    };
    return callFetch(url, 'POST', req);
}
export function GetExecutions(start, limit) {
    var url = baseURL + 'GetExecutions';
    var req = {
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetExecution(id) {
    var url = baseURL + 'GetExecution';
    var req = {
        id: id,
    };
    return callFetch(url, 'GET', req);
}
export function UpdateExecution(execution) {
    var url = baseURL + 'UpdateExecution';
    var req = {
        execution: execution,
    };
    return callFetch(url, 'POST', req);
}
export function GetActivityForDate() {
}
export function GetLogs(start, limit) {
    var url = baseURL + 'GetLogs';
    var req = {
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetLog(id) {
    var url = baseURL + 'GetLog';
    var req = {
        id: id,
    };
    return callFetch(url, 'GET', req);
}
export function GetLogsByEmail(start, limit, email) {
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
