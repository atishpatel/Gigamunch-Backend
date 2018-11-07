<template>
  <mdc-layout-app id="app">
    <mdc-drawer slot="drawer" temporary toggle-on="toggle-drawer">
      <!-- Drawer -->
      <mdc-drawer-list>
        <mdc-drawer-item to="/" start-icon="inbox">
          My Dinners
        </mdc-drawer-item>
        <mdc-drawer-item to="/history" start-icon="send">
          Dinner History
        </mdc-drawer-item>
        <mdc-drawer-item to="/account" start-icon="drafts">
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
      </mdc-toolbar-row>
    </mdc-toolbar>

    <main>
      <!-- Main Content -->
      <router-view />
    </main>
  </mdc-layout-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
// import {GetUser} from './service';

@Component({})
export default class App extends Vue {
  public drawerOpen = false;

  public created() {
    // App ready
  }
}
</script>

<style lang="scss">
@import 'scss/theme';
#app {
  font-family: 'Roboto', 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $mdc-theme-on-primary;
}
// #nav {
//   padding: 30px;
//   a {
//     font-weight: bold;
//     color: #2c3e50;

//     color: #42b983;
//   }
// }
</style>
