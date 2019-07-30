<template>
  <DialogConfirm
    ref="dialog"
    Title="Change Plan Day"
    ButtonText="Change Plan Day"
    ConfirmText="Change"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-select
        label="Plan Day"
        :items="planDayOptions"
        v-model="req.new_plan_day"
      >
      </v-select>
      <v-text>Date after which activities switch to new plan day</v-text>
      <v-date-picker
        landscape
        v-model="req.date"
        outline
      ></v-date-picker>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { ChangeSubscriberPlanDay } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class ButtonSetupActivity extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public req = {
    date: '',
    new_plan_day: '',
  };
  public planDayOptions = ['Monday', 'Thursday'];

  protected submit() {
    if (!this.sub) {
      alert('sub not found');
      return;
    }
    if (!this.req.new_plan_day) {
      alert('select a new plan day');
      return;
    }
    const d = new Date(this.req.date);
    ChangeSubscriberPlanDay(
      this.sub.id,
      this.req.new_plan_day,
      d.toISOString()
    ).then((resp) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }

      (this.$refs.dialog as DialogConfirm).Dismiss();
      this.$emit('dialog-success');
    });
  }

  protected allowedDates(v: string) {
    const d = new Date(v);
    const day = d.getDay();
    if (day === 0 || day === 3) {
      // Monday or Thursday
      return true;
    }
    return false;
  }
}
</script>

<style scoped lang="scss">
.bold {
  font-weight: 500;
}
</style>
