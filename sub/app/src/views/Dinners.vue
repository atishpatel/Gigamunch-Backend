<template>
  <div class="home">
    <div>
      <HelloWorld ref="helloWorld" msg="Welcome to Your Vue.js + TypeScrit App" />
      <button @click="getExecutions">test button</button>
      <ExecutionsList :executions="executions"></ExecutionsList>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import HelloWorld from '../components/HelloWorld.vue';
import ExecutionsList from '../components/ExecutionsList.vue';
import { GetExecutions } from '../ts/service';
import { IsError } from '../ts/errors';

@Component({
  components: {
    HelloWorld,
    ExecutionsList,
  },
})
export default class Dinners extends Vue {
  protected executions: Array<Common.Execution>;

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
      console.log('GetExecutions: ', resp);
      this.executions = resp.executions;
    });
  }
}
</script>
