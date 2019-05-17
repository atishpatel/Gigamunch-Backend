<template>
  <div class="el-container">
    <v-card class="discounts">
      <v-card-title>
        Discounts
        <v-spacer></v-spacer>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="discounts"
        :pagination.sync="pagination"
        item-key="id"
      >
        <template v-slot:items="props">
          <tr>
            <td>{{ props.item.id }}</td>
            <td class="text-secondary">{{ props.item.created_datetime }}</td>
            <td class="text-secondary">{{ props.item.user_id }}</td>
            <td class="text-secondary">{{ props.item.email }}</td>
            <td>{{ getDateUsed(props.item.date_used) }}</td>
            <td>{{ props.item.discount_amount }}</td>
            <td>{{ props.item.discount_percent }}</td>
          </tr>
        </template>

        <template v-slot:no-results>
          <v-alert
            :value="true"
            color="warn"
            icon="warning"
          >
            Subscriber has no discounts.
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
export default class SubscriberDiscountsList extends Vue {
  @Prop()
  public discounts!: Common.Discount[];
  public pagination = {
    rowsPerPage: -1,
  };
  public expand = true;
  public headers = [
    { text: 'ID', value: 'id', sortable: false },
    { text: 'Created datetime', value: 'created_datetime', sortable: false },
    { text: 'UserID', value: 'user_id', sortable: false },
    { text: 'Email', value: 'email', sortable: false },
    { text: 'DateUsed', value: 'date_used', sortable: false },
    { text: 'Amount', value: 'discount_amount', sortable: false },
    { text: 'Percent', value: 'discount_percent', sortable: false },
  ];

  protected getDateUsed(dateUsed: string): string {
    if (dateUsed.includes('0001-01-01')) {
      return '- unused -';
    }
    return dateUsed;
  }
}
</script>

<style scoped lang="scss">
.text-secondary {
  color: #bbbbbb;
}

.el-container {
  background: white;
}

.discounts {
  max-width: 1500px;
  margin: auto;
}
</style>
