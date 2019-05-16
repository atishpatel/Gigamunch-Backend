<template>
  <div>
    <v-timeline
      class="timeline"
      clipped
    >
      <v-timeline-item
        v-for="(log, i) in logs"
        :key="i"
        :color="log.color"
        :icon="log.icon"
        large
        right
      >
        <template v-slot:opposite>
          <span
            :class="`font-weight-bold ${log.color}--text`"
            v-text="log.timestampString"
          ></span>
        </template>
        <v-card
          :color="log.color"
          dark
        >
          <v-card-title class="title">{{log.basic_payload.title}}</v-card-title>
          <v-card-text class="white text--primary">
            <p
              class="description"
              v-html="log.basicPayloadDescriptionHTML"
            ></p>
          </v-card-text>
        </v-card>
      </v-timeline-item>
    </v-timeline>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component({
  components: {},
})
export default class SubscriberLogs extends Vue {
  @Prop()
  public logs!: Types.LogExtended[];
}
</script>

<style scoped lang="scss">
.timeline {
  max-width: 1100px;
}

.description {
  margin: 0;
}
</style>
