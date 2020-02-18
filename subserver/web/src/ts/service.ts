import { GetToken } from './auth';
import { IsDev } from './env';

let baseURL = '/sub/api/v1/';
if (IsDev()) {
  baseURL = 'https://gigamunch-omninexus-dev.appspot.com/sub/api/v1/';
}

// GetUserSummary
export function GetUserSummary(): Promise<SubAPI.GetUserSummaryResp> {
  const url: string = baseURL + 'GetUserSummary';
  const req: SubAPI.GetUserSummaryReq = {};
  return GetToken().then((token) => {
    if (!token) {
      return Promise.resolve({ error: {} });
    }
    return callFetchWithToken(url, 'GET', req, token);
  });
}

// GetAccountInfo
export function GetAccountInfo(): Promise<SubAPI.GetAccountInfoResp> {
  const url: string = baseURL + 'GetAccountInfo';
  const req: SubAPI.GetAccountInfoReq = {};
  return callFetch(url, 'GET', req);
}

// Execution
export function GetExecutions(start: number, limit: number): Promise<SubAPI.GetExecutionsResp> {
  const url: string = baseURL + 'GetExecutions';
  const req: SubAPI.GetExecutionsReq = {
    start,
    limit,
  };
  return callFetch(url, 'GET', req);
}

export function GetExecutionsAfterDate(date: Date | string): Promise<SubAPI.GetExecutionsResp> {
  const url: string = baseURL + 'GetExecutionsAfterDate';
  let dateString = '';
  if (typeof (date) === 'string') {
    dateString = date;
  } else if (typeof (date) === 'object') {
    dateString = date.toISOString();
  }
  const req: SubAPI.GetExecutionsDateReq = {
    date: dateString,
  };
  return callFetch(url, 'GET', req);
}

export function GetExecutionsBeforeDate(date: Date | string): Promise<SubAPI.GetExecutionsResp> {
  const url: string = baseURL + 'GetExecutionsBeforeDate';
  let dateString = '';
  if (typeof (date) === 'string') {
    dateString = date;
  } else if (typeof (date) === 'object') {
    dateString = date.toISOString();
  }
  const req: SubAPI.GetExecutionsDateReq = {
    date: dateString,
  };
  return callFetch(url, 'GET', req);
}

export function GetExecution(idOrDate: string): Promise<SubAPI.GetExecutionResp> {
  const url: string = baseURL + 'GetExecution';
  const req: SubAPI.GetExecutionReq = {
    idOrDate,
  };
  return callFetch(url, 'GET', req);
}

export function SkipActivity(date: string): Promise<SubAPI.ErrorOnlyResp> {
  const url: string = baseURL + 'SkipActivity';
  const req: SubAPI.DateReq = {
    date,
  };
  return callFetch(url, 'POST', req);
}

export function UnskipActivity(date: string): Promise<SubAPI.ErrorOnlyResp> {
  const url: string = baseURL + 'UnskipActivity';
  const req: SubAPI.DateReq = {
    date,
  };
  return callFetch(url, 'POST', req);
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
