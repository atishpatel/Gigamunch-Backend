function GetURLParmas() {
    var vars = {};
    window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function (m, key, value) {
        vars[key] = value;
        return value;
    });
    return vars;
}

var utils = /*#__PURE__*/Object.freeze({
    GetURLParmas: GetURLParmas
});

function SetToken(cvalue) {
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
function GetJWT(tkn) {
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

function Login(token) {
    var url = '/api/v1/Login';
    var req = {
        token: token,
    };
    return callFetch(url, 'POST', req).then(function (resp) {
        if (resp && resp.token) {
            SetToken(resp.token);
        }
        return resp;
    });
}
function callFetch(url, method, body) {
    return fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    }).then(function (resp) {
        return resp.json();
    }).catch(function (err) {
        console.error('failed to callFetch', err);
    });
}

var service = /*#__PURE__*/Object.freeze({
    Login: Login
});

function SignOut() {
    firebase.auth().signOut();
}
function SetupFirebase() {
    var config;
    if (APP.IsProd) {
        config = {
            apiKey: 'AIzaSyC-1vqT4YIKXVmrGkaoVSj1BJnm48NxlT0',
            authDomain: 'gigamunch-omninexus.firebaseapp.com',
            databaseURL: 'https://gigamunch-omninexus.firebaseio.com',
            projectId: 'gigamunch-omninexus',
            storageBucket: 'gigamunch-omninexus.appspot.com',
            messagingSenderId: '837147123677',
        };
    }
    else {
        config = {
            apiKey: 'AIzaSyBHPe4B4k72ljnBXszda6AJGGg21YqEJ4g',
            authDomain: 'gigamunch-omninexus-dev.firebaseapp.com',
            databaseURL: 'https://gigamunch-omninexus-dev.firebaseio.com',
            projectId: 'gigamunch-omninexus-dev',
            storageBucket: 'gigamunch-omninexus-dev.appspot.com',
            messagingSenderId: '108585202286',
        };
    }
    firebase.initializeApp(config);
}
function SetupFirebaseAuthUI(elementID) {
    var uiConfig = {
        tosUrl: '/terms',
        privacyPolicyUrl: '/privacy',
        signInSuccessUrl: 'login',
        signInOptions: [
            firebase.auth.EmailAuthProvider.PROVIDER_ID,
        ],
    };
    var ui = new firebaseui.auth.AuthUI(firebase.auth());
    ui.start(elementID, uiConfig);
}
var Events = {
    SignedOut: 'signed-out',
    SignedIn: 'signed-in',
};
SetupFirebase();
firebase.auth().onAuthStateChanged(function (user) {
    console.log('user', user);
    var eventName;
    if (!user) {
        eventName = Events.SignedOut;
    }
    else {
        eventName = Events.SignedIn;
        user.getIdToken(false).then(function (idToken) {
            Login(idToken);
        });
        APP.User = user;
        var event_1 = document.createEvent('Event');
        event_1.initEvent(eventName, true, true);
        window.dispatchEvent(event_1);
    }
});

var auth = /*#__PURE__*/Object.freeze({
    SignOut: SignOut,
    SetupFirebase: SetupFirebase,
    SetupFirebaseAuthUI: SetupFirebaseAuthUI,
    Events: Events
});

APP.Utils = utils;
APP.Service = service;
APP.Auth = auth;
console.log('app.js loaded');
