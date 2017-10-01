import { Element as PolymerElement } from '../../node_modules/@polymer/polymer/polymer-element.js';
import { Login } from '../service';
import { IsLoggedIn } from '../user';
import template from './login-page.html';

export class LoginPage extends PolymerElement {

  static get template() {
    return template;
  }

  constructor() {
    super();
    this.name = 'Login Page';
  }

  ready() {
    super.ready();
    firebase.auth().onAuthStateChanged((user) => {
      console.log('user', user);
      if (user === null && !IsLoggedIn()) {
        console.log('setup fb login');
        this.setupFBLogin();
      } else {
        user.getToken(false).then((idToken: string) => {
          console.log('signing in');
          Login(idToken).then((resp) => {
            console.log('login resp: ', resp);
            firebase.auth().signOut();
          });
        });
      }
    });
  }

  setupFBLogin() {
    // FirebaseUI config.
    const uiConfig = {
      callbacks: {
        signInSuccess: this.signInSuccess,
      },
      tosUrl: '/terms',
      signInSuccessUrl: 'login',
      signInOptions: [
        // Leave the lines as is for the providers you want to offer your users.
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

    // Initialize the FirebaseUI Widget using Firebase.
    const ui = new firebaseui.auth.AuthUI(firebase.auth());
    // The start method will wait until the DOM is loaded.
    ui.start('#firebaseui-auth-container', uiConfig);
  }

  selected() {
    console.log('login selected');
  }

  signInSuccess(currentUser, credential, redirectUrl) {
    console.log('currentUser:', currentUser);
    console.log('credential:', credential);
    console.log('redirectUrl:', redirectUrl);
    return true;
  }

  static get properties() {
    return {
      name: {
        Type: String,
      },
    };
  }
}

customElements.define('login-page', LoginPage);
