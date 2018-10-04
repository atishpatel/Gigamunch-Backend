<template>
  <div class="home">
    <div>
      <!-- <img alt="Vue logo" src="../assets/logo.png"> -->
      <HelloWorld ref="helloWorld" msg="Welcome to Your Vue.js + TypeScrit App" />
      <button @click="getExecutions">test button</button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import HelloWorld from '../components/HelloWorld.vue';
import { GetExecutions } from '../ts/service';

@Component({
  components: {
    HelloWorld,
  },
})
export default class Dinners extends Vue {
  public created() {
    this.getExecutions();
  }

  public getExecutions() {
    GetExecutions(0, 10).then((resp) => {
      console.log('GetExecutions: ', resp);
    });
  }
}
</script>
