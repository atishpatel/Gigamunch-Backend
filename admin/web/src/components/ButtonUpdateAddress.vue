<template>
  <DialogConfirm
    ref="dialog"
    Title="Update Address And Notes"
    ButtonText="Update Address And Notes"
    ConfirmText="Update"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-layout>
        <v-flex>
          <v-text-field
            class="field-right-padding"
            v-model="address.apt"
            label="Apt"
            type="text"
            outline
            round
          ></v-text-field>
        </v-flex>
      </v-layout>
      <v-flex>
        <vuetify-google-autocomplete
          ref="elAddress"
          id="map"
          append-icon="search"
          classname="form-control"
          placeholder="End address"
          v-on:placechanged="getAddressData"
          country="us"
          outlined
          aria-autocomplete="false"
          autocomplete="false"
        >
        </vuetify-google-autocomplete>
      </v-flex>
      <v-textarea
        class="field-right-padding"
        v-model="deliveryNotes"
        label="Delivery Notes"
        outline
        round
        auto-grow
      ></v-textarea>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { UpdateAddress } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class ButtonChangeServings extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public address: Common.Address = {
    apt: '',
    full_address: '',
  } as Common.Address;
  public deliveryNotes = '';

  public setAddressAndNotes(addr: Common.Address, dNotes: string) {
    if (!addr.full_address) {
      addr.full_address = `${addr.street}, ${addr.city}, ${addr.state}, ${addr.zip}, ${addr.country}`;
    }
    if (this.$refs.elAddress) {
      // @ts-ignore
      this.$refs.elAddress.update(addr.full_address);
    }
    this.address = addr;
    this.deliveryNotes = dNotes;
  }

  protected getAddressData(addressData: any, placeResultData: any) {
    console.log('placeResultData', placeResultData);
    this.address.full_address = placeResultData.formatted_address;
  }

  protected submit() {
    if (!this.address.full_address) {
      alert('Please enter an address');
      return;
    }
    this.address.latitude = 0;
    this.address.longitude = 0;
    UpdateAddress(this.sub.id, this.address, this.deliveryNotes).then(
      (resp: any) => {
        if (IsError(resp)) {
          ErrorAlert(resp);
          return;
        }

        (this.$refs.dialog as DialogConfirm).Dismiss();
        this.$emit('dialog-success');
      }
    );
  }
}
</script>

<style scoped lang="scss">
.bold {
  font-weight: 500;
}

.field-right-padding {
  padding-right: 12px;
}
</style>
