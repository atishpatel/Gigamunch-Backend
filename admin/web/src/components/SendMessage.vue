<template>
  <v-card
    outlined
    class="wrapper"
    elevation="0"
  >
    <v-textarea
      label="UserID, Emails, or Phone Numbers"
      outline
      v-model="ToField"
      @blur="updateEmails"
    ></v-textarea>
    <v-textarea
      label="Message"
      outline
      v-model="req.message"
      auto-grow
    ></v-textarea>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        outline
        round
        @click="submit"
      >Send Message</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { SendCustomerSMS } from '../ts/service';
import { IsError, ErrorAlert } from '../ts/errors';

@Component({
  components: {},
})
export default class SendMessage extends Vue {
  @Prop()
  public ToField!: string;
  @Prop({
    default: {
      emails: '',
      message: '',
    },
  })
  public req!: AdminAPI.SendCustomerSMSReq;

  protected submit() {
    this.updateEmails;
    if (this.req.emails.length <= 0) {
      alert('add emails');
      return;
    }
    if (this.req.message == '') {
      alert('enter a message');
      return;
    }

    SendCustomerSMS(this.req.emails, this.req.message).then((resp) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }

      this.$emit('send-message-success');
    });
  }

  protected updateEmails() {
    let to = this.ToField;
    to = to.replace(/\t/g, ',');
    to = to.replace(/\n/g, ',');
    to = to.replace(/    /g, ',');
    to = to.replace(/ /g, '');
    to = to.replace(/,,/g, ',');
    let toArray = to.split(',');
    this.ToField = to;
    this.req.emails = toArray;
    return toArray;
  }
}
</script>

<style scoped lang="scss">
.wrapper {
  padding: 12px;
}
</style>
