<template>
  <div>
    <div class="content-container">
      <h1>Account</h1>
      <div class="list">
        <AccountChangeName
          :sub="accountInfo.subscriber"
          v-on:get-account-info="getAccountInfo"
        ></AccountChangeName>
        <AccountListItem
          title="Payment Method"
          value="xxxx-xxxx-xxxx-4444"
        ></AccountListItem>
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
        <AccountListItem
          title="Delivery Address"
          value="1835 North Washington Avenue, Cookeville, TN, 38501"
        ></AccountListItem>
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
        <v-btn
          depressed
          large
          color="#E8554E"
          class="white--text"
        >Cancel Account</v-btn>
      </div>
      <hr class="divider-line">
      <div class="footer-message">
        <p class="footer-message-text">Feel free to talk to us at</p>
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
import AccountListItem from '../components/AccountListItem.vue';
import AccountChangeServings from '../components/AccountChangeServings.vue';
import AccountChangeName from '../components/AccountChangeName.vue';
import AccountChangeDeliveryNotes from '../components/AccountChangeDeliveryNotes.vue';
import AccountChangePhoneNumber from '../components/AccountChangePhoneNumber.vue';
import AccountChangeDeliveryDay from '../components/AccountChangeDeliveryDay.vue';

@Component({
  components: {
    AccountListItem,
    AccountChangeServings,
    AccountChangeName,
    AccountChangeDeliveryNotes,
    AccountChangePhoneNumber,
    AccountChangeDeliveryDay,
  },
})
export default class Account extends Vue {
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
  padding: 0 0 50px 0;
  align-content: center;
}

.footer-message-text {
  align-content: center;
  text-align: center;
  margin: 0;
  font-size: 16px;
}
</style>
