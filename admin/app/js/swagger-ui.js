"use strict";
var APP = {};
APP.IsDev = false;
APP.IsStage = false;
APP.IsProd = false;
switch (location.hostname) {
    case '127.0.0.1':
    case 'localhost':
        APP.IsDev = true;
        break;
    case 'gigamunch-omninexus-dev.appspot.com':
        APP.IsStage = true;
        break;
    default:
        APP.IsProd = true;
}
function updateAuthToken() {
    if (ui) {
        console.log("set auth-token");
        APP.Auth.GetToken().then(function (token) {
            ui.preauthorizeApiKey("auth-token", token);
        });
    }
}
setTimeout(updateAuthToken, 3000);
setInterval(updateAuthToken, 10 * 60 * 1000);
