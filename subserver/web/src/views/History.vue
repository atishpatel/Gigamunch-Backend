<template>
  <div>
    <div class="content-container">
      <h1>Dinner History</h1>
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
import { Vue } from 'vue-property-decorator';
// import { GetAccountInfo } from '../ts/service';
// import { IsError } from '../ts/errors';

export default class Dinner extends Vue {
  protected accountInfo!: SubAPI.GetAccountInfoResp;

  public constructor() {
    super();
    // this.accountInfo = {
    //   address: {} as Common.Address,
    //   payment_info: {} as SubAPI.PaymentInfo,
    // } as SubAPI.GetAccountInfoResp;
  }

  // public created() {
  //   this.getAccountInfo();
  // }

  // public getAccountInfo() {}
}
</script>
<style scoped lang="scss">
.content-container {
  max-width: 700px;
  margin: auto;
  padding: 12px;
}
.divider-line {
  margin: 40px 0 40px 0;
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
