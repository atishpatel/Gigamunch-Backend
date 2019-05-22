<template>
  <DialogConfirm
    ref="dialog"
    Title="Setup Activity"
    ButtonText="Setup Activity"
    ConfirmText="Setup"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <!-- TODO: -->
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { DeactivateSubscriber } from '../ts/service';
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
    reason: '',
  };

  protected submit() {
    if (!this.sub) {
      alert('sub not found');
      return;
    }
    DeactivateSubscriber(this.sub.id, this.req.reason).then((resp) => {
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
