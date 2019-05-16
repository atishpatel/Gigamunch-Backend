<template>
  <div>
    <v-card>
    </v-card>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component({
  components: {},
})
export default class SubscriberSummary extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public logs!: Common.Log[];
}
</script>

<style scoped lang="scss">
</style>
