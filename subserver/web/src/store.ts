import Vue from 'vue';
import Vuex from 'vuex';
import { GetUserSummary } from './ts/service';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    // UserSummaryLoaded: false,
    // UserSummary: {
    //   is_logged_in: false,
    //   has_subscribed: false,
    //   is_active: false,
    //   on_probation: false,
    //   error: {} as Common.Error,
    // },
  },
  mutations: {
    // SetUserSummary(state, userSummary: SubAPI.GetUserSummaryResp) {
    //   state.UserSummaryLoaded = true;
    //   state.UserSummary = userSummary;
    // },
  },
  actions: {

  },
  getters: {
    // GetUserSummary(state): Promise<SubAPI.GetUserSummaryResp> {
    //   if (state.UserSummaryLoaded === true) {
    //     return Promise.resolve(state.UserSummary);
    //   }
    //   return GetUserSummary().then((resp) => {
    //     state.UserSummaryLoaded = true;
    //     state.UserSummary = resp;
    //     return resp;
    //   });
    // },
  },
});

// interface AppState {
  // UserSummaryLoaded: boolean,
  // UserSummary: SubAPI.GetUserSummaryResp,
// }
