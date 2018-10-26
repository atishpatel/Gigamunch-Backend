<template>
    <div>
        <div class="image">
            <img :src="src" :alt="computedImageAlt">
        </div>
        <h3 class="title" v-html="title"></h3>
        <p class="description" v-html="description"></p>
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