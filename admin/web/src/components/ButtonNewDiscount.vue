<template>
  <DialogConfirm
    ref="dialog"
    Title="New Discount"
    ButtonText="New Discount"
    ConfirmText="Create"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-layout>
        <v-flex>
          <v-text-field
            class="field-right-padding"
            v-model="req.amount"
            label="Discount Amount"
            prefix="$"
            type="number"
            outline
            round
          ></v-text-field>
        </v-flex>
        <v-flex>
          <v-text-field
            v-model="req.percent"
            label="Discount Percent"
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
import { DiscountSubscriber } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class DialogNewDiscount extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public req = {
    amount: 0.0,
    percent: 0,
  };

  protected submit() {
    if (!this.sub) {
      alert('sub not found');
      return;
    }
    DiscountSubscriber(this.sub.id, this.req.amount, this.req.percent).then(
      (resp) => {
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
.field-right-padding {
  padding-right: 12px;
}
</style>
