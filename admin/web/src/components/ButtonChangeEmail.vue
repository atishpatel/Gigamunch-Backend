<template>
  <DialogConfirm
    ref="dialog"
    Title="Change Email"
    ButtonText="Change Email"
    ConfirmText="Change Email"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-flex>
        <v-text-field
          class="field-right-padding"
          v-model="req.new_email"
          label="New Email"
          outline
          round
        ></v-text-field>
      </v-flex>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { ReplaceSubscriberEmail } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';
import DialogConfirm from './DialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
  },
})
export default class ButtonChangeEmail extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public dialog = false;
  public req = {
    new_email: '',
  };

  protected submit() {
    if (!this.sub) {
      alert('sub not found');
      return;
    }
    if (!this.req.new_email) {
      alert('new email cannot be empty');
      return;
    }
    const old_email = this.sub.emails[0];
    ReplaceSubscriberEmail(old_email, this.req.new_email).then((resp) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      window.location.href = window.location.href.replace(
        old_email,
        this.sub.id
      );
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
