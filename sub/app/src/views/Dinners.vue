<template>
  <div class="home">
    <div>
      <ExecutionsList class="list-component" :executions="executions"></ExecutionsList>
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
  protected executions: Common.Execution[];

  public constructor() {
    super();
    this.executions = [];
  }

  public created() {
    this.getExecutions();
  }

  public getExecutions() {
    GetExecutions(0, 10).then((resp) => {
      if (IsError(resp)) {
        return;
      }
      this.executions = resp.executions;
    });
  }
}
</script>
<style scoped lang="scss">
.list-component {
}
</style>
