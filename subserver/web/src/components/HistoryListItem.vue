<template>
  <div class="list-item-container">
    <div class="top-row">
      <p class="title">{{title}}</p>
      <v-spacer></v-spacer>
      <v-btn
        text
        large
        color="#E8554E"
        flat
        class="edit-button"
        :ripple=false
      >Edit</v-btn>
    </div>
    <p class="value">{{value}}</p>
    <hr class="divider-line">
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class AccountListItem extends Vue {
  @Prop()
  public title!: string;
  @Prop()
  public value!: string;
  @Prop()
  public activity!: Common.Activity;
  @Prop()
  public culture!: string;
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
  margin: 30px 40px 0 0;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}
</style>