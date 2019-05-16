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
          <td>{{ props.item.dateFull }}</td>
          <td class="text-xs-right">{{ props.item.servings_non_vegetarian }}</td>
          <td class="text-xs-right">{{ props.item.servings_vegetarian }}</td>
          <td class="text-xs-right">{{ props.item.transaction_id }}</td>
          <td class="text-xs-right">${{ props.item.discount_amount }} | {{ props.item.discount_percent}}%</td>
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
    { text: 'Servings', value: 'servings_non_vegetarian', sortable: false },
    { text: 'Veg Servings', value: 'servings_vegetarian', sortable: false },
    { text: 'TransactionID', value: 'transaction_id', sortable: false },
    { text: 'Discount', value: 'discountPercent', sortable: false },
    { text: 'Status', value: 'status', sortable: false },
  ];

  getStatusString(act) {
    if (act.refunded) {
      return 'Refunded $' + act.refunded_amount;
    } else if (act.skip) {
      return 'Skipped';
    } else if (act.free) {
      return 'First';
    } else if (act.paid) {
      return 'Paid $' + act.amount_paid;
    }
    const today = new Date();
    const actDate = new Date(act.date);
    if (today < actDate) {
      return 'Pending';
    }
    return 'Owe $' + act.amount;
  }

  getAddress(a) {
    if (a && a.street) {
      let apt = '';
      if (a.apt !== undefined && a.apt !== '') {
        apt = '#' + a.apt + ' ';
      }
      return apt + a.street + ', ' + a.city;
    }
    return '';
  }

  getAddressLink(a) {
    if (a && a.street) {
      return (
        'https://maps.google.com/?q=' +
        encodeURIComponent(
          a.apt + ' ' + a.street + ', ' + a.city + ', ' + a.state + ' ' + a.zip
        )
      );
    }
    return '';
  }

  sublogsObserver() {
    if (this.sublogs) {
      const sublogs = this.sublogs;
      for (var i = 0; i < sublogs.length; i++) {
        sublogs[i].statusString = this._getStatusString(sublogs[i]);
        sublogs[i].addressString = this._getAddress(sublogs[i].address);
        sublogs[i].addressLink = this._getAddressLink(sublogs[i].address);
        sublogs[i].dateString = this._getDateString(sublogs[i].date);
        sublogs[i].paidDateString = this._getPaidDateString(
          sublogs[i].paid_datetime
        );
        sublogs[i].discount_string = this._getDiscountString(
          sublogs[i].discount_amount,
          sublogs[i].discount_percent
        );
      }
      this.sublogs = sublogs;
    }
  }

  getDateString(dateString) {
    const monthNames = [
      'January',
      'February',
      'March',
      'April',
      'May',
      'June',
      'July',
      'August',
      'September',
      'October',
      'November',
      'December',
    ];
    const dayNames = [
      'Sunday',
      'Monday',
      'Tuesday',
      'Wedensday',
      'Thursday',
      'Friday',
      'Saturday',
    ];
    const d = new Date(dateString.substr(0, 10) + 'T12:12:12');
    let day = d.getDay();
    let month = d.getMonth();
    let date = d.getDate();
    let year = d.getFullYear();
    return `${dayNames[day]}, ${monthNames[month]} ${date} ${year}`;
  }

  getPaidDateString(dateString) {
    const today = new Date();
    const sublogDate = new Date(dateString);
    if (today < sublogDate) {
      return ' ';
    }
    const monthNames = [
      'Jan',
      'Feb',
      'Mar',
      'Apr',
      'May',
      'Jun',
      'Jul',
      'Aug',
      'Sep',
      'Oct',
      'Nov',
      'Dec',
    ];
    const d = new Date(dateString.substr(0, 10) + 'T12:12:12');
    let day = d.getDay();
    let month = d.getMonth();
    let date = d.getDate();
    let year = d.getFullYear();
    return `${monthNames[month]} ${date} ${year}`;
  }
}
</script>

<style scoped lang="scss">
</style>
