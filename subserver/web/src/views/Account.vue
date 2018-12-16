<template>
  <div>
    <LoadingView hide="loading"></LoadingView>
    <h1>Account</h1>
    <section>
      <div class="field">
        <div class="field-title">Name</div>
        <div class="field-value">
          {{name}}
        </div>
        <div class="field-action">
          Change
        </div>
      </div>
    </section>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { GetAccountInfo } from '../ts/service';
import { IsError } from '../ts/errors';
import LoadingView from './subviews/LoadingView.vue';

@Component({
  components: { LoadingView },
})
export default class Dinner extends Vue {
  protected accountInfo!: SubAPI.GetAccountInfoResp;
  protected name!: string;
  protected loading!: boolean;

  public constructor() {
    super();
    this.accountInfo = {
      address: {} as Common.Address,
      payment_info: {} as SubAPI.PaymentInfo,
    } as SubAPI.GetAccountInfoResp;
  }

  public created() {
    this.getAccountInfo();
  }

  public getAccountInfo() {
    this.loading = true;
    GetAccountInfo().then((resp) => {
      this.loading = false;
      if (IsError(resp)) {
        // TODO: handle errors
        return;
      }
      this.accountInfo = resp;
      this.name =
        resp.email_prefs[0].first_name + ' ' + resp.email_prefs[0].last_name;
    });
  }
}
</script>
<style scoped lang="scss">
.field {
  display: flex;
  flex-direction: row;
}

.field-title {
  font-weight: 600;
}
</style>
