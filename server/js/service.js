import { SetToken, } from './token';
const baseURL = '/api/v1/';
export function Login(token) {
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
