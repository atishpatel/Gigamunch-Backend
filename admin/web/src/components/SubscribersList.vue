<template>
  <div>
    <v-card>
      <v-card-title>
        Subscribers
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
        :items="subs"
        :search="search"
        :pagination.sync="pagination"
      >
        <template v-slot:items="props">
          <td>{{ props.item.emailsString }}</td>
          <td class="text-xs-right">{{ props.item.namesString }}</td>
          <td class="text-xs-right">{{ props.item.phonenumbersString }}</td>
          <td class="text-xs-right">{{ props.item.addressString }}</td>
          <td class="text-xs-right">{{ props.item.id }}</td>
          <td class="text-xs-right">{{ props.item.active }}</td>
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
export default class SubscribersList extends Vue {
  @Prop()
  public subs!: Common.Subscriber[];
  public search = '';
  public pagination = {
    rowsPerPage: -1,
  };
  public headers = [
    { text: 'Emails', value: 'emailsString' },
    { text: 'Name', value: 'name' },
    { text: 'Phone Numbers', value: 'phonenumbersString' },
    { text: 'Address', value: 'addressString' },
    { text: 'ID', value: 'id' },
    { text: 'Active', value: 'active' },
  ];
}
</script>

<style scoped lang="scss">
</style>
