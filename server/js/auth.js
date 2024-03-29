import { Login } from './service';
export function SignOut(callback) {
    firebase.auth().signOut().then(callback);
}
export function SetupFirebase() {
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
export function SetupFirebaseAuthUI(elementID) {
    var uiConfig = {
        tosUrl: '/terms',
        privacyPolicyUrl: '/privacy',
        signInSuccessUrl: '/login',
        signInOptions: [{
                provider: firebase.auth.EmailAuthProvider.PROVIDER_ID,
                requireDisplayName: false,
            },
        ],
        credentialHelper: firebaseui.auth.CredentialHelper.NONE,
    };
    var ui = new firebaseui.auth.AuthUI(firebase.auth());
    ui.start(elementID, uiConfig);
}
export var Events = {
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
            Login(idToken).then(function () {
                window.location.href = '/sub/';
            });
        });
        APP.User = user;
        var event_1 = document.createEvent('Event');
        event_1.initEvent(eventName, true, true);
        window.dispatchEvent(event_1);
    }
});
