<template>
  <div
    class="image"
    :style="{ backgroundImage: 'url(\'' + src + '\')' }"
    v-bind:class="{ rounded: rounded }"
  >
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Image169 extends Vue {
  @Prop()
  public src!: string;
  @Prop()
  public rounded!: boolean;
}
</script>

<style scoped lang="scss">
.image {
  height: 0;
  padding: 0; /* remove any pre-existing padding, just in case */
  background-position: center center;
  background-size: cover;
  padding-bottom: 56.25%;
  background-repeat: no-repeat;
}

.rounded {
  border-radius: 12px;
}
</style>