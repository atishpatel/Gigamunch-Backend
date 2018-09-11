import { SetToken, } from './token';
var baseURL = '/api/v1/';
export function Login(token) {
    var url = '/api/v1/Login';
    var req = {
        token: token,
    };
    return callFetch(url, 'POST', req).then(function (resp) {
        if (resp && resp.token) {
            SetToken(resp.token);
        }
        return resp;
    });
}
function callFetch(url, method, body) {
    return fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    }).then(function (resp) {
        return resp.json();
    }).catch(function (err) {
        console.error('failed to callFetch', err);
    });
}
