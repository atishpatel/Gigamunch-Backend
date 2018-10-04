import firebase from 'firebase/app';
import 'firebase/auth';
import { IsProd } from './env';

export const Events = {
  UserUpdated: 'user-updated',
  SignedOut: 'signed-out',
  SignedIn: 'signed-in',
};

let userLoaded = false;
let usr: firebase.User | null = null;

function fireUserUpdated() {
  const event = document.createEvent('Event');
  event.initEvent(Events.UserUpdated, true, true);
  window.dispatchEvent(event);
}

function setUser(user: firebase.User | null) {
  usr = user;
  fireUserUpdated();
  userLoaded = true;
}

export function SignOut() {
  firebase.auth().signOut();
}

export function SetupFirebase() {
  let config;
  if (IsProd()) {
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
firebase.auth().onAuthStateChanged((user) => {
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

export function GetUser(): Promise<firebase.User | null> {
  return new Promise((resolve: (FBUser: firebase.User | null) => void, reject: () => void) => {
    if (userLoaded) {
      resolve(usr);
    }

    const unsubscribe = firebase.auth().onAuthStateChanged((user) => {
      if (!usr) {
        setUser(user);
      }
      unsubscribe();
      resolve(usr);
    }, reject);
  });
}

export function IsAdmin(): Promise<boolean> {
  return GetUser().then((user) => {
    if (!user) {
      return false;
    }
    return user.getIdTokenResult(false).then((tokenResult) => {
      const adminClaim = tokenResult.claims.admin;
      if (adminClaim) {
        return true;
      }
      return false;
    });
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