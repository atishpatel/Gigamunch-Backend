<template>
  <div>
    <v-card>
      <v-card-title>
        Unpaid Summary
        <v-spacer></v-spacer>
        <v-text-field
          v-model="search"
          append-icon="search"
          label="Search"
          single-line
          hide-details
        ></v-text-field>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="summaries"
        :search="search"
        :pagination.sync="pagination"
      >
        <template v-slot:items="props">
          <td>
            <router-link :to="'/subscriber/' + props.item.email"> {{ props.item.email }} </router-link>
          </td>
          <td class="text-xs-right">{{ props.item.user_id }}</td>
          <td class="text-xs-right">{{ props.item.min_date }}</td>
          <td class="text-xs-right">{{ props.item.max_date }}</td>
          <td class="text-xs-right">{{ props.item.name }}</td>
          <td class="text-xs-right">{{ props.item.num_unpaid }}</td>
          <td class="text-xs-right">{{ props.item.amount_due }}</td>
        </template>
        <template v-slot:no-results>
          <v-alert
            :value="true"
            color="error"
            icon="warning"
          >
            Your search for "{{ search }}" found no results.
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
export default class UnpaidSummaryList extends Vue {
  @Prop()
  public summaries!: Types.UnpaidSummaryExtended[];
  public search = '';
  public pagination = {
    rowsPerPage: -1,
  };
  public headers = [
    { text: 'First Unpaid', value: 'min_date' },
    { text: 'Last Unpaid', value: 'max_date' },
    { text: 'Name', value: 'name' },
    { text: 'Email', value: 'email' },
    { text: 'ID', value: 'user_id' },
    { text: 'Num Unpaid', value: 'num_unpaid' },
    { text: 'Amount Due', value: 'amount_due' },
  ];
}
</script>

<style scoped lang="scss">
</style>
