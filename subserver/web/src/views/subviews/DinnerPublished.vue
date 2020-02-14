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

  get cultureCookName(): string {
    return (
      this.exe.culture_cook.first_name + ' ' + this.exe.culture_cook.last_name
    );
  }
}
</script>
<style scoped lang="scss">
.view {
  max-width: 1024px;
  margin: auto;
}

$view-edge-padding: 24px;
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
    rgba(100, 100, 100, 0.45),
    rgba(0, 0, 0, 0)
  );
  height: 15vw;
  width: 100%;
  content: '';
}
.hero-image-text {
  position: absolute;
  color: white;
  font-size: 24px;
  font-weight: 700;
  top: 12px;
  left: 12px;
}
// host-action
.host-action {
}
.host {
  display: flex;
  position: relative;
  top: -35px;
  padding-left: $view-edge-padding;
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
  margin: 0;
  opacity: 0.75;
}
.host-text-name {
  margin: 6px 0;
}
</style>
