var baseURLOld = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/spi/Service.';
if (APP.IsDev) {
    baseURLOld = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/spi/Service.';
}
else if (APP.IsStage) {
    baseURLOld = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/spi/Service.';
}
function GetToken() {
    return APP.Auth.GetToken();
}
function logError(fnName, err) {
    if (err && (err.code === undefined || err.code !== 0)) {
        var desc = "Function: " + fnName + " | Message: " + err.message + " | Details: " + err.detail;
        console.error(desc);
        ga('send', 'exception', {
            exDescription: desc,
            exFatal: false,
        });
        return true;
    }
    return false;
}
export function getSubLogs(callback) {
    var url = baseURLOld + 'getSubLogs';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('getSubLogs', resp.err);
            callback(resp.sublogs, resp.err);
        });
    });
}
export function getSubLogsForDate(date, callback) {
    var url = baseURLOld + 'getSubLogsForDate';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            date: date.toISOString(),
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('getSubLogsForDate', resp.err);
            callback(resp.sublogs, resp.err);
        });
    });
}
export function getSubEmails(callback) {
    var url = baseURLOld + 'getSubEmails';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('getSubEmails', resp.err);
            callback(resp.sub_emails, resp.err);
        });
    });
}
export function getSubEmailsAndSubs(callback) {
    var url = baseURLOld + 'getSubEmails';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('getSubEmails', resp.err);
            callback(resp.sub_emails, resp.subscribers, resp.err);
        });
    });
}
export function skipSubLog(date, subEmail, callback) {
    var url = baseURLOld + 'skipSubLog';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            date: date.toISOString(),
            sub_email: subEmail,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('skipSubLog', resp.err);
            callback(resp.err);
        });
    });
}
export function CancelSub(email, callback) {
    var url = baseURLOld + 'CancelSub';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            email: email,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('CancelSub', resp.err);
            callback(resp.err);
        });
    });
}
export function discountSubLog(date, subEmail, amount, percent, overrideDiscount, callback) {
    var url = baseURLOld + 'DiscountSubLog';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            date: date.toISOString(),
            sub_email: subEmail,
            amount: amount,
            percent: percent,
            override_discount: overrideDiscount,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('DiscountSubLog', resp.err);
            callback(resp.err);
        });
    });
}
export function ChangeServingsForDate(date, subEmail, servings, callback) {
    var url = baseURLOld + 'ChangeServingsForDate';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            date: date.toISOString(),
            sub_email: subEmail,
            servings: servings,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('ChangeServingForDate', resp.err);
            callback(resp.err);
        });
    });
}
export function ChangeServingsPermanently(email, servings, vegetarian, callback) {
    var url = baseURLOld + 'ChangeServingsPermanently';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            email: email,
            servings: servings,
            vegetarian: vegetarian,
        };
        callOldFetch(url, 'POST', request).then(function (resp) {
            logError('ChangeServingsPermanently', resp.err);
            callback(resp.err);
        });
    });
}
export function GetGeneralStats(start_date_min, start_date_max, callback) {
    var url = baseURLOld + 'GetGeneralStats';
    GetToken().then(function (token) {
        var request = {
            gigatoken: token,
            start_date_min: start_date_min.toISOString(),
            start_date_max: start_date_max.toISOString(),
        };
        callOldFetch(url, 'POST', request).then(function (resp) { callback(resp); });
    });
}
function callOldFetch(url, method, body) {
    var config = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
        },
    };
    var URL = url;
    if (method === 'GET') {
        URL += '?' + serializeParams(body);
    }
    else {
        config.body = JSON.stringify(body);
    }
    return fetch(URL, config)
        .then(function (resp) {
        return resp.json();
    })
        .catch(function (err) {
        console.error('failed to callOldFetch', err);
    });
}
function serializeParams(obj) {
    var str = [];
    var p;
    p = 0;
    for (p in obj) {
        if (obj.hasOwnProperty(p)) {
            var k = p;
            var v = obj[p];
            str.push((v !== null && typeof v === 'object') ?
                serializeParams(v) :
                encodeURIComponent(k) + '=' + encodeURIComponent(v));
        }
    }
    return str.join('&');
}
