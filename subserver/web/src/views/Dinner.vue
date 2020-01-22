<template>
  <div>
    <!-- v-if="userSummary.has_subscribed" -->
    <DinnerPublished
      :exe="exe"
      :activity="activity"
      :userSummary="userSummary"
    ></DinnerPublished>

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
}
</script>
<style scoped lang="scss">
</style>
