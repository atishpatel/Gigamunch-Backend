<template>
  <div>
    <!-- v-if="userSummary.has_subscribed" -->
    <DinnerPublished
      :exe="exe"
      :activity="activity"
      :userSummary="userSummary"
    ></DinnerPublished>
    <div class="hero-image unset-link">
      <img
        :src="landscape_image_src"
        :alt="computedLandscapeImageAlt"
      >
    </div>

  </div>
</template>

<script lang="ts">
import { Prop, Component, Vue } from 'vue-property-decorator';
import DinnerPublished from './subviews/DinnerPublished.vue';
// import DinnerLead from '../subview/DinnerLead.vue';
import { GetExecution } from '../ts/service';
import { IsError } from '../ts/errors';

@Component({
  components: {
    DinnerPublished,
    // DinnerLead,
  },
})
export default class Dinner extends Vue {
  @Prop()
  public userSummary!: SubAPI.GetUserSummaryResp;

  protected exe!: Common.Execution;
  protected loading!: boolean;
  protected activity!: Common.Activity;

  public constructor() {
    super();
    this.exe = {
      culture: {} as Common.Culture,
      content: {} as Common.Content,
      culture_cook: {} as Common.CultureCook,
      culture_guide: {} as Common.CultureGuide,
      dishes: {} as Common.Dish[],
      stickers: {} as Common.Sticker[],
      notifications: {} as Common.Notifications,
      email: {} as Common.Email,
    } as Common.Execution;
    this.activity = {} as Common.Activity;
  }

  public created() {
    this.getExecution(this.$route.params.date);
  }

  public getExecution(id: string) {
    this.loading = true;
    GetExecution(id).then((resp) => {
      this.loading = false;
      if (IsError(resp)) {
        return;
      }
      this.exe = resp.execution_and_activity.execution;
      this.activity = resp.execution_and_activity.activity;
    });
  }
  get computedLandscapeImageAlt() {
    return this.exe.culture.country + 'landscape image';
  }
  get landscape_image_src() {
    return this.exe.content.landscape_image_url;
  }
}
</script>
<style scoped lang="scss">
.hero-image {
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
</style>
