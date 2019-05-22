<template>
  <v-dialog
    v-model="dialog"
    max-width="500"
  >
    <template v-slot:activator="{ on }">
      <v-btn
        outline
        round
        v-on="on"
      >
        New Discount
      </v-btn>
    </template>
    <v-card class="dialog-card">
      <v-card-title>
        <span class="dialog-title">New Discount</span>
      </v-card-title>
      <v-card-text>
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
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          outline
          round
          @click="dialog = false"
        >Close</v-btn>
        <v-btn
          outline
          round
          @click="submit"
        >Create</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { DiscountSubscriber } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';

@Component({
  components: {},
})
export default class DialogNewDiscount extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public dialog = false;
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

        this.dialog = false;
        this.$emit('dialog-success');
      }
    );
  }
}
</script>

<style scoped lang="scss">
.dialog-card {
  padding: 12px;
}

.dialog-title {
  font-weight: 600;
  font-size: 24px;
}

.field-right-padding {
  padding-right: 12px;
}
</style>
