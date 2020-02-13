<template>
  <div>
    <h2 class="date">{{titleDate}}</h2>
    <div class="cards">
      <!-- Culture Card -->
      <ExecutionsCard
        class="card"
        :cook_name="cookName"
        :dinner_image_src="execution.content.hands_plate_non_veg_image_url"
        :cook_face_image_src="execution.email.cook_face_image_url"
        :to="{path: 'dinner/'+executionURLID+'#culture'}"
      ></ExecutionsCard>

    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import ExecutionsCard from './ExecutionsCard.vue';
import { GetDayMonthDayDate } from '../ts/utils';

@Component({
  components: {
    ExecutionsCard,
  },
})
export default class ExecutionsItem extends Vue {
  @Prop()
  public execution!: Common.Execution;

  // computed
  get titleDate() {
    return GetDayMonthDayDate(this.execution.date);
  }

  get cookName() {
    return (
      this.execution.culture_cook.first_name +
      ' ' +
      this.execution.culture_cook.last_name
    );
  }

  get cultureTitle() {
    return this.execution.culture.greeting;
  }

  get cultureDescription() {
    if (this.execution.culture.description.length > 60) {
      return this.execution.culture.description.substr(0, 60) + '...';
    }
    return this.execution.culture.description;
  }

  get cookDescription() {
    if (this.execution.culture_cook.story.length > 60) {
      return this.execution.culture_cook.story.substr(0, 60) + '...';
    }
    return this.execution.culture_cook.story;
  }

  get executionURLID() {
    if (this.execution.date) {
      return this.execution.date;
    }
    return this.execution.id;
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.date {
  font-weight: 500;
  font-size: 1.5em;
  padding-left: 24px;
}
.cards {
  display: flex;
  flex-wrap: nowrap;
  overflow-x: auto;
  transition: 0.5s ease 0s;
  padding: 0px 24px 30px;
  -webkit-overflow-scrolling: touch;
  .card {
    flex: 0 0 auto;
    max-width: 77vw;
    min-height: 250px;
  }
  @media (min-width: 800px) {
    .card {
      max-width: 500px;
      width: 500px;
    }
  }
  .card:not(:first-child) {
    margin-left: 24px;
  }
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
