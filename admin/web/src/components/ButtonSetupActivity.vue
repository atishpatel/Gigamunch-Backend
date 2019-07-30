<template>
  <DialogConfirm
    ref="dialog"
    Title="Setup Activity"
    ButtonText="Setup Activity"
    ConfirmText="Setup"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-date-picker
        landscape
        v-model="req.date"
        :allowed-dates="allowedDates"
      ></v-date-picker>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { SetupActivity } from '../ts/service';
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
  };

  protected submit() {
    if (!this.sub) {
      alert('sub not found');
      return;
    }
    const d = new Date(this.req.date);
    SetupActivity(this.sub.id, d.toISOString()).then((resp) => {
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
