<template>
  <router-link
    :to="to"
    tag="div"
  >
    <div class="image unset-link">
      <img
        :src="dinner_image_src"
        :alt="computedDinnerImageAlt"
      >
    </div>
    <div class="cook-info-row">
      <div class="cook-face-image unset-link">
        <img
          :src="cook_face_image_src"
          :alt="computedCookImageAlt"
        >
      </div>
      <div>
        <h3
          class="cook_name unset-link"
          v-html="cook_name"
        ></h3>
        <p
          class="subtitle"
          v-html="subtitle"
        ></p>
      </div>
    </div>

  </router-link>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ExecutionsItem extends Vue {
  @Prop()
  public cook_name!: string;
  @Prop()
  public dinner_image_src!: string;
  @Prop()
  public cook_face_image_src!: string;
  @Prop()
  public nationality!: string;
  @Prop()
  public to!: object;

  get computedCookImageAlt(): string {
    return this.cook_name + ' image';
  }

  get computedDinnerImageAlt(): string {
    return this.nationality + 'dinner image';
  }

  get subtitle(): string {
    return 'your ' + this.nationality + ' host';
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.cook-info-row {
  display: flex;
  flex-direction: row;
  margin: 15px 0 0 0;
}
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
.cook-face-image {
  height: 100px;
  width: 100px;
  margin: 0 15px 0 0;
  img {
    width: 100%;
    margin: 0;
    border-radius: 50%;
  }
}
.cook_name {
  font-weight: 500;
  padding: 16px 0;
  margin: 0;
}
.subtitle {
  padding: 0px;
  margin: 0px;
}
</style>
