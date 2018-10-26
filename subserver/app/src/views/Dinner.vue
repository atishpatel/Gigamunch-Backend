<template>
  <div>

  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import ExecutionsList from '../components/ExecutionsList.vue';
import { GetExecution } from '../ts/service';
import { IsError } from '../ts/errors';

@Component({
  components: {
    ExecutionsList,
  },
})
export default class Dinners extends Vue {
  protected execution!: Common.Execution;

  public constructor() {
    super();
  }

  public created() {
    this.getExecution(this.$route.params.date);
  }

  public getExecution(id: string) {
    GetExecution(id).then((resp) => {
      if (IsError(resp)) {
        return;
      }
      this.execution = resp.execution;
    });
  }
}
</script>
<style scoped lang="scss">
</style>
