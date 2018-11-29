<template>
  <div>
    <h1 class="title">Upcomming Dinners</h1>
    <div v-for="executionAndActivity in executionAndActivityList" v-bind:key="executionAndActivity.execution.id">
      <ExecutionsItem :execution="executionAndActivity.execution"></ExecutionsItem>
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
  public executionAndActivityList!: Common.Execution[];
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
@import '../scss/theme.scss';
.title {
  display: inline-block;
  text-align: start;
  margin: 0 24px;
  padding-bottom: 6px;
  font-size: 1.5em;
  color: $mdc-theme-secondary;
  border-bottom: 4px solid $mdc-theme-accent;
}
</style>
