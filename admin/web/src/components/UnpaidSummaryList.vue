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
        :expand="expand"
        item-key="user_id"
      >
        <template v-slot:items="props">
          <tr @click="props.expanded = !props.expanded">
            <td class="text-xs-right">{{ props.item.min_date }}</td>
            <td class="text-xs-right">{{ props.item.max_date }}</td>
            <td class="text-xs-right">{{ props.item.name }}</td>
            <td>
              <router-link :to="'/subscriber/' + props.item.user_id"> {{ props.item.email }} </router-link>
            </td>
            <td class="text-xs-right">{{ props.item.num_unpaid }}</td>
            <td class="text-xs-right">${{ props.item.amount_due }}</td>
          </tr>
        </template>
        <template v-slot:expand="props">
          <table>
            <tr>
              <td>ID</td>
              <td>{{props.item.user_id}}</td>
            </tr>
          </table>
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
  public expand = true;
  public headers = [
    { text: 'First Unpaid', value: 'min_date' },
    { text: 'Last Unpaid', value: 'max_date' },
    { text: 'Name', value: 'name' },
    { text: 'Email', value: 'email' },
    { text: 'Num Unpaid Dinners', value: 'num_unpaid' },
    { text: 'Amount Due', value: 'amount_due' },
  ];
}
</script>

<style scoped lang="scss">
</style>
