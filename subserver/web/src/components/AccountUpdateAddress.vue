<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>
      <AccountDialogConfirm
        ref="dialog"
        Title="Update Address"
        ButtonText="Edit"
        ConfirmText="Update"
        v-on:dialog-success="submit"
      >
        <template v-slot:dialog-content>
          <v-layout>
            <v-flex>
              <vuetify-google-autocomplete
                ref="elAddress"
                id="map"
                append-icon="search"
                classname="form-control"
                placeholder="Select Address"
                v-on:placechanged="getAddressData"
                country="us"
                outlined
                outline
                aria-autocomplete="false"
                autocomplete="false"
              >

              </vuetify-google-autocomplete>
            </v-flex>
          </v-layout>
        </template>
      </AccountDialogConfirm>
    </div>
    <p class="value">{{value}}</p>
    <hr class="divider-line">
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import AccountDialogConfirm from '../components/AccountDialogConfirm.vue';
import { IsError, ErrorAlert } from '../ts/errors';
import { UpdateSubscriber } from '../ts/service';
@Component({
  components: {
    AccountDialogConfirm,
  },
})
export default class AccountUpdateAddress extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;

  get title(): string {
    return 'Address';
  }

  get value(): string {
    if (this.sub && this.sub.address) {
      const addr = this.sub.address;
      return `${addr.street}, ${addr.city}, ${addr.state}, ${addr.zip}, ${addr.country}`;
    } else {
      return '';
    }
  }

  public req = {
    address: {} as Common.Address,
  };

  protected getAddressData(addressData: any, placeResultData: any) {
    this.req.address.full_address = placeResultData.formatted_address;
  }

  protected submit() {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      (this.$refs.dialog as AccountDialogConfirm).Dismiss();
      this.$emit('get-account-info');
    };
    if (this.req.address.full_address === '') {
      alert('Address is no selected.');
      return;
    }
    if (!this.sub) {
      alert('account info not loaded in address section');
      return;
    }
    UpdateSubscriber(
      this.sub.email_prefs[0].first_name,
      this.sub.email_prefs[0].last_name,
      this.req.address,
      this.sub.delivery_notes,
      this.sub.phonenumbersString
    ).then(handler);
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
  margin: 30px 35px 0 0;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}
</style>