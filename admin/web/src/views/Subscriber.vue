<template>
  <div class="subscriber">
    <SubscriberSummary :sub="sub"></SubscriberSummary>
    <SubscriberDiscountsList :discounts="discounts"></SubscriberDiscountsList>
    <SubscriberActivitiesList :activities="acts"></SubscriberActivitiesList>
    <SubscriberLogs :logs="logs"></SubscriberLogs>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import SubscriberActivitiesList from '../components/SubscriberActivitiesList.vue';
import SubscriberDiscountsList from '../components/SubscriberDiscountsList.vue';
import SubscriberSummary from '../components/SubscriberSummary.vue';
import SubscriberLogs from '../components/SubscriberLogs.vue';
import {
  GetSubscriber,
  GetSubscriberActivities,
  GetSubscriberDiscounts,
  GetLogsForUser,
} from '../ts/service';
import {
  GetSubscriberExtended,
  GetActivitiesExtended,
  GetLogsExtended,
} from '../ts/extended';
import { IsError } from '../ts/errors';

@Component({
  components: {
    SubscriberActivitiesList,
    SubscriberDiscountsList,
    SubscriberSummary,
    SubscriberLogs,
  },
})
export default class Subscriber extends Vue {
  protected sub: Types.SubscriberExtended;
  protected discounts: Common.Discount[];
  protected acts: Types.ActivitiyExtended[];
  protected logs: Types.LogExtended[];

  public constructor() {
    super();
    this.sub = {} as Types.SubscriberExtended;
    this.acts = [];
    this.discounts = [];
    this.logs = [];
  }

  public created() {
    const tmp = window.location.pathname.split('/subscriber/');
    const idOrEmail = decodeURIComponent(tmp[1]);
    this.getSubscriber(idOrEmail);
    this.getActivities(idOrEmail);
    this.getLogs(idOrEmail);
    this.getDiscounts(idOrEmail);
  }

  public getActivities(idOrEmail: string) {
    GetSubscriberActivities(idOrEmail).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }

      this.acts = GetActivitiesExtended(resp.activities);
    });
  }

  public getDiscounts(id: string) {
    GetSubscriberDiscounts(id).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }

      this.discounts = resp.discounts;
    });
  }

  public getSubscriber(idOrEmail: string) {
    GetSubscriber(idOrEmail).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }

      this.sub = GetSubscriberExtended(resp.subscriber);
    });
  }

  public getLogs(idOrEmail: string) {
    GetLogsForUser(0, 1000, idOrEmail).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      this.logs = GetLogsExtended(resp.logs);
    });
  }
}
</script>
<style lang="scss">
</style>
