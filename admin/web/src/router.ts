import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/subscribers',
      name: 'subscribers',
      component: () => import(/* webpackChunkName: "subscribers" */ './views/Subscribers.vue'),
    },
    {
      path: '/subscriber/:id',
      name: 'subscriber',
      component: () => import(/* webpackChunkName: "subscriber" */ './views/Subscriber.vue'),
    },
    {
      path: '/unpaid-summary',
      name: 'unpaid-summary',
      component: () => import(/* webpackChunkName: "unpaid-summary" */ './views/UnpaidSummary.vue'),
    },
  ],
});
