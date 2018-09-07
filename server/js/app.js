function GetURLParmas() {
    let vars = {};
    let parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, (m, key, value) => {
        vars[key] = value;
    });
    return vars;
}


var utils = Object.freeze({
	GetURLParmas: GetURLParmas
});

function SetToken(cvalue) {
    const jwt = GetJWT(cvalue);
    const d = new Date(0);
    if (jwt) {
        d.setUTCSeconds(jwt.exp);
    }
    document.cookie = `AUTHTKN=${cvalue}; expires=${d.toUTCString()}; path=/`;
    if (location.hostname === 'localhost') {
        window.localStorage.setItem('AUTHTKN', cvalue);
    }
}
function GetJWT(tkn) {
    if (!tkn) {
        return null;
    }
    const tknConv = tkn.replace(/[+\/]/g, (m0) => {
        return m0 === '+' ? '-' : '_';
    }).replace(/=/g, '');
    const userString = tknConv.split('.')[1].replace(/\s/g, '');
    return JSON.parse(window.atob(userString.replace(/[-_]/g, (m0) => {
        return m0 === '-' ? '+' : '/';
    }).replace(/[^A-Za-z0-9\+\/]/g, '')));
}

function Login(token) {
    const url = '/api/v1/Login';
    const req = {
        token,
    };
    return callFetch(url, 'POST', req).then((resp) => {
        if (resp && resp.token) {
            SetToken(resp.token);
        }
        return resp;
    });
}
function callFetch(url, method, body) {
    return fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    }).then((resp) => {
        return resp.json();
    }).catch((err) => {
        console.error('failed to callFetch', err);
        console.error('details: ', err.code, err.name, err.message, err.detail);
    });
}


var service = Object.freeze({
	Login: Login
});

function SignOut() {
    firebase.auth().signOut();
}
function SetupFirebase() {
    let config;
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
        signInSuccessUrl: 'sub',
        signInOptions: [
            firebase.auth.GoogleAuthProvider.PROVIDER_ID,
            {
                provider: firebase.auth.FacebookAuthProvider.PROVIDER_ID,
                scopes: [
                    'public_profile',
                    'email',
                    'user_likes',
                    'user_friends',
                ],
            },
            firebase.auth.EmailAuthProvider.PROVIDER_ID,
        ],
    };
    const ui = new firebaseui.auth.AuthUI(firebase.auth());
    ui.start(elementID, uiConfig);
}
const EventSignedOut = 'signed-out';
const EventSignedIn = 'signed-in';
SetupFirebase();
firebase.auth().onAuthStateChanged((user) => {
    console.log('user', user);
    let eventName;
    if (!user) {
        eventName = EventSignedOut;
    }
    else {
        eventName = EventSignedIn;
        user.getIdToken(false).then((idToken) => {
            console.log('login in');
            Login(idToken).then((resp) => {
                console.log('login resp: ', resp);
            });
        });
        APP.User = user;
        const event = document.createEvent('Event');
        event.initEvent(eventName, true, true);
        window.dispatchEvent(event);
    }
});


var auth = Object.freeze({
	SignOut: SignOut,
	SetupFirebase: SetupFirebase,
	SetupFirebaseAuthUI: SetupFirebaseAuthUI,
	EventSignedOut: EventSignedOut,
	EventSignedIn: EventSignedIn
});

APP.Utils = utils;
APP.Service = service;
APP.Auth = auth;
console.log('app.js loaded');
