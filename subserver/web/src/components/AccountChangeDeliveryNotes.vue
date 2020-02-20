<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>

      <AccountDialogConfirm
        ref="dialog"
        Title="Update Delivery Notes"
        ButtonText="Edit"
        ConfirmText="Update"
        v-on:dialog-success="submit"
      >
        <template v-slot:dialog-content>
          <v-layout>
            <v-flex>
              <v-text-field
                class="field-right-padding"
                v-model="req.delivery_notes"
                label="Delivery Notes"
                outline
                round
              ></v-text-field>
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
export default class AccountChangeDeliveryNotes extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;

  get title(): string {
    return 'Delivery Notes';
  }

  get value(): string {
    if (this.sub) {
      if (!this.sub.delivery_notes || this.sub.delivery_notes === '') {
        return 'Not provided';
      } else {
        return this.sub.delivery_notes;
      }
    }
    return '';
  }

  public req = {
    delivery_notes: '',
  };

  protected submit() {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      (this.$refs.dialog as AccountDialogConfirm).Dismiss();
      this.$emit('get-account-info');
    };
    if (!this.sub) {
      alert('account info not loaded in deliery notes section');
      return;
    }
    if (!this.sub.email_prefs[0]) {
      alert('account info not loaded email_prefs in deliery notes section');
      return;
    }
    UpdateSubscriber(
      this.sub.email_prefs[0].first_name,
      this.sub.email_prefs[0].last_name,
      this.sub.address,
      this.req.delivery_notes,
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