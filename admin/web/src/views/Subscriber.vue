<template>
  <div class="subscriber">
    <SubscriberSummary :sub="sub"></SubscriberSummary>
    <SubscriberActivitiesList :activities="acts"></SubscriberActivitiesList>
    <SubscriberLogs :logs="logs"></SubscriberLogs>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import SubscriberActivitiesList from '../components/SubscriberActivitiesList.vue';
import SubscriberSummary from '../components/SubscriberSummary.vue';
import SubscriberLogs from '../components/SubscriberLogs.vue';
import {
  GetSubscriber,
  GetSubscriberActivities,
  GetLogsForUser,
} from '../ts/service';
import { GetAddressLink, GetAddress } from '../ts/utils';
import { IsError } from '../ts/errors';
import { GetDayFullDate } from '../ts/utils';

@Component({
  components: {
    SubscriberActivitiesList,
    SubscriberSummary,
    SubscriberLogs,
  },
})
export default class Subscriber extends Vue {
  protected sub: Types.SubscriberExtended;
  protected acts: Types.ActivitiyExtended[];
  protected logs: Types.LogExtended[];

  public constructor() {
    super();
    this.sub = {} as Types.SubscriberExtended;
    this.acts = [];
    this.logs = [];
  }

  public created() {
    const tmp = window.location.pathname.split('/subscriber/');
    const idOrEmail = decodeURIComponent(tmp[1]);
    this.getSubscriber(idOrEmail);
    this.getActivities(idOrEmail);
    this.getLogs(idOrEmail);
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

  public getLogs(idOrEmail: string) {
    GetLogsForUser(0, 1000, idOrEmail).then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      this.logs = resp.logs as Types.LogExtended[];
    });
  }
}
</script>
<style lang="scss">
</style>
