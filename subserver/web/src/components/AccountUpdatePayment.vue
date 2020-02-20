<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>
      <DialogConfirmPayment
        ref="dialog"
        Title="Update Payment"
        ButtonText="Edit"
        ConfirmText="Update"
        v-on:dialog-success="submit"
      >
        <template v-slot:dialog-content>
          <v-layout>
            <v-flex>
              <!-- https://francoislevesque.github.io/vue-braintree/configuration.html#enable-3d-secure -->
              <v-braintree
                :authorization="getAuthorization"
                @success="onSuccess"
                @error="onError"
              >
                <template v-slot:button="slotProps">
                  <v-btn
                    @click="slotProps.submit"
                    color="#E8554E"
                    depressed
                    class="white--text"
                    style="margin: 10px 0 -60px 0;"
                  >Update</v-btn>
                </template>
              </v-braintree>
            </v-flex>
          </v-layout>
        </template>
      </DialogConfirmPayment>
    </div>
    <p class="value">{{value}}</p>
    <hr class="divider-line">
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import DialogConfirmPayment from '../components/DialogConfirmPayment.vue';
import { IsError, ErrorAlert } from '../ts/errors';
import { UpdatePayment } from '../ts/service';
@Component({
  components: {
    DialogConfirmPayment,
  },
})
export default class AccountUpdatePayment extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  protected disableButton = false;

  get title(): string {
    return 'Payment Method';
  }

  get value(): string {
    return '';
  }

  get getAuthorization(): string {
    let authorization = 'production_tv5qygvt_wsgmypp8c46cnbpc';
    if (
      window.location.hostname === 'localhost' ||
      window.location.hostname === 'gigamunch-omninexus-dev.appspot.com'
    ) {
      authorization = 'sandbox_vprqjq87_4j6rdqcz74z7rt92';
    }
    return authorization;
  }

  public onSuccess(payload: any) {
    let nonce = payload.nonce;
    this.submit(nonce);
  }

  public onError(error: any) {
    let message = error.message;
    console.error(error);
    alert(message);
  }

  public req = {
    nonce: '',
  };

  protected submit(payment_method_nonce: string) {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      (this.$refs.dialog as DialogConfirmPayment).Dismiss();
      this.$emit('get-account-info');
    };
    if (payment_method_nonce == '') {
      alert('Payment is no selected.');
      return;
    }
    if (!this.sub) {
      alert('account info not loaded in payment section');
      return;
    }
    UpdatePayment(payment_method_nonce).then(handler);
  }
}
</script>

<style scoped lang="scss">
.list-item-container {
  margin: 16px 0;
}

.top-row {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  margin: 0;
}

.title {
  font-weight: 500;
  color: #333333;
}

.edit-button {
  cursor: pointer;
}

.value {
  color: #869995;
  font-size: 18px;
}

.divider-line {
  margin: 30px 10px 0 0;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}
</style>