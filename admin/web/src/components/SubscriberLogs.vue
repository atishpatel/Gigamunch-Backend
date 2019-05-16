<template>
  <div>
    <v-timeline
      dense
      clipped
    >
      <v-timeline-item
        v-for="(log, i) in logs"
        :key="i"
        :color="log.color"
        :icon="log.icon"
        small
      >
        <template v-slot:opposite>
          <span
            :class="`headline font-weight-bold ${log.color}--text`"
            v-text="log.typestampString"
          ></span>
        </template>
        <v-card
          :color="log.color"
          dark
        >
          <v-card-title class="title">{{log.basic_payload.title}}</v-card-title>
          <v-card-text class="white text--primary">
            <p v-html="log.basicPayloadDescriptionHTML">Lorem ipsum dolor sit amet, no nam oblique veritus. Commune scaevola imperdiet nec ut, sed euismod convenire principes at. Est et nobis iisque percipit, an vim zril disputando voluptatibus, vix an salutandi sententiae.</p>
          </v-card-text>
        </v-card>
      </v-timeline-item>
    </v-timeline>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';

@Component({
  components: {},
})
export default class SubscriberLogs extends Vue {
  @Prop()
  public logs!: Types.LogExtended[];

  @Watch('logs')
  logsWatcher(l: Types.LogExtended[]) {
    const logs = l;
    for (var i = 0; i < logs.length; i++) {
      logs[i].color = this.getColor(logs[i].type, logs[i].action);
      logs[i].icon = this.getIcon(logs[i].type, logs[i].action);
      logs[i].timestampString = this.getTimestampString(logs[i].timestamp);
      logs[i].basicPayloadDescriptionHTML = logs[
        i
      ].basic_payload.description.replace(/;;;/g, '<br>');
    }
    this.logs = logs;
  }

  getColor(type: string, action: string) {
    switch (action) {
      case 'skip':
        return 'orange';
      case 'unskip':
        return 'pink';
      case 'message':
        return 'green';
      case 'rating':
        return 'cyan';
      case 'update':
        return 'amber';
    }
    return 'bubble_chart';
  }

  getIcon(type: string, action: string) {
    switch (action) {
      case 'skip':
        return 'remove_shopping_cart';
      case 'unskip':
        return 'add_shopping_cart';
      case 'message':
        return 'message';
      case 'rating':
        return 'star_rate';
      case 'update':
        return 'cloud_upload';
    }
    return 'bubble_chart';
  }

  getTimestampString(dateString: string) {
    const monthNames = [
      'Jan',
      'Feb',
      'Mar',
      'Apr',
      'May',
      'Jun',
      'Jul',
      'Aug',
      'Sep',
      'Oct',
      'Nov',
      'Dec',
    ];
    const dayNames = [
      'Sunday',
      'Monday',
      'Tuesday',
      'Wedensday',
      'Thursday',
      'Friday',
      'Saturday',
    ];
    const d = new Date(dateString);
    let day = d.getDay();
    let month = d.getMonth();
    let date = d.getDate();
    let year = d.getFullYear();
    return `${dayNames[day]}, ${
      monthNames[month]
    } ${date}, ${year} @ ${d.toLocaleTimeString()}`;
  }
}
</script>

<style scoped lang="scss">
</style>
