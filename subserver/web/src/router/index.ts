import Vue from 'vue';
import VueRouter from 'vue-router';
import Dinners from '../views/Dinners.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'dinners',
    component: Dinners,
  },
  {
    path: '/dinner/:date',
    name: 'dinner',
    component: () => import(/* webpackChunkName: "dinner" */ '../views/Dinner.vue'),
  },
  {
    path: '/history',
    name: 'history',
    component: () => import(/* webpackChunkName: "history" */ '../views/History.vue'),
  },
  {
    path: '/account',
    name: 'account',
    component: () => import(/* webpackChunkName: "account" */ '../views/Account.vue'),
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
