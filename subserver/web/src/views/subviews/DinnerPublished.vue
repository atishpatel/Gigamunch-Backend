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
    <div class="host-action">
      <div class="host">
        <div class="host-image">
          <div
            class="host-image-image"
            :style="{ backgroundImage: 'url(\'' + exe.email.cook_face_image_url + '\')' }"
          ></div>
        </div>
        <div class="host-text">
          <p class="host-text-hosted-by">hosted by</p>
          <h2 class="host-text-name">{{cultureCookName}}</h2>
        </div>
      </div>
      <!-- <div class="action">
        <div v-if="userSummary.is_logged_in === true">signed in view</div>
        <div v-else>singed out</div>
        <p>action</p>
      </div> -->
    </div>
    <div class="footer-message">
      <p class="footer-message-text">Feel free to talk to us at</p>
      <p class="footer-message-text"><a href="mailto:hello@eatgigamunch.com"><strong>hello@eatgigamunch.com</strong></a></p>
      <p
        class="footer-message-text"
        style="margin-top: 12px;"
      ><strong>We're here for you.</strong></p>
      <p
        class="footer-message-text"
        style="margin-top: 32px;"
      >💛&nbsp;&nbsp;The Gigamunch Team</p>
    </div>
  </div>

</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Image169 from '../../components/Image169.vue';
import Dish from '../../components/Dish.vue';

@Component({
  components: {
    Image169,
    Dish,
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
    return 'Your Journey to ' + this.exe.culture.country;
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

  get dishes(): Common.Dish[] {
    return this.exe.dishes;
  }

  get playlistTitle(): string {
    return this.exe.culture.nationality + ' Music Playlist';
  }

  get spotifyUrl(): string {
    return this.exe.content.spotify_url;
  }

  get youtubeUrl(): string {
    return this.exe.content.youtube_url;
  }

  get cookImageSrc(): string {
    return this.exe.content.cook_image_url;
  }

  get cookTitle(): string {
    return this.exe.culture_cook.first_name + "'s Story";
  }

  get cookStory(): string {
    return this.exe.culture_cook.story;
  }
}
</script>
<style scoped lang="scss">
.view {
  max-width: 850px;
  margin: auto;
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
  font-size: 20px;
  font-weight: 700;
  top: 24px;
  left: 28px;
}
@media (min-width: 700px) {
  .hero-image-text {
    font-size: 40px;
  }
}

.content-container {
  padding: 0 $view-edge-padding-mobile 25px $view-edge-padding-mobile;
  box-shadow: 0 0 15px grey;
}
@media (min-width: 700px) {
  .content-container {
    padding: 0 $view-edge-padding-desktop 25px $view-edge-padding-desktop;
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
  width: 75px;
  height: 75px;
  background-size: cover;
  border-radius: 50%;
}
.host-text {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding-left: 12px;
}
.host-text-name {
  font-size: 15px;
  margin: 6px 0 0 0;
}
.host-text-hosted-by {
  font-size: 13px;
  margin: 0;
  opacity: 0.75;
}
@media (min-width: 500px) {
  .host-text-name {
    font-size: 24px;
  }
  .host-text-hosted-by {
    margin: 0 0 6px 0;
    font-size: 15px;
  }
  .host-image-image {
    width: 100px;
    height: 100px;
  }
}

.dinner-image-title-text {
  margin: 24px 0 12px 0;
}

.divider-line {
  margin: 50px 0px;
  border: 0;
  border-bottom: 1px solid #dadfe1;
}

.section-title {
  margin: 0 0 12px 0;
}

@media (min-width: 500px) {
  .section-title {
    font-size: 20px;
  }
}

.cook-story {
  margin: 24px 6px;
  text-align: justify;
}

.footer-message {
  padding: 50px 0 50px 0;
  // background-color: #dadfe1;
  align-content: center;
}

.footer-message-text {
  align-content: center;
  text-align: center;
  margin: 0;
  font-size: 16px;
}
</style>
