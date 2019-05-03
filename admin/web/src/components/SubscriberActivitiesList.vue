<template>
  <div>
    <v-card>
      <v-card-title>
        Activities
        <v-spacer></v-spacer>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="activities"
        :pagination.sync="pagination"
      >
        <template v-slot:items="props">
          <td>{{ props.item.dateString }}</td>
          <td class="text-xs-right">{{ props.item.non_veg_servings }}</td>
          <td class="text-xs-right">{{ props.item.veg_servings }}</td>
          <td class="text-xs-right">{{ props.item.transaction_id }}</td>
          <td class="text-xs-right">${{ props.item.discountAmount }} | {{ props.item.discountPercent}}%</td>
          <td class="text-xs-right">{{ props.item.status }}</td>
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
  public headers = [
    { text: 'Date', value: 'date', sortable: false },
    { text: 'Servings', value: 'non_veg_servings', sortable: false },
    { text: 'Veg Servings', value: 'veg_servings', sortable: false },
    { text: 'TransactionID', value: 'transaction_id', sortable: false },
    { text: 'Discount', value: 'discountPercent', sortable: false },
    { text: 'Status', value: 'status', sortable: false },
  ];
}
</script>

<style scoped lang="scss">
</style>
