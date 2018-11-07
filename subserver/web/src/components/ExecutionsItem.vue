<template>
  <div>
    <h2 class="date">{{titleDate}}</h2>
    <div class="cards">
      <!-- Culture Card -->
      <ExecutionsCard class="card" :title="cultureTitle" :description="cultureDescription" :src="execution.content.hero_image_url"></ExecutionsCard>
      <!-- Cook Card -->
      <ExecutionsCard class="card" :title="cookName" :description="cookDescription" :src="execution.content.cook_image_url"></ExecutionsCard>

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

  get cultureTitle() {
    return this.execution.culture.greeting;
  }

  get cultureDescription() {
    if (this.execution.culture.description.length > 60) {
      return this.execution.culture.description.substr(0, 60) + '...';
    }
    return this.execution.culture.description;
  }

  get cookName() {
    return `Meet ${this.execution.culture_cook.first_name} ${
      this.execution.culture_cook.last_name
    }`;
  }

  get cookDescription() {
    if (this.execution.culture_cook.story.length > 60) {
      return this.execution.culture_cook.story.substr(0, 60) + '...';
    }
    return this.execution.culture_cook.story;
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
  // -webkit-overflow-scrolling: touch;
  .card {
    flex: 0 0 auto;
    max-width: 77vw;
  }
  @media (min-width: 800px) {
    .card {
      max-width: 500px;
    }
  }
  .card:not(:first-child) {
    padding-left: 24px;
  }
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
