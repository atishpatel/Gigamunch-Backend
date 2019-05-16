<template>
  <div class="el-container">
    <v-card class="activities">
      <v-card-title>
        Activities
        <v-spacer></v-spacer>
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
          <v-card flat>
            <v-btn
              outline
              round
              :disabled="props.item.skip"
              @click="skip(props.item)"
            >Skip</v-btn>
            <v-btn
              outline
              round
              :disabled="!props.item.skip"
              @click="unskip(props.item)"
            >Unskip</v-btn>
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
    </v-card>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
// import { SkipActivity, UnskipActivity } from '../ts/service';

@Component({
  components: {},
})
export default class SubscriberActivitiesList extends Vue {
  @Prop()
  public activities!: Types.ActivitiyExtended[];
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

  detailHTML(act: any) {
    console.log(act);
    let table = '<table>';
    const keys = Object.keys(act);
    for (let i = 0; i < keys.length; i++) {
      if (act[keys[i]]) {
        table += `<tr><td class="font-weight-bold">${keys[i]}</td><td>${
          act[keys[i]]
        }</td></tr>`;
      }
    }
    return table + '</table>';
  }

  skip(act: Types.ActivitiyExtended) {
    console.log(act);
  }

  unskip(act: Types.ActivitiyExtended) {
    console.log(act);
  }
}
</script>

<style scoped lang="scss">
.paid-date {
  color: #bbbbbb;
}

.el-container {
  background: white;
}

.activities {
  max-width: 1500px;
  margin: auto;
}
</style>
