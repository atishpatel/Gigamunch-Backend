<template>
  <div>
    <h2 class="culture-title">{{date}} – {{culture}}</h2>
    <div class="execution-item-row">
      <!-- <div class="date-section">
        <h2 class="date">{{date}}</h2>
      </div> -->
      <!-- Culture Card -->
      <ExecutionsCard
        class="card"
        :cookName="cookName"
        :nationality="nationality"
        :dinnerImageSource="execution.content.hands_plate_non_veg_image_url"
        :cookFaceImageSource="execution.email.cook_face_image_url"
        :to="{path: 'dinner/'+executionURLID+'#culture'}"
      ></ExecutionsCard>
      <!-- <div class="dishes-section">
      </div> -->
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
  get date() {
    return GetDayMonthDayDate(this.execution.date);
  }

  get culture() {
    return this.execution.culture.country;
  }

  get cookName() {
    return (
      this.execution.culture_cook.first_name +
      ' ' +
      this.execution.culture_cook.last_name
    );
  }

  get nationality() {
    return this.execution.culture.nationality;
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
// .date {
//   flex: 0 0 auto;
//   font-weight: 500;
//   font-size: 1.5em;
//   padding-left: 24px;
//   margin: 40px 0 15px 0;
// }

.culture-title {
  font-weight: 500;
  font-size: 1.2em;
  padding-left: 24px;
  margin: 40px 0 15px 0;
}
@media (min-width: 800px) {
  .culture-title {
    font-size: 2em;
  }
}
.execution-item-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  flex-wrap: nowrap;
  overflow-x: auto;
  transition: 0.5s ease 0s;
  padding: 0px 24px 30px;
  -webkit-overflow-scrolling: touch;
  .card {
    flex: 0 0 auto;
    width: 100%;
    max-width: 100%;
    min-height: 250px;
  }
  @media (min-width: 800px) {
    .card {
      max-width: 500px;
      width: 500px;
    }
  }

  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
