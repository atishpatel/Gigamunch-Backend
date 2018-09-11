export function GetToken() {
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
export function SetToken(cvalue) {
    var jwt = GetJWT(cvalue);
    var d = new Date(0);
    if (jwt) {
        d.setUTCSeconds(jwt.exp);
    }
    document.cookie = "AUTHTKN=" + cvalue + "; expires=" + d.toUTCString() + "; path=/";
    if (location.hostname === 'localhost') {
        window.localStorage.setItem('AUTHTKN', cvalue);
    }
}
export function GetJWT(tkn) {
    if (!tkn) {
        return null;
    }
    var tknConv = tkn.replace(/[+\/]/g, function (m0) {
        return m0 === '+' ? '-' : '_';
    }).replace(/=/g, '');
    var userString = tknConv.split('.')[1].replace(/\s/g, '');
    return JSON.parse(window.atob(userString.replace(/[-_]/g, function (m0) {
        return m0 === '-' ? '+' : '/';
    }).replace(/[^A-Za-z0-9\+\/]/g, '')));
}
