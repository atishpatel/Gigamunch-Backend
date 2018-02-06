import {
  Fire,
  UserUpdated,
} from './utils/event';
import {
  GetToken,
  SetToken,
} from './utils/token';

const baseURL = '/admin/api/v1/';

// SubLog
export function GetUnpaidSublogs(limit: number): Promise < any > {
  const url: string = baseURL + 'GetUnpaidSublogs';
  const req: GetUnpaidSublogsReq = {
    limit,
  };
  return callFetch(url, 'GET', req);
}

export function ProcessSublog(date: string, email: string): Promise < any > {
  const url: string = baseURL + 'ProcessSublog';
  const req: ProcessSublogsReq = {
    date,
    email,
  };
  return callFetch(url, 'POST', req);
}

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
  return callFetch(url, 'GET', req);
}

export function GetLog(id: number): Promise < any > {
  const url: string = baseURL + 'GetLog';
  const req: GetLogReq = {
    id,
  };

  return callFetch(url, 'GET', req);
}

function callFetch(url: string, method: string, body: object): Promise < APIResponse > {
  const config: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'auth-token': GetToken(),
    },
  };
  let URL = url;
  if (method === 'GET') {
    URL += '?' + serializeParams(body);
  } else {
    config.body = JSON.stringify(body);
  }
  return fetch(URL, config).then((resp: Response) => {
    return resp.json();
  }).catch((err: any) => {
    console.error('failed to callFetch', err);
  });
}

function serializeParams(obj: any):string {
  const str = [];
  let p: any;
  p = 0;
  for (p in obj) {
    if (obj.hasOwnProperty(p)) {
      const k:any = p;
      const v:any = obj[p];
      str.push((v !== null && typeof v === 'object') ?
      serializeParams(v) :
        encodeURIComponent(k) + '=' + encodeURIComponent(v));
    }
  }
  return str.join('&');
}

interface APIResponse {
  token: string;
  json(): APIResponse;
}
