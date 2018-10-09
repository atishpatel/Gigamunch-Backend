import { GetToken } from './auth';

let baseURL = '/sub/api/v1/';
if (location.hostname === 'localhost') {
  baseURL = 'https://gigamunch-omninexus-dev.appspot.com/sub/api/v1/';
}

// Execution
export function GetExecutions(start: number, limit: number): Promise<GetExecutionsResp> {
  const url: string = baseURL + 'GetExecutions';
  const req: GetExecutionsReq = {
    start,
    limit,
  };
  return callFetch(url, 'GET', req);
}

export function GetExecution(id: number): Promise<GetExecutionResp> {
  const url: string = baseURL + 'GetExecution';
  const req: GetExecutionReq = {
    id,
  };
  return callFetch(url, 'GET', req);
}

function callFetch(url: string, method: string, body: object): Promise<any> {
  return GetToken().then((token) => {
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
