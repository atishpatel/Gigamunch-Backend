<template>
  <div>
    <h1>{{execution.date}}</h1>
    <p>{{execution.culture}}</p>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ExecutionsItem extends Vue {
  @Prop()
  public execution!: Common.Execution;

  public created() {}
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
</style>
