<template>
  <DialogConfirm
    ref="dialog"
    Title="If you don't mind, please let us know why you wish to cancel."
    ButtonText="Cancel Account"
    ConfirmText="Cancel Account"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-layout>
        <v-flex>
          <v-textarea
            class="field-right-padding"
            v-model="req.cancel_reason"
            outline
            round
          ></v-textarea>
        </v-flex>
      </v-layout>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import DialogConfirm from '../components/DialogConfirm.vue';
import { IsError, ErrorAlert } from '../ts/errors';
import { DeactivateSubscriber } from '../ts/service';
@Component({
  components: {
    DialogConfirm,
  },
})
export default class AccountChangeName extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;

  public req = {
    cancel_reason: '',
  };

  protected submit() {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }
      (this.$refs.dialog as DialogConfirm).Dismiss();
      this.$emit('get-account-info');
      this.$emit('get-ser-summary');
    };

    if (!this.sub) {
      alert('account info not loaded in cancel section');
      return;
    }
    DeactivateSubscriber(this.req.cancel_reason).then(handler);
  }
}
</script>

<style scoped lang="scss">
.list-item-container {
  margin: 16px 0;
}

.top-row {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  margin: 0;
}

.title {
  font-weight: 500;
  color: #333333;
}

.edit-button {
  cursor: pointer;
}

.value {
  color: #869995;
  font-size: 18px;
}

.divider-line {
  margin: 30px 10px 0 0;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}
</style>