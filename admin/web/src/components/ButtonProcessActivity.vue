<template>
  <DialogConfirm
    ref="dialog"
    Title="Process Activity"
    ButtonText="Process"
    ConfirmText="Process"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { ProcessActivity } from '../ts/service';
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
  public req = {};

  protected submit() {
    if (!this.activity) {
      alert('activity not found');
      return;
    }
    ProcessActivity(this.activity.user_id, this.activity.date).then((resp) => {
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
