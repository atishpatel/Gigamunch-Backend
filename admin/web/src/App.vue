<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      fixed
      clipped
      app
    >
      <v-list>
        <v-list-tile to="subscribers">
          <v-list-tile-action>
            <v-icon></v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>
              Subscribers
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
      <v-toolbar-title class="headline">
        <span>Gigamunch Admin</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>

    <v-content>
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
</style>
