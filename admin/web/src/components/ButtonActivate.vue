<template>
  <DialogConfirm
    ref="dialog"
    Title="Activate Subscriber"
    ButtonText="Activate"
    ConfirmText="Activate"
    :ButtonDisabled="sub.active"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <!-- TODO: -->
      <v-card-text>Activate <span class="bold">{{sub.namesString}}</span>?</v-card-text>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { ActivateSubscriber } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class ButtonActivate extends Vue {
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
    ActivateSubscriber(this.sub.id, this.req.date).then((resp) => {
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
