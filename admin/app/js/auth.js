export var Events = {
    UserUpdated: 'user-updated',
    SignedOut: 'signed-out',
    SignedIn: 'signed-in',
};
var userLoaded = false;
function fireUserUpdated() {
    var event = document.createEvent('Event');
    event.initEvent(Events.UserUpdated, true, true);
    window.dispatchEvent(event);
}
function setUser(user) {
    APP.User = user;
    if (user) {
        user.getIdTokenResult(false)
            .then(function (tokenResult) {
            var adminClaim = tokenResult.claims['admin'];
            if (adminClaim) {
                user.Admin = true;
            }
            else {
                user.Admin = false;
            }
            user.IsAdmin = function () {
                return user.Admin;
            };
            APP.User = user;
        })
            .then(function () {
            console.log('user', user);
            fireUserUpdated();
        });
    }
    else {
        console.log('user', user);
        fireUserUpdated();
    }
    userLoaded = true;
}
export function IsAdmin() {
    return GetUser().then(function (user) {
        if (!user) {
            return false;
        }
        return user.getIdTokenResult(false).then(function (tokenResult) {
            var adminClaim = tokenResult.claims['admin'];
            if (adminClaim) {
                return true;
            }
            return false;
        });
    });
}
export function GetUser() {
    return new Promise(function (resolve, reject) {
        if (userLoaded) {
            resolve(APP.User);
        }
        var unsubscribe = firebase.auth().onAuthStateChanged(function (user) {
            if (!APP.User) {
                setUser(user);
            }
            resolve(APP.User);
        }, reject);
    });
}
export function GetToken() {
    return GetUser().then(function (user) {
        if (!user) {
            return '';
        }
        return user.getIdToken(false);
    });
}
export function SignOut() {
    firebase.auth().signOut();
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
SetupFirebase();
firebase.auth().onAuthStateChanged(function (user) {
    var eventName;
    if (!user) {
        eventName = Events.SignedOut;
    }
    else {
        eventName = Events.SignedIn;
        setUser(user);
    }
    var event = document.createEvent('Event');
    event.initEvent(eventName, true, true);
    window.dispatchEvent(event);
});
