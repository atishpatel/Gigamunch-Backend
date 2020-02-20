<template>
  <div class="view">
    <div class="content-container">
      <h1>Upcoming Dinners</h1>
      <div
        v-for="executionAndActivity in executionAndActivityList"
        v-bind:key="executionAndActivity.execution.id"
      >
        <ExecutionsItem :execution="executionAndActivity.execution"></ExecutionsItem>
      </div>
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
$view-edge-padding-desktop: 120px;
$view-edge-padding-mobile: 24px;
.content-container {
  background-color: white;
  padding: 0 $view-edge-padding-mobile 25px $view-edge-padding-mobile;
}
@media (min-width: 700px) {
  .content-container {
    padding: 0 $view-edge-padding-desktop 25px $view-edge-padding-desktop;
  }
}

.view {
  max-width: 850px;
  margin: auto;
}
</style>
