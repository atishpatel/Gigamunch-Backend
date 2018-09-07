import { Login } from './service';
export function SignOut() {
    firebase.auth().signOut();
}
export function SetupFirebase() {
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
export function SetupFirebaseAuthUI(elementID) {
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
export const EventSignedOut = 'signed-out';
export const EventSignedIn = 'signed-in';
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
