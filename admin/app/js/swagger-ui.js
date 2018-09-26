"use strict";
setTimeout(function () {
    if (ui) {
        console.log("set auth-token");
        APP.Auth.GetToken().then(function (token) {
            ui.preauthorizeApiKey("auth-token", token);
        });
    }
}, 3000);
