<template>
  <div class="subscribers">
    <SubscribersList :subs="subs"></SubscribersList>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import SubscribersList from '../components/SubscribersList.vue';
import { GetHasSubscribed } from '../ts/service';
import { GetAddressLink, GetAddress } from '../ts/utils';
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
      const subs = resp.subscribers as Types.SubscriberExtended[];
      for (let i = 0; i < subs.length; i++) {
        subs[i].addressString = GetAddress(subs[i].address);
        subs[i].addressLink = GetAddressLink(subs[i].address);
        subs[i].emails = [];
        subs[i].names = [];
        subs[i].email_prefs.reduce((emails, emailPrefs, a, b) => {
          emails.push(emailPrefs.email);
          subs[i].names.push(
            `${emailPrefs.first_name} ${emailPrefs.last_name}`
          );
          return emails;
        }, subs[i].emails);
        subs[i].emailsString = subs[i].emails.toString();
        subs[i].namesString = subs[i].names.toString();

        subs[i].phonenumbers = [];
        subs[i].phone_prefs.reduce((numbers, phonePrefs, a, b) => {
          numbers.push(phonePrefs.number);
          return numbers;
        }, subs[i].phonenumbers);
        subs[i].phonenumbersString = subs[i].phonenumbers.toString();
      }

      this.subs = subs;
    });
  }
}
</script>
<style lang="scss">
</style>
