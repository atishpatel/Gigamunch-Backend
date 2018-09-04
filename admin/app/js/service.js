import { GetToken, } from './utils/token';
let baseURL = '/admin/api/v1/';
if (location.hostname === 'localhost') {
    baseURL = 'https://gigamunch-omninexus-dev.appspot.com/admin/api/v1/';
}
export function GetSubscriber(email) {
    const url = baseURL + 'GetSubscriber';
    const req = {
        email,
    };
    return callFetch(url, 'GET', req);
}
export function GetHasSubscribed(date) {
    const url = baseURL + 'GetHasSubscribed';
    const req = {
        date: date.toISOString(),
    };
    return callFetch(url, 'GET', req);
}
export function GetUnpaidSublogs(limit) {
    const url = baseURL + 'GetUnpaidSublogs';
    const req = {
        limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetSubscriberSublogs(email) {
    const url = baseURL + 'GetSubscriberSublogs';
    const req = {
        email,
    };
    return callFetch(url, 'GET', req);
}
export function ProcessSublog(date, email) {
    const url = baseURL + 'ProcessSublog';
    const req = {
        date,
        email,
    };
    return callFetch(url, 'POST', req);
}
export function GetExecutions(start, limit) {
    const url = baseURL + 'GetExecutions';
    const req = {
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetExecution(id) {
    const url = baseURL + 'GetExecution';
    const req = {
        id,
    };
    return callFetch(url, 'GET', req);
}
export function UpdateExecution(execution) {
    const url = baseURL + 'UpdateExecution';
    const req = {
        execution,
    };
    return callFetch(url, 'POST', req);
}
export function GetActivityForDate() {
}
export function GetLogs(start, limit) {
    const url = baseURL + 'GetLogs';
    const req = {
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}
export function GetLog(id) {
    const url = baseURL + 'GetLog';
    const req = {
        id,
    };
    return callFetch(url, 'GET', req);
}
export function GetLogsByEmail(start, limit, email) {
    const url = baseURL + 'GetLogsByEmail';
    const req = {
        email,
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}
function callFetch(url, method, body) {
    const config = {
        method,
        headers: {
            'Content-Type': 'application/json',
            'auth-token': GetToken(),
            'Access-Control-Allow-Origin': '*',
        },
    };
    let URL = url;
    if (method === 'GET') {
        URL += '?' + serializeParams(body);
    }
    else {
        config.body = JSON.stringify(body);
    }
    return fetch(URL, config)
        .then((resp) => {
        return resp.json();
    })
        .catch((err) => {
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
