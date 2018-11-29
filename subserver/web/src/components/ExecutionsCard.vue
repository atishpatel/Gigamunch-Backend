<template>
  <div>
    <router-link :to="to" tag="div" class="image unset-link">
      <img :src="src" :alt="computedImageAlt">
    </router-link>
    <div>
      <router-link :to="to" tag="h3" class="title unset-link" v-html="title"></router-link>
      <p class="description" v-html="description"></p>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ExecutionsItem extends Vue {
  @Prop()
  public title!: string;
  @Prop()
  public description!: string;
  @Prop()
  public src!: string;
  @Prop()
  public to!: object;

  get computedImageAlt(): string {
    return this.title + ' image';
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.image {
  display: block;
  width: 100%;
  position: relative;
  height: 0;
  padding: 56.25% 0 0 0;
  overflow: hidden;
  background-color: white;
  border-radius: 5px;
  img {
    position: absolute;
    display: block;
    width: 100%;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    margin: auto;
  }
}
.title {
  font-weight: 500;
  padding: 16px 0;
  margin: 0;
}
.description {
  padding: 0px;
  margin: 0px;
}
</style>