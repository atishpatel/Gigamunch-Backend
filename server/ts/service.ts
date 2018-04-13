import {
  GetToken,
  SetToken,
} from './token';

const baseURL = '/api/v1/';

// TODO: setup fetch loading

// Auth
export function Login(token: string): Promise<any> {
  const url: string = '/api/v1/Login';
  const req: TokenOnlyReq = {
    token,
  };
  return callFetch(url, 'POST', req).then((resp) => {
    if (resp && resp.token) {
      SetToken(resp.token);
    }
    return resp;
  });
}

function callFetch(url: string, method: string, body: object): Promise<Response> {
  return fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json',
      // 'auth-token': GetToken(),
    },
    body: JSON.stringify(body),
  }).then((resp: Response) => {
    return resp.json();
  }).catch((err: Error) => {
    console.error('failed to callFetch', err);
    console.error('details: ', err.code, err.name, err.message, err.detail);
  });
}
