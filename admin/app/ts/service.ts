import { GetToken } from './auth';

let baseURL = '/admin/api/v1/';
if (location.hostname === 'localhost') {
    baseURL = 'https://gigamunch-omninexus-dev.appspot.com/admin/api/v1/';
}

// Subscriber
export function GetSubscriber(email: string): Promise<any> {
    const url: string = baseURL + 'GetSubscriber';
    const req: GetSubscriberReq = {
        email,
    };
    return callFetch(url, 'GET', req);
}

export function GetHasSubscribed(date: Date): Promise<any> {
    const url: string = baseURL + 'GetHasSubscribed';
    const req: GetHasSubscribedReq = {
        date: date.toISOString(),
    };
    return callFetch(url, 'GET', req);
}

// SubLog
export function GetUnpaidSublogs(limit: number): Promise<any> {
    const url: string = baseURL + 'GetUnpaidSublogs';
    const req: GetUnpaidSublogsReq = {
        limit,
    };
    return callFetch(url, 'GET', req);
}

export function GetSubscriberSublogs(email: string): Promise<any> {
    const url: string = baseURL + 'GetSubscriberSublogs';
    const req: GetSubscriberSublogsReq = {
        email,
    };
    return callFetch(url, 'GET', req);
}

export function ProcessSublog(date: string, email: string): Promise<any> {
    const url: string = baseURL + 'ProcessSublog';
    const req: ProcessSublogsReq = {
        date,
        email,
    };
    return callFetch(url, 'POST', req);
}

// Execution
export function GetExecutions(start: number, limit: number): Promise<any> {
    const url: string = baseURL + 'GetExecutions';
    const req: GetExecutionsReq = {
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}

export function GetExecution(idOrDate: string): Promise<any> {
    const url: string = baseURL + 'GetExecution';
    const req: GetExecutionReq = {
        idOrDate,
    };
    return callFetch(url, 'GET', req);
}

export function UpdateExecution(mode: string, execution: Common.Execution): Promise<any> {
    const url: string = baseURL + 'UpdateExecution';
    const req: UpdateExecutionReq = {
        mode,
        execution,
    };
    return callFetch(url, 'POST', req);
}

// Activity
export function GetActivityForDate() { }

// Logs
export function GetLogs(start: number, limit: number): Promise<any> {
    const url: string = baseURL + 'GetLogs';
    const req: GetLogsReq = {
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}

export function GetLog(id: number): Promise<any> {
    const url: string = baseURL + 'GetLog';
    const req: GetLogReq = {
        id,
    };

    return callFetch(url, 'GET', req);
}

export function GetLogsByEmail(start: number, limit: number, email: string): Promise<any> {
    const url: string = baseURL + 'GetLogsByEmail';
    const req: GetLogsByEmailReq = {
        email,
        start,
        limit,
    };
    return callFetch(url, 'GET', req);
}

export function GetLogsByExecution(execution_id: number): Promise<any> {
    const url: string = baseURL + 'GetLogsByExecution';
    const req: GetLogsByExecutionReq = {
        execution_id,
    };
    return callFetch(url, 'GET', req);
}

function callFetch(url: string, method: string, body: object): Promise<APIResponse> {
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
                // if (resp.status == 500 || resp.status == 404) {
                try {
                    return resp.json();
                } catch (err) {
                    return {
                        error: {
                            code: resp.status,
                            message: 'Unknown server error',
                        }
                    }
                }
                // }
                // return resp.json();
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
