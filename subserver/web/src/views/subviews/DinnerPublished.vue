<template>
  <div
    class="background"
    :style="{ backgroundImage: 'url(\'' + patternImageSrc + '\')' }"
  >

    <div class="white-filter">
      <div class="view">
        <div class="hero-image">
          <div class="hero-image-gradient">
            <Image169 :src="exe.content.landscape_image_url"></Image169>
          </div>
          <div
<<<<<<< HEAD
            class="host-image-image"
=======
            class="hero-image-text"
>>>>>>> feature/sub-app-v1
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
          <div class="servings">
            <p class="servings-text">
              {{servingText}}
            </p>
          </div>
          <div
            v-if="!userSummary.has_subscribed"
            class="buttons-row"
          >
            <v-btn
              depressed
              color="#E8554E"
              class="white--text"
            >Get this dinner</v-btn>
          </div>
          <div
            v-else
            class="buttons-row"
          >
            <v-btn
              depressed
              color="#E8554E"
              class="white--text"
              :disabled="disableSkip"
              @click="skipClicked"
            >{{skipButtonText}}</v-btn>
            <ButtonChangeServings
              :activity="activity"
              v-on:dialog-success="updatedServings"
              depressed
              color="#E8554E"
              class="white--text"
              :ButtonDisabled="disableChangeServings"
            ></ButtonChangeServings>
            <v-btn
              depressed
              color="#E8554E"
              class="white--text"
              @click="seeVegClicked"
            >{{seeVegButtonText}}</v-btn>
          </div>
          <div
            v-for="dish in dishes"
            v-bind:key="dish.name"
          >
            <Dish
              :name="dish.name"
              :description="dish.description"
              :ingredients="dish.ingredients"
            ></Dish>
          </div>
          <hr class="divider-line">
          <div class="section-title">
            <h2 class="playlist-title-text">{{playlistTitle}}</h2>
          </div>
          <p>A Gigamunch dinner isn’t complete without some cultural music to listen to while you eat.</p>
          <div class="buttons-row">
            <v-btn
              depressed
              large
              color="#E8554E"
              class="white--text"
              :href="spotifyUrl"
              target="_blank"
            >Listen on Spotfiy</v-btn>
            <v-btn
              depressed
              large
              color="#E8554E"
              class="white--text"
              :href="youtubeUrl"
              target="_blank"
            >Listen on Youtube</v-btn>
          </div>
          <hr class="divider-line">
          <div class="section-title">
            <h2 class="cook-title-text">{{cookTitle}}</h2>
          </div>
          <Image169
            :src="cookImageSrc"
            :rounded=true
          ></Image169>
          <div>
            <p class="cook-story">{{cookStory}}</p>
          </div>
          <hr class="divider-line">
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
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Image169 from '../../components/Image169.vue';
import Dish from '../../components/Dish.vue';
import ButtonChangeServings from '../../components/ButtonChangeServings.vue';
import { GetDayMonthDayDate } from '../../ts/utils';
import { IsError, ErrorAlert } from '../../ts/errors';
import { SkipActivity } from '../../ts/service';
import { UnskipActivity } from '../../ts/service';

@Component({
  components: {
    Image169,
    Dish,
    ButtonChangeServings,
  },
})
export default class DinnerPublished extends Vue {
  @Prop()
  public exe!: Common.Execution;
  @Prop()
  public activity!: Common.Activity;
  @Prop()
  public userSummary!: SubAPI.GetUserSummaryResp;

  public showingVegetarianDinner = false;
  protected disableSkip = false;

  public setVegShowing(v: boolean) {
    this.showingVegetarianDinner = v;
  }

