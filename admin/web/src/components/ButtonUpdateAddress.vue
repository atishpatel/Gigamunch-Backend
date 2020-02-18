<template>
  <DialogConfirm
    ref="dialog"
    Title="Update Address"
    ButtonText="Update Address"
    ConfirmText="Update"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-card-text>{{dialogText}}</v-card-text>
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
          id="map"
          append-icon="search"
          classname="form-control"
          placeholder="End address"
          v-on:placechanged="getAddressData"
          country="us"
          outlined
        >
        </vuetify-google-autocomplete>
      </v-flex>
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
  public dialogText = '';

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
    UpdateAddress(this.sub.id, this.address).then((resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }

      (this.$refs.dialog as DialogConfirm).Dismiss();
      this.$emit('dialog-success');
    });
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
