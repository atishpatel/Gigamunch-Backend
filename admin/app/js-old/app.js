var baseURLOld = 'https://cookapi-dot-gigamunch-omninexus.appspot.com/_ah/spi/Service.';
if (APP.IsDev) {
    baseURLOld = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/spi/Service.';
}
else if (APP.IsStage) {
    baseURLOld = 'https://cookapi-dot-gigamunch-omninexus-dev.appspot.com/_ah/spi/Service.';
}
function getToken() {
    return COOK.User.token;
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
function getSubLogs(callback) {
    var url = baseURLOld + 'getSubLogs';
    var request = {
        gigatoken: getToken(),
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('getSubLogs', resp.err);
        callback(resp.sublogs, resp.err);
    });
}
function getSubLogsForDate(date, callback) {
    var url = baseURLOld + 'getSubLogsForDate';
    var request = {
        gigatoken: getToken(),
        date: date.toISOString(),
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('getSubLogsForDate', resp.err);
        callback(resp.sublogs, resp.err);
    });
}
function getSubEmails(callback) {
    var url = baseURLOld + 'getSubEmails';
    var request = {
        gigatoken: getToken(),
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('getSubEmails', resp.err);
        callback(resp.sub_emails, resp.err);
    });
}
function getSubEmailsAndSubs(callback) {
    var url = baseURLOld + 'getSubEmailsAndSubs';
    var request = {
        gigatoken: getToken(),
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('getSubEmails', resp.err);
        callback(resp.sub_emails, resp.subscribers, resp.err);
    });
}
function skipSubLog(date, subEmail, callback) {
    var url = baseURLOld + 'skipSubLog';
    var request = {
        gigatoken: getToken(),
        date: date.toISOString(),
        sub_email: subEmail,
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('skipSubLog', resp.err);
        callback(resp.err);
    });
}
function CancelSub(email, callback) {
    var url = baseURLOld + 'CancelSub';
    var request = {
        gigatoken: getToken(),
        email: email,
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('CancelSub', resp.err);
        callback(resp.err);
    });
}
function discountSubLog(date, subEmail, amount, percent, overrideDiscount, callback) {
    var url = baseURLOld + 'DiscountSubLog';
    var request = {
        gigatoken: getToken(),
        date: date.toISOString(),
        sub_email: subEmail,
        amount: amount,
        percent: percent,
        override_discount: overrideDiscount,
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('DiscountSubLog', resp.err);
        callback(resp.err);
    });
}
function ChangeServingsForDate(date, subEmail, servings, callback) {
    var url = baseURLOld + 'ChangeServingsForDate';
    var request = {
        gigatoken: getToken(),
        date: date.toISOString(),
        sub_email: subEmail,
        servings: servings,
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('ChangeServingForDate', resp.err);
        callback(resp.err);
    });
}
function ChangeServingsPermanently(email, servings, vegetarian, callback) {
    var url = baseURLOld + 'ChangeServingsPermanently';
    var request = {
        gigatoken: getToken(),
        email: email,
        servings: servings,
        vegetarian: vegetarian,
    };
    callFetch(url, 'POST', request).then(function (resp) {
        logError('ChangeServingsPermanently', resp.err);
        callback(resp.err);
    });
}
function GetGeneralStats(callback) {
    var url = baseURLOld + 'GetGeneralStats';
    var request = {
        gigatoken: getToken(),
    };
    callFetch(url, 'POST', request).then(function (resp) { callback(resp); });
}
function callFetch(url, method, body) {
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
        console.error('failed to callFetch', err);
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

var Service = /*#__PURE__*/Object.freeze({
    getSubLogs: getSubLogs,
    getSubLogsForDate: getSubLogsForDate,
    getSubEmails: getSubEmails,
    getSubEmailsAndSubs: getSubEmailsAndSubs,
    skipSubLog: skipSubLog,
    CancelSub: CancelSub,
    discountSubLog: discountSubLog,
    ChangeServingsForDate: ChangeServingsForDate,
    ChangeServingsPermanently: ChangeServingsPermanently,
    GetGeneralStats: GetGeneralStats
});

COOK.Service = Service;
