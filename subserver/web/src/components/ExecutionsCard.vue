<template>
  <router-link
    :to="to"
    tag="div"
    class="card-container"
  >
    <div class="dinner-image unset-link">
      <img
        :src="dinnerImageSource"
        :alt="computedDinnerImageAlt"
      >
    </div>
    <div class="cook-info-row">
      <div class="cook-face-image unset-link">
        <img
          :src="cookFaceImageSource"
          :alt="computedCookImageAlt"
        >
      </div>
      <div>
        <h3 class="cook-name unset-link">{{cookName}}</h3>
        <p class="subtitle">{{subtitle}}</p>
      </div>
    </div>

  </router-link>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ExecutionsItem extends Vue {
  @Prop()
  public cookName!: string;
  @Prop()
  public dinnerImageSource!: string;
  @Prop()
  public cookFaceImageSource!: string;
  @Prop()
  public nationality!: string;
  @Prop()
  public to!: object;

  get computedCookImageAlt(): string {
    return this.cookName + ' image';
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
.card-container {
  display: flex;
  flex-direction: column;
  cursor: pointer;
}
.cook-info-row {
  flex: none;
  display: flex;
  flex-direction: row;
  margin: 15px 0 0 0;
}
.dinner-image {
  flex: none;
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
  height: 85px;
  width: 85px;
  position: relative;
  margin: 0 15px 0 0;
  img {
    position: absolute;
    display: block;
    width: 100%;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    margin: auto;
    border-radius: 50%;
  }
}
.cook-name {
  font-weight: 500;
  font-size: 17px;
  padding: 12px 0 6px 0;
  margin: 0;
}
.subtitle {
  padding: 0px;
  margin: 0px;
  color: #869995;
}
@media (min-width: 600px) {
  .cook-name {
    font-size: 22px;
  }
  .subtitle {
    font-size: 18px;
  }
}
</style>
