<template>
  <DialogConfirm
    ref="dialog"
    Title="Refund And Skip"
    ButtonText="Refund And Skip"
    ConfirmText="Refund And Skip"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-card-text>Refund and skip <span class="bold">{{activity.first_name}} {{activity.last_name}}</span> for <span class="bold">{{activity.dateFull}}</span>?</v-card-text>
      <v-layout>
        <v-flex>
          <v-text-field
            class="field-right-padding"
            v-model="req.amount"
            label="Amount"
            prefix="$"
            type="number"
            outline
            round
          ></v-text-field>
        </v-flex>
        <v-flex>
          <v-text-field
            v-model="req.percent"
            label="Percent"
            suffix="%"
            type="number"
            max="100"
            outline
            round
          ></v-text-field>
        </v-flex>
      </v-layout>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { RefundAndSkipActivity } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class ButtonSkip extends Vue {
  @Prop()
  public activity!: Types.ActivityExtended;
  public dialog = false;
  public req = {
    amount: 0.0,
    percent: 0,
  };

  protected submit() {
    if (!this.activity) {
      alert('activity not found');
      return;
    }
    RefundAndSkipActivity(
      this.activity.user_id,
      this.activity.date,
      this.req.amount,
      this.req.percent
    ).then((resp) => {
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
</style>
