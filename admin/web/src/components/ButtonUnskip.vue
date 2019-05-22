<template>
  <DialogConfirm
    ref="dialog"
    Title="Unskip Activity"
    ButtonText="Unskip"
    ConfirmText="Unskip"
    :ButtonDisabled="!activity.skip"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-card-text>Unskip <span class="bold">{{activity.first_name}} {{activity.last_name}}</span> for <span class="bold">{{activity.dateFull}}</span>?</v-card-text>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { UnskipActivity } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class ButtonUnskip extends Vue {
  @Prop()
  public activity!: Types.ActivityExtended;
  public req = {};

  protected submit() {
    if (!this.activity) {
      alert('activity not found');
      return;
    }
    UnskipActivity(this.activity.user_id, this.activity.date).then((resp) => {
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
