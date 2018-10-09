<template>
    <div>
        <h1>Executions</h1>
        <div v-for="execution in executions" v-bind:key="execution.id">
            <ExecutionsItem :execution="execution"></ExecutionsItem>
        </div>
    </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import ExecutionsItem from './ExecutionsItem.vue';

@Component({
  components: {
    ExecutionsItem,
  },
})
export default class ExecutionsList extends Vue {
  @Prop()
  public executions!: Array<Common.Execution>;

  public created() {}
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
</style>
