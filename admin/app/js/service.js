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
export function GetExecution(idOrDate) {
    var url = baseURL + 'GetExecution';
    var req = {
        idOrDate: idOrDate,
    };
    return callFetch(url, 'GET', req);
}
export function UpdateExecution(mode, execution) {
    var url = baseURL + 'UpdateExecution';
    var req = {
        mode: mode,
        execution: execution,
    };
    return callFetch(url, 'POST', req);
}
export function GetActivityForDate() { }
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
export function GetLogsByEmail(start, limit, id) {
    var url = baseURL + 'GetLogsByEmail';
    var req = {
        id: id,
        start: start,
        limit: limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetLogsByExecution(execution_id) {
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