  get disableChangeServings(): boolean {
    if (this.activity) {
      if (this.activity.skip) {
        return true;
      }
    }
    return false;
  }

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
    if (this.showingVegetarianDinner) {
      return `${this.exe.culture_cook.first_name}'s Vegetarian Dinner`;
    } else {
      return `${this.exe.culture_cook.first_name}'s Dinner`;
    }
  }

  get seeVegButtonText(): string {
    if (this.showingVegetarianDinner) {
      return 'See meat option';
    } else {
      return 'See vegetarian option';
    }
  }

  get dinnerImageSrc(): string {
    if (this.showingVegetarianDinner) {
      return this.exe.content.hands_plate_veg_image_url;
    } else {
      return this.exe.content.hands_plate_non_veg_image_url;
    }
  }

  get dishes(): Common.Dish[] {
    const vegDishes: Common.Dish[] = [];
    const meatDishes: Common.Dish[] = [];

    if (this.exe && this.exe.dishes) {
      for (let i = 0; i < this.exe.dishes.length; i++) {
        const dish = this.exe.dishes[i];
        if (dish.is_for_non_vegetarian) {
          meatDishes.push(dish);
        }
        if (dish.is_for_vegetarian) {
          vegDishes.push(dish);
        }
      }
    }

    if (this.showingVegetarianDinner) {
      return vegDishes;
    } else {
      return meatDishes;
    }
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
    return `${this.exe.culture_cook.first_name}'s Story`;
  }

  get cookStory(): string {
    return this.exe.culture_cook.story;
  }

  get servingText(): string {
    if (!this.userSummary.has_subscribed) {
      return `You could receive this dinner on ${GetDayMonthDayDate(
        this.exe.date
      )}`;
    } else {
      if (this.activity) {
        if (this.activity.skip) {
          return `You will not receive this dinner on ${GetDayMonthDayDate(
            this.exe.date
          )}`;
        } else {
          if (
            this.activity.servings_non_vegetarian > 0 &&
            this.activity.servings_vegetarian > 0
          ) {
            return `You will receive ${
              this.activity.servings_non_vegetarian
            } meat servings and ${
              this.activity.servings_vegetarian
            } vegetarian servings on ${GetDayMonthDayDate(this.exe.date)}`;
          } else if (this.activity.servings_vegetarian > 0) {
            return `You will receive ${
              this.activity.servings_vegetarian
            } vegetarian servings on ${GetDayMonthDayDate(this.exe.date)}`;
          } else if (this.showingVegetarianDinner) {
            return `You will receive ${
              this.activity.servings_non_vegetarian
            } meat servings on ${GetDayMonthDayDate(this.exe.date)}`;
          } else {
            return `You will receive ${
              this.activity.servings_non_vegetarian
            } servings on ${GetDayMonthDayDate(this.exe.date)}`;
          }
        }
      }
      return '';
    }
  }

  get skipButtonText(): string {
    if (this.activity) {
      if (this.activity.skip) {
        return 'Unskip';
      } else {
        return 'Skip';
      }
    }
    return 'Skip';
  }

  get patternImageSrc(): string {
    if (this.exe && this.exe.content && this.exe.content.cook_image_url) {
      const originalSrc = this.exe.content.cover_image_url;
      if (originalSrc.includes('cook.jpg')) {
        return originalSrc.replace('cook.jpg', 'cover.jpg');
      }
    }
    return '';
  }

  protected skipClicked() {
    if (!this.activity) {
      alert('activity not found');
      return;
    }
    this.disableSkip = true;
    if (this.activity.skip) {
      UnskipActivity(this.activity.date).then((resp) => {
        if (IsError(resp)) {
          ErrorAlert(resp);
          window.location.reload();
          return;
        }
        this.$emit('get-activity');
        this.disableSkip = false;
      });
    } else {
      SkipActivity(this.activity.date).then((resp) => {
        if (IsError(resp)) {
          ErrorAlert(resp);
          window.location.reload();
          return;
        }

        this.$emit('get-activity');
        this.disableSkip = false;
      });
    }
    this.activity.skip = true;
  }

  protected updatedServings() {
    this.$emit('get-activity');
  }

  protected seeVegClicked() {
    if (this.showingVegetarianDinner) {
      this.showingVegetarianDinner = false;
    } else {
      this.showingVegetarianDinner = true;
    }
  }
}
</script>
<style scoped lang="scss">
.background {
  background-repeat: repeat;
  background-size: 800px;
}
.white-filter {
  background-color: rgba(255, 255, 255, 0.9);
}
.view {
  max-width: 850px;
  margin: auto;
}

$view-edge-padding-desktop: 120px;
$view-edge-padding-mobile: 24px;

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
  background-color: white;
  padding: 0 $view-edge-padding-mobile 25px $view-edge-padding-mobile;
  box-shadow: 0 0 15px grey;
}
@media (min-width: 700px) {
  .content-container {
    padding: 0 $view-edge-padding-desktop 25px $view-edge-padding-desktop;
  }
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
  color: #869995;
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

.servings-text {
  text-align: center;
  margin: 28px 0 8px 0;
  font-size: 14px;
  font-weight: 700;
}
@media (min-width: 550px) {
  .servings-text {
    font-size: 20px;
  }
}

.buttons-row {
  text-align: center;
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
  padding: 0 0 50px 0;
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