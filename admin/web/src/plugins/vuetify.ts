import Vue from 'vue';
import Vuetify from 'vuetify';
import 'vuetify/src/stylus/app.styl';

Vue.use(Vuetify, {
  iconfont: 'md',
  theme: {
    primary: '#009688',
    secondary: '#f44336',
    accent: '#03a9f4',
    error: '#FF5252',
    info: '#2196F3',
    success: '#4CAF50',
    warning: '#FFC107',
  },
});

// @ts-ignore
import VuetifyGoogleAutocomplete from 'vuetify-google-autocomplete/lib';

Vue.use(VuetifyGoogleAutocomplete, {
  apiKey: 'AIzaSyCDOw6QXpThS7dm3rl79wDdEvwPlLWsi0Y',
});
