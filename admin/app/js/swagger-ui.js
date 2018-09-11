"use strict";
var authTkn = GetToken();
setTimeout(function () {
    if (ui) {
        console.log("set auth-token");
        ui.preauthorizeApiKey("auth-token", authTkn);
    }
}, 3000);
function GetToken() {
    var name = 'AUTHTKN=';
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length).replace(/\n/g, '');
        }
    }
    if (location.hostname === 'localhost') {
        var tnk = window.localStorage.getItem('AUTHTKN');
        if (!tnk) {
            return '';
        }
        return tnk;
    }
    return '';
}
