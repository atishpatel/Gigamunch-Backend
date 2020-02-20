<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>
      <DialogConfirm
        ref="dialog"
        Title="Change Default Delivery Day"
        ButtonText="Edit"
        ConfirmText="Update"
        v-on:dialog-success="submit"
      >
        <template v-slot:dialog-content>
          <v-layout>
            <v-flex>
              <v-select
                :items="availableDeliveryDays"
                class="field-right-padding"
                v-model="req.day"
                label="Delivery Day"
                outline
                round
              ></v-select>
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
import { ChangePlanDay } from '../ts/service';
@Component({
  components: {
    DialogConfirm,
  },
})
export default class AccountChangeDeliveryDay extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;

  get title(): string {
    return 'Default Delivery Day';
  }

  get value(): string {
    return this.sub.plan_weekday;
  }

  //   public availableDeliveryDays = ['Monday', 'Thursday'];
  public availableDeliveryDays = ['Monday'];

  public req = {
    day: '',
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
      alert('account info not loaded in delivery day section');
      return;
    }
    ChangePlanDay(this.req.day).then(handler);
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