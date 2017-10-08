import { Fire, UserUpdated } from './utils/event';
import { GetToken, SetToken } from './utils/token';

const baseURL = '/admin/api/v1/';

// Auth
export function Login(token: string): Promise < any > {
  const url: string = baseURL + 'Login';
  const req: TokenOnlyReq = {
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

export function Refresh(token: string): Promise < any > {
  const url: string = baseURL + 'Refresh';
  const req: TokenOnlyReq = {
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

// Activity
export function GetActivityForDate() {

}

// Logs
export function GetLogs(start: number, limit: number): Promise < any > {
  const url: string = baseURL + 'GetLogs';
  const req: GetLogsReq = {
    start,
    limit,
  };
  return callFetch(url, 'POST', req);
}

export function GetLog(id: number): Promise < any > {
  const url: string = baseURL + 'GetLog';
  const req: GetLogReq = {
    id,
  };

  return callFetch(url, 'POST', req);
}

function callFetch(url: string, method: string, body: object): Promise < Response > {
  return fetch(url, {
    method,
    headers: {
      'Content-Type': 'application/json',
      'auth-token': GetToken(),
    },
    body: JSON.stringify(body),
  }).then((resp: Response) => {
    return resp.json();
  }).catch((err: any) => {
    console.error('failed to callFetch', err);
  });
}
