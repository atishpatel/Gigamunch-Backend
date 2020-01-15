<template>
  <div class="wrapper">
    <SendMessage v-on:send-message-success="getLogs"></SendMessage>
    <SubscriberLogs :logs="logs"></SubscriberLogs>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { GetLogsByAction } from '../ts/service';
import { GetLogsExtended } from '../ts/extended';
import { IsError } from '../ts/errors';
import SubscriberLogs from '../components/SubscriberLogs.vue';
import SendMessage from '../components/SendMessage.vue';

@Component({
  components: {
    SubscriberLogs,
    SendMessage,
  },
})
export default class Messages extends Vue {
  protected logs: Types.LogExtended[];

  public constructor() {
    super();
    this.logs = [];
  }

  public created() {
    this.getLogs();
  }

  public getLogs() {
    GetLogsByAction(0, 100, 'Message').then((resp) => {
      if (IsError(resp)) {
        console.error(resp);
        return;
      }
      this.logs = GetLogsExtended(resp.logs);
    });
  }
}
</script>
<style lang="scss">
.wrapper {
  background-color: white;
  max-width: 1500px;
  margin: auto;
  padding: 24px;
}
</style>
