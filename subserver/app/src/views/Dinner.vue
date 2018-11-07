<template>
  <div>
    <div v-html="execution"></div>
    <div v-html="execution.culture_cook.story"></div>
    <div v-html="execution.culture_cook.story_preview"></div>
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
export default class Dinner extends Vue {
  protected execution!: Common.Execution;

  public constructor() {
    super();
    this.execution = {} as Common.Execution;
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
