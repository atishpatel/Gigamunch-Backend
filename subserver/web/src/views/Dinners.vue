<template>
  <div>
    <div>
      <ExecutionsList
        class="list-component"
        :executionAndActivityList="executionAndActivityList"
      ></ExecutionsList>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import ExecutionsList from '../components/ExecutionsList.vue';
import { GetExecutionsAfterDate } from '../ts/service';
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
    let today = new Date();
    today.setHours(today.getHours() - 6);
    GetExecutionsAfterDate(today).then((resp) => {
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
