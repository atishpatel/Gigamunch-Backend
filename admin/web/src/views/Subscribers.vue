<template>
  <div class="subscribers">
    <SubscribersList :subs="subs"></SubscribersList>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import SubscribersList from '../components/SubscribersList.vue';
import { GetHasSubscribed } from '../ts/service';
import { GetSubscribersExtended } from '../ts/extended';
import { IsError } from '../ts/errors';

@Component({
  components: {
    SubscribersList,
  },
})
export default class Subscribers extends Vue {
  protected subs: Types.SubscriberExtended[];

  public constructor() {
    super();
    this.subs = [];
  }

  public created() {
    this.getSubscribers();
  }

  public getSubscribers() {
    GetHasSubscribed(0, 10000).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      this.subs = GetSubscribersExtended(resp.subscribers);
    });
  }
}
</script>
<style lang="scss">
</style>
