<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>

      <DialogConfirm
        ref="dialog"
        Title="Update Phone Number"
        ButtonText="Edit"
        ConfirmText="Update"
        v-on:dialog-success="submit"
      >
        <template v-slot:dialog-content>
          <v-layout>
            <v-flex>
              <v-text-field
                class="field-right-padding"
                v-model="req.phonenumbersString"
                label="Phone Number"
                outline
                placeholder="555-555-5555"
                maxlength="12"
                round
              ></v-text-field>
            </v-flex>
          </v-layout>
        </template>
      </DialogConfirm>
    </div>
    <p class="value">{{value}}</p>
    <hr class="divider-line">
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import DialogConfirm from '../components/DialogConfirm.vue';
import { IsError, ErrorAlert } from '../ts/errors';
import { UpdateSubscriber } from '../ts/service';
@Component({
  components: {
    DialogConfirm,
  },
})
export default class AccountChangePhoneNumber extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;

  get title(): string {
    return 'Phone Number';
  }

  get value(): string {
    if (this.sub) {
      if (
        !this.sub.phone_prefs ||
        !this.sub.phone_prefs[0] ||
        !this.sub.phone_prefs[0].number ||
        this.sub.phone_prefs[0].number === ''
      ) {
        return 'Not provided - You will miss out on delivery texts';
      } else {
        return this.sub.phone_prefs[0].number;
      }
    }
    return '';
  }

  public req = {
    phonenumbersString: '',
  };

  protected submit() {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      (this.$refs.dialog as DialogConfirm).Dismiss();
      this.$emit('get-account-info');
    };
    if (!this.sub) {
      alert('account info not loaded in phone number section');
      return;
    }
    if (!this.sub.email_prefs[0]) {
      alert('account info not loaded email_prefs in phone number section');
      return;
    }
    UpdateSubscriber(
      this.sub.email_prefs[0].first_name,
      this.sub.email_prefs[0].last_name,
      this.sub.address,
      this.sub.delivery_notes,
      this.req.phonenumbersString
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
  margin: 30px 10px 0 0;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}
</style>