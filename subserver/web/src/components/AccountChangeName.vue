<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>
      <!-- <ButtonChangeServings
        changePermanently=true
        v-on:dialog-success="updatedServings"
        :sub="sub"
        text
        large
        color="#E8554E"
        flat
        class="edit-button"
        :ripple=false
      ></ButtonChangeServings> -->
      <DialogConfirm
        ref="dialog"
        Title="Change Name"
        ButtonText="Edit"
        ConfirmText="Update"
        v-on:dialog-success="submit"
      >
        <template v-slot:dialog-content>
          <v-layout>
            <v-flex>
              <v-text-field
                class="field-right-padding"
                v-model="req.first_name"
                label="First Name"
                outline
                round
              ></v-text-field>
            </v-flex>
            <v-flex>
              <v-text-field
                v-model="req.last_name"
                label="Last Name"
                outline
                round
              ></v-text-field>
            </v-flex>
          </v-layout>
        </template>
      </DialogConfirm>
    </div>
    <p class="value">{{value}}</p>
    <hr class="divider-line">
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import DialogConfirm from '../components/DialogConfirm.vue';
import { IsError, ErrorAlert } from '../ts/errors';
@Component({
  components: {
    DialogConfirm,
  },
})
export default class AccountChangeName extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;

  get title(): string {
    return 'Name';
  }

  get value(): string {
    if (this.sub && this.sub.email_prefs) {
      return `${this.sub.email_prefs[0].first_name} ${this.sub.email_prefs[0].last_name}`;
    } else {
      return '';
    }
  }

  public req = {
    first_name: '',
    last_name: '',
  };

  protected submit() {
    // const handler = (resp: any) => {
    //   if (IsError(resp)) {
    //     ErrorAlert(resp);
    //     return;
    //   }
    // if this.req.first_name == ''
    //   (this.$refs.dialog as DialogConfirm).Dismiss();
    //   this.$emit('dialog-success');
    // };
    this.$emit('get-account-info');
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