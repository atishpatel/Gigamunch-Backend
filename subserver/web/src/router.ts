import Vue from 'vue';
import Router from 'vue-router';
import Dinners from './views/Dinners.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'dinners',
      component: Dinners,
    },
    {
      path: '/dinner/:date',
      name: 'dinner',
      component: () => import(/* webpackChunkName: "dinner" */ './views/Dinner.vue'),
    },
    {
      path: '/account',
      name: 'account',
      component: () => import(/* webpackChunkName: "account" */ './views/Account.vue'),
    },
  ],
});
