<template>
  <mdc-layout-app id="app">
    <mdc-drawer
      slot="drawer"
      temporary
      toggle-on="toggle-drawer"
    >
      <!-- Drawer -->
      <div class="drawer-top">
        <div class="drawer-top-logo">Gigamunch</div>
      </div>
      <mdc-drawer-list>
        <mdc-drawer-item
          to="/"
          start-icon="inbox"
        >
          My Dinners
        </mdc-drawer-item>
        <mdc-drawer-item
          to="/history"
          start-icon="send"
        >
          Dinner History
        </mdc-drawer-item>
        <mdc-drawer-item
          to="/account"
          start-icon="drafts"
        >
          Account
        </mdc-drawer-item>
      </mdc-drawer-list>
    </mdc-drawer>

    <mdc-toolbar slot="toolbar">
      <!-- Toolbar -->
      <mdc-toolbar-row>
        <mdc-toolbar-section align-start>
          <mdc-toolbar-menu-icon event="toggle-drawer"></mdc-toolbar-menu-icon>
          <mdc-toolbar-title>Gigamunch</mdc-toolbar-title>
        </mdc-toolbar-section>
        <mdc-toolbar-section align-end>
          <a
            v-if="userSummary.has_subscribed === true && userSummary.is_active === false"
            class="nav-link"
            href="account"
          >Sign up</a>
          <a
            v-if="userSummary.has_subscribed === false"
            class="nav-link"
            href="/checkout"
          >Sign up</a>
          <a
            v-if="userSummary.is_logged_in === false"
            class="nav-link"
            href="/login"
          >Login</a>
        </mdc-toolbar-section>
      </mdc-toolbar-row>
    </mdc-toolbar>

    <main>
      <div
        class="main-loading-screen"
        :hidden-fade="hideLoadingScreen"
      >
        TODO: Add loading animation
      </div>
      <!-- Main Content -->
      <router-view :userSummary="userSummary" />
    </main>
  </mdc-layout-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { GetUserSummary } from './ts/service';

@Component({})
export default class App extends Vue {
  public drawerOpen = false;
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
      // this.userSummary = resp;
      console.log(resp);
    });
  }
}
</script>
<style lang="scss">
// global
@import 'scss/theme';
@import 'scss/shared-styles';

#app {
  font-family: 'Roboto', 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $mdc-theme-on-primary;
}

.main-loading-screen {
  background-color: pink;
  height: 100%;
  width: 100%;
  z-index: 1000;
  position: absolute;
  top: 0;
  left: 0;
  margin: auto;
  opacity: 1;
  transition: visibility 0s 2s, opacity 2s ease-in-out;
}
.main-loading-screen[hidden-fade] {
  opacity: 0;
  visibility: hidden;
}

.drawer-top {
  min-height: 200px;
  background-color: $mdc-theme-accent;
}
.drawer-top-logo {
  padding: 12px;
}
</style>
<style lang="scss" scoped>
.nav-link {
  padding: 12px;
}
</style>
