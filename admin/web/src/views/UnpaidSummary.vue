<template>
  <div>
    <UnpaidSummaryList :subs="subs"></UnpaidSummaryList>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import UnpaidSummaryList from '../components/UnpaidSummaryList.vue';
import { GetUnpaidSummaries } from '../ts/service';
import { GetUnpaidSummariesExtended } from '../ts/extended';
import { IsError } from '../ts/errors';

@Component({
  components: {
    UnpaidSummaryList,
  },
})
export default class UnpaidSummary extends Vue {
  protected summaries: Types.UnpaidSummaryExtended[];

  public constructor() {
    super();
    this.summaries = [];
  }

  public created() {
    this.getUnpaidSummaries();
  }

  public getUnpaidSummaries() {
    GetUnpaidSummaries().then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      this.summaries = GetUnpaidSummariesExtended(resp.unpaid_summaries);
    });
  }
}
</script>
<style lang="scss">
</style>
