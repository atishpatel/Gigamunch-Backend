<template>
  <div>
    <div>
      <ExecutionsList class="list-component" :executionAndActivityList="executionAndActivityList"></ExecutionsList>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import ExecutionsList from '../components/ExecutionsList.vue';
import { GetExecutions } from '../ts/service';
import { IsError } from '../ts/errors';

@Component({
  components: {
    ExecutionsList,
  },
})
export default class Dinners extends Vue {
  protected executionAndActivityList: Common.ExecutionAndActivity[];

  public constructor() {
    super();
    this.executionAndActivityList = [];
  }

  public created() {
    this.getExecutions();
  }

  public getExecutions() {
    GetExecutions(0, 10).then((resp) => {
      if (IsError(resp)) {
        return;
      }
      this.executionAndActivityList = resp.execution_and_activity;
    });
  }
}
</script>
<style scoped lang="scss">
.list-component {
}
</style>
