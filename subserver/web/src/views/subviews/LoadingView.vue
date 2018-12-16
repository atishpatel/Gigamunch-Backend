<template>
  <div class="loading-view">
    <div class="loader"></div>
    Loading... TODO:
  </div>
</template>

<script lang="ts">
import { Vue } from 'vue-property-decorator';

export default class LoadingView extends Vue {}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.loading-view {
  text-align: center;
  padding: 25vh;
  // TODO:
}

.loading-view[hide] {
  display: none;
}

.loader {
  margin: auto;
  border: 16px solid #f3f3f3; /* Light grey */
  border-top: 16px solid #3498db; /* Blue */
  border-radius: 50%;
  width: 120px;
  height: 120px;
  animation: spin 2s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>