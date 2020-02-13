<template>
  <div class="view">
    <div class="hero-image">
      <div class="hero-image-gradient">
        <Image169 :src="exe.content.landscape_image_url"></Image169>
      </div>
      <div
        class="hero-image-text"
        v-html="heroImageText"
      ></div>
    </div>
    <div class="content-container">
      <div class="host-action">
        <div class="host">
          <div class="host-image">
            <div
              class="host-image-image"
              :style="{ backgroundImage: 'url(\'' + exe.email.cook_face_image_url + '\')' }"
            ></div>
          </div>
          <div class="host-text">
            <h2 class="host-text-name">{{cultureCookName}}</h2>
            <p class="host-text-hosted-by">{{hostSubtitle}}</p>
          </div>
        </div>
        <!-- <div class="action">
        <div v-if="userSummary.is_logged_in === true">signed in view</div>
        <div v-else>singed out</div>
        <p>action</p>
      </div> -->
      </div>
      <div class="culture-description">
        <p class="culture-description-text">{{cultureDescription}}</p>
      </div>
      <div class="section-title">
        <h2 class="dinner-image-title-text">{{dinnerImageTitle}}</h2>
      </div>
      <Image169
        :src="dinnerImageSrc"
        :rounded=true
      ></Image169>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Image169 from '../../components/Image169.vue';

@Component({
  components: {
    Image169,
  },
})
export default class DinnerPublished extends Vue {
  @Prop()
  public exe!: Common.Execution;
  @Prop()
  public activity!: Common.Activity;
  @Prop()
  public userSummary!: SubAPI.GetUserSummaryResp;

  get heroImageText(): string {
    return 'Your Jorney to ' + this.exe.culture.country;
  }

  get hostSubtitle(): string {
    return 'your ' + this.exe.culture.nationality + ' host';
  }

  get cultureCookName(): string {
    return (
      this.exe.culture_cook.first_name + ' ' + this.exe.culture_cook.last_name
    );
  }

  get cultureDescription(): string {
    return this.exe.culture.description;
  }

  get dinnerImageTitle(): string {
    return this.exe.culture_cook.first_name + "'s Dinner";
  }

  get dinnerImageSrc(): string {
    return this.exe.content.hands_plate_non_veg_image_url;
  }
}
</script>
<style scoped lang="scss">
.view {
  max-width: 850px;
  margin: auto;
  box-shadow: 0 0 15px grey;
}

$view-edge-padding-desktop: 120px;
$view-edge-padding-mobile: 24px;
// hero image
.hero-image {
  position: relative;
}
.hero-image-img {
  width: 100%;
}
.hero-image-gradient::before {
  display: block;
  position: absolute;
  top: 0;
  background-image: linear-gradient(
    rgba(100, 100, 100, 0.69),
    rgba(0, 0, 0, 0)
  );
  height: 15vw;
  width: 100%;
  content: '';
}
.hero-image-text {
  position: absolute;
  color: white;
  font-size: 28px;
  font-weight: 700;
  top: 24px;
  left: 28px;
}
@media (min-width: 800px) {
  .hero-image-text {
    font-size: 40px;
  }
}

.content-container {
  padding: 0 $view-edge-padding-mobile;
}
@media (min-width: 800px) {
  .content-container {
    padding: 0 $view-edge-padding-desktop;
  }
}

// host-action
.host-action {
}
.host {
  display: flex;
  position: relative;
  top: -35px;
}
.host-image {
  background-color: white;
  border-radius: 50%;
  padding: 4px;
}
.host-image-image {
  margin: auto;
  width: 100px;
  height: 100px;
  background-size: cover;
  border-radius: 50%;
}
.host-text {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding-left: 12px;
}
.host-text-hosted-by {
  margin: 0 0 10px 0;
  opacity: 0.75;
}
.host-text-name {
  margin: 6px 0 0 0;
}
.section-title {
  margin: 36px 0 18px 0;
}
</style>
