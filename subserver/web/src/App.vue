<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      fixed
      clipped
      app
    >
      <div class="drawer-nav-header">Gigamunch</div>
      <v-list>
        <v-list-tile to="/">
          <v-list-tile-action>
            <v-icon>calendar_today</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Upcoming Dinners
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <!-- <v-list-tile to="/history">
          <v-list-tile-action>
            <v-icon>drafts</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Dinner History
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile> -->
        <v-list-tile
          to="/account"
          v-if="userSummary.has_subscribed === true"
        >
          <v-list-tile-action>
            <v-icon>account_circle</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Account
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile
          v-if="userSummary.is_logged_in === true"
          href="#"
          @click="signOut"
        >
          <v-list-tile-action>
            <v-icon>logout</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Log out
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile
          href="https://eatgigamunch.com/checkout"
          target="_blank"
          v-if="userSummary.has_subscribed === false"
        >
          <v-list-tile-action>
            <v-icon>emoji_emotions</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Sign up
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile
          href="/login"
          v-if="userSummary.is_logged_in === false"
        >
          <v-list-tile-action>
            <v-icon>account_circle</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Log in
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar
      fixed
      flat
      clipped-left
      app
      color="white"
      class="app-toolbar"
    >
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>
        <span class="logo-headline">Gigamunch</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <div>
        <v-btn
          v-if="userSummary.has_subscribed === false"
          depressed
          color="#E8554E"
          target="_blank"
          class="white--text"
          href="https://eatgigamunch.com/checkout"
        >Sign up</v-btn>
        <a
          v-if="userSummary.is_logged_in === false"
          class="nav-link"
          href="/login"
        >Login</a>
      </div>
    </v-toolbar>

    <v-content class="main">
      <router-view :userSummary="userSummary" />

    </v-content>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { GetUserSummary } from './ts/service';
import { SignOut } from './ts/auth';

@Component({})
export default class App extends Vue {
  public drawer = false;
  public hideLoadingScreen = false;
  public userSummary = {
    is_active: false,
    is_logged_in: false,
    has_subscribed: false,
    is_admin: false,
    on_probation: false,
    error: {} as Common.Error,
  } as SubAPI.GetUserSummaryResp;

  public created() {
    // App ready
    GetUserSummary().then((resp) => {
      this.hideLoadingScreen = true;
      if (resp.is_logged_in) {
        this.userSummary = resp;
      }
      // console.log(resp);
    });
  }

  public signOut() {
    SignOut().then(() => {
      window.location.href = '/sub/';
    });
  }
}
</script>

<style lang="scss">
$body-font-family: 'Poppins', 'Avenir', Helvetica, Arial, sans-serif;
$title-font: 'Laila', 'Poppins', 'Avenir', Helvetica, Arial, sans-serif;

v-app {
  font-family: $body-font-family !important;
  .title {
    // To pin point specific classes of some components
    font-family: $title-font !important;
  }
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

.main {
  background: white;
}

.drawer-nav-header {
  background-color: #d0782c;
  color: white;
  font-weight: bold;
  font-size: 24px;
  padding: 24px 20px;
  font-family: 'Laila', serif;
  font-weight: 500;
}

.nav-tile-content {
  flex-direction: row;
  align-items: center;
}

.theme--light.v-list .v-list__tile--link:hover {
  background-color: #dce0e2;
}

.v-list__tile--active {
  background-color: #eceff1;
}

.logo-headline {
  color: #d0782c;
  font-family: 'Laila', serif;
  font-weight: 500;
}
</style>
