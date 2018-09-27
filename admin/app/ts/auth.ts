
declare var APP: any;
declare var app: any;
declare var firebase: any;

export const Events = {
  UserUpdated: 'user-updated',
  SignedOut: 'signed-out',
  SignedIn: 'signed-in',
}

let userLoaded = false;

function setUser(user: FBUser | null) {
  APP.User = user;
  if (user) {
    user.getIdTokenResult(false).then((tokenResult) => {
      let adminClaim = tokenResult.claims['admin'];
      if (adminClaim) {
        user.Admin = true;
      } else {
        user.Admin = false;
      }
      user.IsAdmin = function () {
        return user.Admin;
      }
      APP.User = user;
    })
  }
  const event = document.createEvent('Event');
  event.initEvent(Events.UserUpdated, true, true);
  window.dispatchEvent(event);
  console.log('user', user);
  userLoaded = true;
}

export function GetUser(): Promise<FBUser> {
  return new Promise((resolve, reject) => {
    if (userLoaded) {
      resolve(APP.User);
    }

    const unsubscribe = firebase.auth().onAuthStateChanged((user) => {
      if (!APP.User) {
        setUser(user);
      }
      resolve(APP.User);
    }, reject);
  });
}

export function GetToken(): Promise<string> {
  return GetUser().then((user) => {
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
  } else {
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

// Called when user signs in or signs out
firebase.auth().onAuthStateChanged((user: FBUser) => {
  let eventName: string;
  if (!user) {
    // isn't signed in
    eventName = Events.SignedOut;
  } else {
    // is signed in
    eventName = Events.SignedIn;
    setUser(user);
  }
  // fire event
  const event = document.createEvent('Event');
  event.initEvent(eventName, true, true);
  window.dispatchEvent(event);
});

interface FBUser {
  Admin: boolean
  getIdToken(frocerefresh: boolean): Promise<any>
  getIdTokenResult(forcerefresh: boolean): Promise<any>
  IsAdmin(): boolean
}