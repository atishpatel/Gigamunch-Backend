<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      fixed
      clipped
      app
    >
      <div class="drawer-nav-header">Menu</div>
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
          to="/logout"
          v-if="userSummary.has_subscribed === true"
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
          to="/checkout"
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
          to="/login"
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
      class="app-toolbar"
    >
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>
        <span class="logo-headline">Gigamunch</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <div>
        <a
          v-if="userSummary.has_subscribed === false"
          class="nav-link"
          href="/checkout"
        > Sign up </a>
        <a
          v-if="userSummary.is_logged_in === false"
          class="nav-link"
          href="/login"
        > Login </a>
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
      this.userSummary = resp;
      // console.log(resp);
    });
  }
}
</script>

<style lang="scss">
v-app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

.app-toolbar div {
  background: white;
}

.main {
  background: white;
}

.drawer-nav-header {
  background-color: #009688;
  color: white;
  font-weight: bold;
  font-size: 24px;
  padding: 24px 20px;
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
