<template>
  <AccountDialogConfirm
    v-if="changePermanently"
    ref="dialog"
    :Title="computedText"
    :ButtonText="buttonText"
    :ConfirmText="confirmText"
    v-on:dialog-success="submit"
    :ButtonDisabled="ButtonDisabled"
  >
    <template v-slot:dialog-content>
      <v-card-text>{{dialogText}}</v-card-text>
      <v-layout>
        <v-flex>
          <v-select
            :items="servingSizes"
            class="field-right-padding"
            v-model="req.servings_non_veg"
            label="Meat Servings"
            :placeholder="nonvegPlaceholder"
            outline
            round
          ></v-select>
        </v-flex>
        <v-flex>
          <v-select
            :items="servingSizes"
            v-model="req.servings_veg"
            label="Veg Servings"
            :placeholder="vegPlaceholder"
            outline
            round
          ></v-select>
        </v-flex>
      </v-layout>
    </template>
  </AccountDialogConfirm>

  <DialogConfirm
    v-else
    ref="dialog"
    :Title="computedText"
    :ButtonText="buttonText"
    :ConfirmText="confirmText"
    v-on:dialog-success="submit"
    :ButtonDisabled="ButtonDisabled"
  >
    <template v-slot:dialog-content>
      <v-card-text>{{dialogText}}</v-card-text>
      <v-layout>
        <v-flex>
          <v-select
            :items="servingSizes"
            class="field-right-padding"
            v-model="req.servings_non_veg"
            label="Meat Servings"
            :placeholder="nonvegPlaceholder"
            outline
            round
          ></v-select>
        </v-flex>
        <v-flex>
          <v-select
            :items="servingSizes"
            v-model="req.servings_veg"
            label="Veg Servings"
            :placeholder="vegPlaceholder"
            outline
            round
          ></v-select>
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
import AccountDialogConfirm from './AccountDialogConfirm.vue';

@Component({
  components: {
    DialogConfirm,
    AccountDialogConfirm,
  },
})
export default class ButtonChangeServings extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  @Prop()
  public activity!: Types.ActivityExtended;
  @Prop({ default: false })
  public changePermanently!: boolean;
  @Prop()
  public ButtonDisabled!: boolean;
  public req = {
    servings_non_veg: '',
    servings_veg: '',
  };
  public buttonText = '';
  public confirmText = '';
  public dialogText = '';
  public servingSizes = [0, 2, 4, 6, 8, 10];

  get computedText() {
    if (this.changePermanently) {
      this.buttonText = 'Edit';
      this.confirmText = 'Change';
      if (this.sub) {
        this.dialogText = `Right now, you will receive ${this.sub.servings_non_vegetarian} meat servings and ${this.sub.servings_vegetarian} vegetarian servings by default.`;
      } else {
        this.dialogText = '';
      }
      return 'Change defalt serving size';
    }
    this.buttonText = 'Change Servings';
    this.confirmText = 'Update';
    if (this.activity) {
      this.dialogText = `Right now, you will receive ${this.activity.servings_non_vegetarian} meat servings and ${this.activity.servings_vegetarian} vegetarian servings for this day.`;
    } else {
      this.dialogText = '';
    }

    return 'Select servings for this day';
  }

  get vegPlaceholder(): string {
    if (this.changePermanently) {
      if (this.sub) {
        return `${this.sub.servings_vegetarian}`;
      } else {
        return '0';
      }
    } else {
      if (this.activity) {
        return `${this.activity.servings_vegetarian}`;
      } else {
        return '0';
      }
    }
  }

  get nonvegPlaceholder(): string {
    if (this.changePermanently) {
      if (this.sub) {
        return `${this.sub.servings_non_vegetarian}`;
      } else {
        return '0';
      }
    } else {
      if (this.activity) {
        return `${this.activity.servings_non_vegetarian}`;
      } else {
        return '0';
      }
    }
  }

  get vegServings(): string {
    if (this.activity) {
      return `${this.activity.servings_vegetarian}`;
    }
    return '';
  }

  get nonvegServings(): string {
    if (this.activity) {
      return `${this.activity.servings_non_vegetarian}`;
    }
    return '';
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
    if (servings_non_veg === 0 && servings_veg === 0) {
      if (this.changePermanently) {
        alert(
          'Try selecting more than zero servings. Or press the cancel button at the bottom of the page if you would like to stop receiving dinners.'
        );
      } else {
        alert(
          'Try selecting more than zero servings. Or press the skip button if you would like to skip this dinner.'
        );
      }
      return;
    }
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
