<template>
  <DialogConfirm
    ref="dialog"
    Title="Deactivate Subscriber"
    ButtonText="Deactivate"
    ConfirmText="Deactivate"
    :ButtonDisabled="!sub.active"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <!-- TODO:  -->
      <v-card-text>Deactivate <span class="bold">{{sub.namesString}}</span>?</v-card-text>
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
export default class ButtonDeactivate extends Vue {
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
