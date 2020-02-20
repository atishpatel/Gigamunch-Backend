<template>
  <div>
    <div class="content-container">
      <h1 v-if="userSummary.is_active === true">Account</h1>
      <h1
        style="color: #869995;"
        v-if="userSummary.is_active === false"
      >Inactive Account</h1>
      <v-btn
        v-if="userSummary.is_active === false"
        depressed
        color="#E8554E"
        class="white--text"
        @click="reactivateClicked"
      >Re-activate Account</v-btn>
      <div class="list">
        <AccountChangeName
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        ></AccountChangeName>
        <AccountUpdatePayment
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        ></AccountUpdatePayment>
        <AccountChangeServings
          v-on:get-account-info="getAccountInfo"
          :sub="accountInfo.subscriber"
        >
        </AccountChangeServings>
        <AccountChangeDeliveryDay
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        >
        </AccountChangeDeliveryDay>
        <AccountUpdateAddress
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        ></AccountUpdateAddress>
        <AccountChangeDeliveryNotes
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        ></AccountChangeDeliveryNotes>
        <AccountChangePhoneNumber
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        ></AccountChangePhoneNumber>
      </div>
      <div class="cancel">
        <CancelButton
          v-if="userSummary.is_active === true"
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
          v-on:get-user-summary="getUserSummary"
        ></CancelButton>
      </div>

      <div class="footer-message">
        <p class="footer-message-text">If you have any questions, feel free to talk to us at</p>
        <p class="footer-message-text"><a href="mailto:hello@eatgigamunch.com"><strong>hello@eatgigamunch.com</strong></a></p>
        <p
          class="footer-message-text"
          style="margin-top: 12px;"
        ><strong>We're here for you.</strong></p>
        <p
          class="footer-message-text"
          style="margin-top: 32px;"
        >💛&nbsp;&nbsp;The Gigamunch Team</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Prop, Component, Vue } from 'vue-property-decorator';
import { GetAccountInfo } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import { GetUserSummary } from '../ts/service';
import AccountListItem from '../components/AccountListItem.vue';
import AccountChangeServings from '../components/AccountChangeServings.vue';
import AccountChangeName from '../components/AccountChangeName.vue';
import AccountChangeDeliveryNotes from '../components/AccountChangeDeliveryNotes.vue';
import AccountUpdatePayment from '../components/AccountUpdatePayment.vue';
import AccountChangePhoneNumber from '../components/AccountChangePhoneNumber.vue';
import AccountChangeDeliveryDay from '../components/AccountChangeDeliveryDay.vue';
import AccountUpdateAddress from '../components/AccountUpdateAddress.vue';
import CancelButton from '../components/CancelButton.vue';
import DialogConfirm from '../components/DialogConfirm.vue';
import { ActivateSubscriber } from '../ts/service';

@Component({
  components: {
    AccountListItem,
    AccountChangeServings,
    AccountChangeName,
    AccountChangeDeliveryNotes,
    AccountUpdatePayment,
    AccountChangePhoneNumber,
    AccountChangeDeliveryDay,
    AccountUpdateAddress,
    CancelButton,
    DialogConfirm,
  },
})
export default class Account extends Vue {
  @Prop()
  public userSummary!: SubAPI.GetUserSummaryResp;

  public accountInfo!: SubAPI.GetAccountInfoResp;
  protected loading!: boolean;

  public constructor() {
    super();
    this.accountInfo = {
      subscriber: {} as Common.Subscriber,
      payment_info: {} as SubAPI.PaymentInfo,
    } as SubAPI.GetAccountInfoResp;
  }

  public created() {
    this.getAccountInfo();
    window.scrollTo(0, 0);
  }

  public getAccountInfo() {
    this.loading = true;
    GetAccountInfo().then((resp) => {
      this.loading = false;
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      this.accountInfo = resp;
    });
  }
  public getUserSummary() {
    GetUserSummary().then((resp) => {
      this.userSummary = resp;
      // console.log(resp);
    });
  }
  protected reactivateClicked() {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      this.getAccountInfo();
      this.getUserSummary();
    };
    if (!this.accountInfo || !this.accountInfo.subscriber) {
      alert('failed to load account info when reactivating');
      return;
    }
    ActivateSubscriber('').then(handler);
  }
}
</script>
<style scoped lang="scss">
.content-container {
  max-width: 600px;
  margin: auto;
  padding: 12px;
}
.list {
  margin: 40px 0 40px 0;
}
.cancel {
  text-align: center;
}
.divider-line {
  margin: 40px 40px 40px 0;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}
.footer-message {
  padding: 69px 0 50px 0;
  align-content: center;
}

.footer-message-text {
  align-content: center;
  text-align: center;
  margin: 0;
  font-size: 16px;
}
</style>
