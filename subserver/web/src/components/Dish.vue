<template>
  <!-- <div
    class="dish-name"
    :style="{ backgroundImage: 'url(\'' + src + '\')' }"
    v-bind:class="{ rounded: rounded }"
  >
  </div> -->
  <div class="dish-container">
    <p class="description"><strong>{{name}}</strong> – {{description}}</p>
    <p class="ingredients">{{ingredients}}</p>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Dish extends Vue {
  @Prop()
  public name!: string;
  @Prop()
  public description!: string;
  @Prop()
  public ingredients!: string;
}
</script>

<style scoped lang="scss">
.dish-container {
}
.name {
}

.description {
  margin: 30px 0 5px 0;
}

.ingredients {
  margin: 0;
  color: #869995;
}
</style>