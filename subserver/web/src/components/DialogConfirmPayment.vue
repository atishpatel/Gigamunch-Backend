<template>
  <v-dialog
    v-model="dialog"
    max-width="500"
  >
    <template v-slot:activator="{ on }">
      <v-btn
        text
        large
        color="#E8554E"
        flat
        class="edit-button"
        :ripple=false
        v-on="on"
        :disabled="ButtonDisabled"
      >
        {{ButtonText}}
      </v-btn>
    </template>
    <v-card class="dialog-card">
      <v-card-title>
        <span class="dialog-title">{{Title}}</span>
      </v-card-title>
      <slot name="dialog-content">
      </slot>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          depressed
          color="#E8554E"
          class="white--text"
          @click="dialog = false"
        >{{CancelText}}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component({
  components: {},
})
export default class DialogConfirm extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  @Prop({ default: false })
  public ButtonDisabled!: boolean;
  @Prop()
  public ButtonText!: string;
  @Prop({ default: 'Submit' })
  public ConfirmText!: string;
  @Prop({ default: 'Close' })
  public CancelText!: string;
  @Prop()
  public Title!: string;
  public dialog = false;

  public Dismiss() {
    this.dialog = false;
  }
}
</script>

<style scoped lang="scss">
.dialog-card {
  padding: 12px;
}

.dialog-title {
  font-weight: 600;
  font-size: 24px;
}

.field-right-padding {
  padding-right: 12px;
}
</style>
