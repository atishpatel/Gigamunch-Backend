import { GetToken, IsAdmin } from './auth';
import { IsDev } from './env';

let baseURL = '/admin/api';
if (IsDev()) {
  baseURL = 'https://gigamunch-omninexus-dev.appspot.com/admin/api';
}

// Subscribers
export function GetHasSubscribed(start: number, limit: number): Promise<AdminAPI.GetHasSubscribedRespV2> {
  const url: string = baseURL + '/v2/GetHasSubscribed';
  const req: AdminAPI.GetHasSubscribedReq = {
    start,
    limit,
  };
  return callFetch(url, 'GET', req);
}

export function GetExecution(idOrDate: string): Promise<AdminAPI.GetExecutionResp> {
  const url: string = baseURL + '/v1/GetExecution';
  const req: AdminAPI.GetExecutionReq = {
    idOrDate,
  };
  return callFetch(url, 'GET', req);
}

function callFetch(url: string, method: string, body: object): Promise<any> {
  return GetToken().then((token) => {
    return callFetchWithToken(url, method, body, token);
  });
}

function callFetchWithToken(url: string, method: string, body: object, token: string) {
  const config: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'auth-token': token,
      'Access-Control-Allow-Origin': '*',
    },
  };
  let URL = url;
  if (method === 'GET') {
    URL += '?' + serializeParams(body);
  } else {
    config.body = JSON.stringify(body);
  }
  return fetch(URL, config)
    .then((resp: Response) => {
      return resp.json();
    })
    .catch((err: any) => {
      console.error('failed to callFetch', err);
    });
}

function serializeParams(obj: any): string {
  const str = [];
  let p: any;
  p = 0;
  for (p in obj) {
    if (obj.hasOwnProperty(p)) {
      const k: any = p;
      const v: any = obj[p];
      str.push(v !== null && typeof v === 'object' ? serializeParams(v) : encodeURIComponent(k) + '=' + encodeURIComponent(v));
    }
  }
  return str.join('&');
}

interface APIResponse {
  token: string;
  json(): APIResponse;
}
