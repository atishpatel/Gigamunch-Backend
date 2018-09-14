import {
  GetToken,
} from './utils/token';

let baseURLOld = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/spi/Service.';
if (APP.IsDev) {
  baseURLOld = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/spi/Service.';
} else if (APP.IsStage) {
  baseURLOld = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/spi/Service.';
}

function getToken() {
  return COOK.User.token;
}

function logError(fnName: string, err: ErrorWithCode) {
  if (err && (err.code === undefined || err.code !== 0)) {
    const desc = `Function: ${fnName} | Message: ${err.message} | Details: ${err.detail}`;
    console.error(desc);
    ga('send', 'exception', {
      exDescription: desc,
      exFatal: false,
    });
    // if (err.code && err.code === 452 && !COOK.isDev) { // code signout
    //   window.location.href = '/signout';
    // }
    return true;
  }
  return false;
}

export function getSubLogs(callback: (sublogs: SubLogs[], err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'getSubLogs';
  const request = {
    gigatoken: getToken(),
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('getSubLogs', resp.err);
        callback(resp.sublogs, resp.err);
  })
}

export function getSubLogsForDate(date: Date, callback: (sublogs: SubLogs, err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'getSubLogsForDate';
  const request = {
    gigatoken: getToken(),
    date: date.toISOString(),
  };
  callFetch(url, 'POST', request).then((resp)=>{logError('getSubLogsForDate', resp.err);
  callback(resp.sublogs, resp.err);})
}

export function getSubEmails(callback: (subEmails: String[], err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'getSubEmails';
  const request = {
    gigatoken: getToken(),
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('getSubEmails', resp.err);
        callback(resp.sub_emails, resp.err);
  })
}

export function getSubEmailsAndSubs(callback: (subEmails: String[], subs: Object[], err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'getSubEmails';
  const request = {
    gigatoken: getToken(),
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('getSubEmails', resp.err);
        callback(resp.sub_emails, resp.subscribers, resp.err);
  })
}

export function skipSubLog(date: Date, subEmail: string, callback: (err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'skipSubLog';
  const request = {
    gigatoken: getToken(),
    date: date.toISOString(),
    sub_email: subEmail,
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('skipSubLog', resp.err);
        callback(resp.err);
  })
}

export function CancelSub(email: string, callback: (err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'CancelSub';
  const request = {
    gigatoken: getToken(),
    email: email,
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('CancelSub', resp.err);
        callback(resp.err);
  })
}

export function discountSubLog(date: Date, subEmail: string, amount: number, percent: number, overrideDiscount: boolean, callback: (err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'DiscountSubLog';
  const request = {
    gigatoken: getToken(),
    date: date.toISOString(),
    sub_email: subEmail,
    amount: amount,
    percent: percent,
    override_discount: overrideDiscount,
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('DiscountSubLog', resp.err);
        callback(resp.err);
  })
}

export function ChangeServingsForDate(date: Date, subEmail: string, servings: number, callback: (err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'ChangeServingsForDate';
  const request = {
    gigatoken: getToken(),
    date: date.toISOString(),
    sub_email: subEmail,
    servings: servings,
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('ChangeServingForDate', resp.err);
    callback(resp.err);
  })
}

export function ChangeServingsPermanently(email: string, servings: number, vegetarian: boolean, callback: (err: ErrorWithCode) => void) {
  const url: string = baseURLOld + 'ChangeServingsPermanently';
  const request = {
    gigatoken: getToken(),
    email: email,
    servings: servings,
    vegetarian: vegetarian,
  };
  callFetch(url, 'POST', request).then((resp)=>{
    logError('ChangeServingsPermanently', resp.err);
        callback(resp.err);
  })
}

export function GetGeneralStats(callback: (resp: Response) => void) {
  const url: string = baseURLOld + 'GetGeneralStats';
  const request = {
    gigatoken: getToken(),
  };

  callFetch(url, 'POST', request).then((resp)=>{callback(resp)})
}



function callFetch(url: string, method: string, body: object): Promise<APIResponse> {
  const config: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
      // 'auth-token': GetToken(),
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
