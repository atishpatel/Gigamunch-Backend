<template>
  <div class="el-container">
    <div class="activities">
      <v-card-title>
        Activities
        <v-spacer></v-spacer>
        <ButtonSetupActivity v-on:dialog-success="getActivities"></ButtonSetupActivity>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="activities"
        :pagination.sync="pagination"
        :expand="expand"
        item-key="date"
      >
        <template v-slot:items="props">
          <tr @click="props.expanded = !props.expanded">
            <td>{{ props.item.dateFull }}</td>
            <td>{{ props.item.servings_non_vegetarian }}</td>
            <td>{{ props.item.servings_vegetarian }}</td>
            <td>{{ props.item.transaction_id }}</td>
            <td>{{ props.item.discountString }}</td>
            <td>{{ props.item.status }}<br><span class="paid-date">{{ props.item.paidDate }}</span></td>
          </tr>
        </template>

        <template v-slot:expand="props">
          <v-card
            flat
            class="actions"
          >
            <ButtonSkip
              :activity="props.item"
              v-on:dialog-success="getActivities"
            ></ButtonSkip>
            <ButtonUnskip
              :activity="props.item"
              v-on:dialog-success="getActivities"
            ></ButtonUnskip>
            <v-btn
              outline
              round
              :disabled="!props.item.paid"
              @click="processAcitivity(props.item)"
            >Process</v-btn>

            <ButtonChangeServings
              :activity="props.item"
              v-on:dialog-success="getActivities"
            ></ButtonChangeServings>
          </v-card>
          <v-card-text v-html="detailHTML(props.item)"></v-card-text>
        </template>

        <template v-slot:no-results>
          <v-alert
            :value="true"
            color="warn"
            icon="warning"
          >
            Subscriber has no actvitities.
          </v-alert>
        </template>
      </v-data-table>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import ButtonSkip from './ButtonSkip.vue';
import ButtonUnskip from './ButtonUnskip.vue';
import ButtonSetupActivity from './ButtonSetupActivity.vue';
import ButtonChangeServings from './ButtonChangeServings.vue';

@Component({
  components: {
    ButtonSkip,
    ButtonUnskip,
    ButtonSetupActivity,
    ButtonChangeServings,
  },
})
export default class SubscriberActivitiesList extends Vue {
  @Prop()
  public activities!: Types.ActivityExtended[];
  public search = '';
  public pagination = {
    rowsPerPage: -1,
  };
  public expand = true;
  public headers = [
    { text: 'Date', value: 'date', sortable: false },
    { text: 'Servings', value: 'servings_non_vegetarian', sortable: false },
    { text: 'Veg Servings', value: 'servings_vegetarian', sortable: false },
    { text: 'TransactionID', value: 'transaction_id', sortable: false },
    { text: 'Discount', value: 'discountPercent', sortable: false },
    { text: 'Status', value: 'status', sortable: false },
  ];

  protected detailHTML(act: any) {
    let table = '<table>';
    const keys = Object.keys(act);
    for (let i = 0; i < keys.length; i++) {
      table += '<tr>';
      while (!act[keys[i]] && i < keys.length) {
        i++;
      }
      if (act[keys[i]]) {
        table += `<td class="font-weight-bold">${keys[i]}</td><td>${
          act[keys[i]]
        }</td>`;
      }
      i++;
      while (!act[keys[i]] && i < keys.length) {
        i++;
      }
      if (act[keys[i]]) {
        table += `<td class="font-weight-bold">${keys[i]}</td><td>${
          act[keys[i]]
        }</td>`;
      }
      table += '</tr>';
    }
    return table + '</table>';
  }

  protected processAcitivity(act: Types.ActivityExtended) {
    console.log(act);
  }

  protected getActivities() {
    this.$emit('get-activities');
  }
}
</script>

<style scoped lang="scss">
.paid-date {
  color: #bbbbbb;
}

.actions {
  display: flex;
  justify-content: center;
  padding: 24px;
}

.el-container {
  background: white;
}

.activities {
  // max-width: 1500px;
  // margin: auto;
  border: 1px solid #dadce0;
  border-radius: 8px;
  overflow: hidden;
}
</style>
