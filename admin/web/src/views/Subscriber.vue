<template>
  <div class="subscriber">
    Subscriber
    <SubscriberActivitiesList :activities="acts"></SubscriberActivitiesList>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import SubscriberActivitiesList from '../components/SubscriberActivitiesList.vue';
import SubscriberSummary from '../components/SubscriberSummary.vue';
import { GetSubscriber, GetSubscriberActivities } from '../ts/service';
import { GetAddressLink, GetAddress } from '../ts/utils';
import { IsError } from '../ts/errors';
import { GetDayFullDate } from '../ts/utils';

@Component({
  components: {
    SubscriberActivitiesList,
    SubscriberSummary,
  },
})
export default class Subscriber extends Vue {
  protected sub: Types.SubscriberExtended;
  protected acts: Types.ActivitiyExtended[];

  public constructor() {
    super();
    this.sub = {} as Types.SubscriberExtended;
    this.acts = [];
  }

  public created() {
    const tmp = window.location.pathname.split('/subscriber/');
    const idOrEmail = decodeURIComponent(tmp[1]);
    this.getSubscriber(idOrEmail);
    this.getActivities(idOrEmail);
  }

  public getActivities(idOrEmail: string) {
    GetSubscriberActivities(idOrEmail).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      const acts = resp.activities as Types.ActivitiyExtended[];
      for (let i = 0; i < acts.length; i++) {
        acts[i].dateFull = GetDayFullDate(acts[i].date);
      }

      this.acts = acts;
    });
  }

  public getSubscriber(idOrEmail: string) {
    GetSubscriber(idOrEmail).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      const sub = resp.subscriber as Types.SubscriberExtended;
      sub.addressString = GetAddress(sub.address);
      sub.addressLink = GetAddressLink(sub.address);
      sub.emails = [];
      sub.names = [];
      sub.email_prefs.reduce((emails, emailPrefs, a, b) => {
        emails.push(emailPrefs.email);
        sub.names.push(`${emailPrefs.first_name} ${emailPrefs.last_name}`);
        return emails;
      }, sub.emails);
      sub.emailsString = sub.emails.toString();
      sub.namesString = sub.names.toString();

      sub.phonenumbers = [];
      sub.phone_prefs.reduce((numbers, phonePrefs, a, b) => {
        numbers.push(phonePrefs.number);
        return numbers;
      }, sub.phonenumbers);
      sub.phonenumbersString = sub.phonenumbers.toString();

      this.sub = sub;
    });
  }
}
</script>
<style lang="scss">
</style>
