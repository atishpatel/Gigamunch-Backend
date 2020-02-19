<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>
      <ButtonChangeServings
        changePermanently=true
        v-on:dialog-success="updatedServings"
        :sub="sub"
        text
        large
        color="#E8554E"
        flat
        class="edit-button"
        :ripple=false
      ></ButtonChangeServings>
    </div>
    <p class="value">{{value}}</p>
    <hr class="divider-line">
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import ButtonChangeServings from '../components/ButtonChangeServings.vue';

@Component({
  components: {
    ButtonChangeServings,
  },
})
export default class AccountChangeServings extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  @Prop()
  public activity!: Types.ActivityExtended;
  @Prop()
  protected accountInfo!: SubAPI.GetAccountInfoResp;

  get title(): string {
    return 'Default Serving Size';
  }

  get value(): string {
    if (this.sub) {
      return `${this.sub.servings_non_vegetarian} meat servings and ${this.sub.servings_vegetarian} vegetarian servings`;
    } else {
      return '';
    }
  }

  protected updatedServings() {
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