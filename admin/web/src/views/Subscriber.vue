<template>
  <div class="subscriber">
    <SubscriberSummary
      :sub="sub"
      v-on:get-subscriber="getSubscriber"
    ></SubscriberSummary>
    <SubscriberDiscountsList
      class="list"
      :discounts="discounts"
      :sub="sub"
      v-on:dialog-success="getDiscounts"
      v-on:get-subscriber="getSubscriber"
    ></SubscriberDiscountsList>
    <SubscriberActivitiesList
      class="list"
      :activities="acts"
      :sub="sub"
      v-on:get-activities="getActivities"
    ></SubscriberActivitiesList>
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
  protected acts: Types.ActivityExtended[];
  protected logs: Types.LogExtended[];
  protected id: string;

  public constructor() {
    super();
    this.sub = {
      address: {},
    } as Types.SubscriberExtended;
    this.acts = [];
    this.discounts = [];
    this.logs = [];
    const tmp = window.location.pathname.split('/subscriber/');
    const idOrEmail = decodeURIComponent(tmp[1]);
    this.id = idOrEmail;
  }

  public created() {
    this.getSubscriber();
    this.getActivities();
    this.getLogs();
    this.getDiscounts();
  }

  public getActivities() {
    GetSubscriberActivities(this.id).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }

      this.acts = GetActivitiesExtended(resp.activities);
    });
  }

  public getDiscounts() {
    GetSubscriberDiscounts(this.id).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }

      this.discounts = resp.discounts;
    });
  }

  public getSubscriber() {
    GetSubscriber(this.id).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }

      this.sub = GetSubscriberExtended(resp.subscriber);
    });
  }

  public getLogs() {
    GetLogsForUser(0, 1000, this.id).then((resp) => {
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
.subscriber {
  background-color: white;
  max-width: 1500px;
  margin: auto;
  padding: 24px;
}

.list {
  padding: 24px 0;
  // border: 1px solid #dadce0;
  // border-radius: 8px;
  // overflow: hidden;
}
</style>
