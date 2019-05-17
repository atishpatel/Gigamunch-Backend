import { GetToken, IsAdmin } from './auth';
import { IsDev } from './env';

let baseURL = '/admin/api';
if (IsDev()) {
  baseURL = 'https://gigamunch-omninexus-dev.appspot.com/admin/api';
}

// Log

export function GetLogsForUser(start: number, limit: number, id: string): Promise<AdminAPI.GetLogsResp> {
  const url: string = baseURL + '/v1/GetLogsForUser';
  const req: AdminAPI.GetLogsForUserReq = {
    id,
    start,
    limit,
  };
  return callFetch(url, 'GET', req);
}

// Discount
export function GetSubscriberDiscounts(id: string): Promise<AdminAPI.GetSubscriberDiscountsResp> {
  const url: string = baseURL + '/v1/GetSubscriberDiscounts';
  const req: AdminAPI.UserIDReq = {
    ID: id,
  };
  return callFetch(url, 'GET', req);
}

export function DiscountSubscriber(user_id: string, discount_amount: number, discount_percent: number): Promise<AdminAPI.ErrorOnlyResp> {
  const url: string = baseURL + '/v1/DiscountSubscriber';
  const req: AdminAPI.DiscountSubscriberReq = {
    user_id,
    discount_amount,
    discount_percent,
  };
  return callFetch(url, 'GET', req);
}

// Activitiy
export function GetSubscriberActivities(id: string): Promise<AdminAPI.GetSubscriberActivitiesResp> {
  const url: string = baseURL + '/v1/GetSubscriberActivities';
  const req: AdminAPI.UserIDReq = {
    ID: id,
  };
  return callFetch(url, 'GET', req);
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

export function GetSubscriber(id: string): Promise<AdminAPI.GetSubscriberRespV2> {
  const url: string = baseURL + '/v2/GetSubscriber';
  const req: AdminAPI.UserIDReq = {
    ID: id,
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

