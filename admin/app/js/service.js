import { Fire, UserUpdated, } from './utils/event';
import { GetToken, SetToken, } from './utils/token';
const baseURL = '/admin/api/v1/';
export function Login(token) {
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
export function Refresh(token) {
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
export function GetActivityForDate() {
}
export function GetLogs(start, limit) {
    const url = baseURL + 'GetLogs';
    const req = {
        start,
        limit,
    };
    return callFetch(url, 'POST', req);
}
export function GetLog(id) {
    const url = baseURL + 'GetLog';
    const req = {
        id,
    };
    return callFetch(url, 'POST', req);
}
function callFetch(url, method, body) {
    return fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
            'auth-token': GetToken(),
        },
        body: JSON.stringify(body),
    }).then((resp) => {
        return resp.json();
    }).catch((err) => {
        console.error('failed to callFetch', err);
    });
}
