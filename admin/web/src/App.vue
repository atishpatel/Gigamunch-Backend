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
        <v-list-tile to="/subscribers">
          <v-list-tile-action>
            <v-icon>people</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Subscribers
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile to="/message">
          <v-list-tile-action>
            <v-icon>message</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Message
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile to="/unpaid-summary">
          <v-list-tile-action>
            <v-icon>sentiment_dissatisfied</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Unpaid Summary
            </v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile
          target="_blank"
          href="/admin/executions"
        >
          <v-list-tile-action>
            <v-icon>description</v-icon>
          </v-list-tile-action>
          <v-list-tile-content class="nav-tile-content">
            <v-list-tile-title>
              Culture Exections
            </v-list-tile-title>
            <v-spacer></v-spacer>
            <v-icon>open_in_new</v-icon>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile
          target="_blank"
          href="/admin/swagger/index.html"
        >
          <v-list-tile-action>
            <v-icon>code</v-icon>
          </v-list-tile-action>
          <v-list-tile-content class="nav-tile-content">
            <v-list-tile-title>
              Admin API
            </v-list-tile-title>
            <v-spacer></v-spacer>
            <v-icon>open_in_new</v-icon>
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
      <v-toolbar-title class="headline">
        <span>Gigamunch Admin</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>

    <v-content class="main">
      <router-view />
    </v-content>
  </v-app>
</template>


<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { IsAdmin } from './ts/auth';

@Component({})
export default class App extends Vue {
  public drawer = false;

  public created() {
    // App ready
    if (!IsAdmin()) {
      alert('User is not admin');
    }
  }
}
</script>
<style lang="scss">
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
</style>
