import Vue from 'vue';
import router from './router';
import store from './store';
import App from './App.vue';
import './registerServiceWorker';

import VueMDCAdapter from 'vue-mdc-adapter';
import './scss/theme.scss';

Vue.config.productionTip = false;

Vue.use(VueMDCAdapter);

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
