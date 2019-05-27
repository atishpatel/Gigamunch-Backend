<template>
  <DialogConfirm
    ref="dialog"
    :Title="computedText"
    :ButtonText="buttonText"
    :ConfirmText="confirmText"
    v-on:dialog-success="submit"
  >
    <template v-slot:dialog-content>
      <v-card-text>{{dialogText}}</v-card-text>
      <v-layout>
        <v-flex>
          <v-text-field
            class="field-right-padding"
            v-model="req.servings_non_veg"
            label="Non-veg Servings"
            type="number"
            outline
            round
          ></v-text-field>
        </v-flex>
        <v-flex>
          <v-text-field
            v-model="req.servings_veg"
            label="Veg Servings"
            type="number"
            outline
            round
          ></v-text-field>
        </v-flex>
      </v-layout>
    </template>
  </DialogConfirm>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import {
  ChangeActivityServings,
  ChangeSubscriberServings,
} from '../ts/service';
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
  @Prop()
  public activity!: Types.ActivityExtended;
  @Prop({ default: false })
  public changePermanently!: boolean;
  public req = {
    servings_non_veg: '',
    servings_veg: '',
  };
  public buttonText = '';
  public confirmText = '';
  public dialogText = '';

  get computedText() {
    if (this.changePermanently) {
      this.buttonText = 'Change Servings Permanently';
      this.confirmText = 'Change';
      this.dialogText = `Subscriber currently has ${
        this.sub.servings_non_vegetarian
      } non-veg and ${this.sub.servings_vegetarian} veg servings.`;
      return 'Change Servings Permanently';
    }
    this.buttonText = 'Change Servings For Date';
    this.confirmText = 'Change';
    this.dialogText = `Activity currently has ${
      this.activity.servings_non_vegetarian
    } non-veg and ${this.activity.servings_vegetarian} veg servings.`;

    return 'Change Servings for Date';
  }

  protected submit() {
    const handler = (resp: any) => {
      if (IsError(resp)) {
        ErrorAlert(resp);
        return;
      }

      (this.$refs.dialog as DialogConfirm).Dismiss();
      this.$emit('dialog-success');
    };
    const servings_non_veg = Number(this.req.servings_non_veg);
    const servings_veg = Number(this.req.servings_veg);
    if (this.changePermanently) {
      if (!this.sub) {
        alert('sub not found');
        return;
      }
      ChangeSubscriberServings(
        this.sub.id,
        servings_non_veg,
        servings_veg
      ).then(handler);
    } else {
      if (!this.activity) {
        alert('activity not found');
        return;
      }
      ChangeActivityServings(
        this.activity.user_id,
        servings_non_veg,
        servings_veg,
        this.activity.date
      ).then(handler);
    }
  }
}
</script>

<style scoped lang="scss">
.bold {
  font-weight: 500;
}

.field-right-padding {
  padding-right: 12px;
}
</style>
