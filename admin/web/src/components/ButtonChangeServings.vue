<template>
  <DialogConfirm
    ref="dialog"
    :Title="computedText"
    :ButtonText="buttonText"
    :ConfirmText="confirmText"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <!-- TODO: -->
      <v-card-text>Activate <span class="bold"></span>?</v-card-text>
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
export default class ButtonChangeServings extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  @Prop({ default: false })
  public changePermanently!: boolean;
  public req = {
    date: '',
  };
  public buttonText = '';
  public confirmText = '';

  get computedText() {
    if (this.changePermanently) {
      this.buttonText = 'Change Servings Permanently';
      this.confirmText = 'Change';
      return 'Change Servings Permanently';
    }
    this.buttonText = 'Change Servings For Date';
    this.confirmText = 'Change';
    return 'Change Servings for Date';
  }

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
