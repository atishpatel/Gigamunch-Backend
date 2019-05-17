<template>
  <div>
    <div class="summary">
      <div>
        <div class="name-container">
          <div class="subscriber-name">{{sub.namesString}}</div>
          <div class="edit-button">
            <a
              target="_blank"
              :href="computedDatastoreLink"
            > <i
                class="material-icons"
                style="font-size: 42px;"
              >edit</i>
            </a>
          </div>
        </div>
        <div :class="[{ 'sub-active': sub.active, 'sub-deactive': !sub.active},'subscriber-status']">
          {{computedSubStatus}}
        </div>

        <!-- Subsriber Info Table -->
        <div class="subscriber-table-info">
          <div class="info-row">
            <div class="info-label">Email:</div>
            <div class="info-value subscriber-email">
              {{sub.emailsString}}
            </div>
          </div>

          <div class="info-row">
            <div class="info-label">Phone number:</div>
            <div class="info-value subscriber-phone-number">
              {{sub.phonenumbersString}}
            </div>
          </div>

          <div class="info-row">
            <div class="info-label">Address:</div>
            <div class="info-value subscriber-address">
              <a
                target="_blank"
                :href="sub.addressLink"
              >{{sub.addressString}}</a>
            </div>
          </div>

          <div class="info-row">
            <div class="info-label">Delivery Tip:</div>
            <div class="info-value subscriber-delivery-tip">
              {{sub.delivery_notes}}
            </div>
          </div>

          <div class="info-row">
            <div class="info-label">Servings:</div>
            <div class="info-value servings">
              {{computedServings}}
            </div>
          </div>
          <div class="info-row">
            <div class="info-label">Customer ID:</div>
            <div class="info-value customer-id">
              <a
                target="_blank"
                :href="'https://www.braintreegateway.com/merchants/wsgmypp8c46cnbpc/customers/'+sub.payment_customer_id"
              >{{sub.payment_customer_id}}</a>
            </div>
          </div>
          <div class="info-row">
            <div class="info-label">Subscription Date:</div>
            <div class="info-value subscription-date">
              {{sub.signUpDatetimeTimestamp}}
            </div>
          </div>
          <div class="info-row">
            <div class="info-label">Subscription Day:</div>
            <div class="info-value subscription-day">
              {{sub.plan_weekday}}
            </div>
          </div>
          <div class="info-row">
            <div class="info-label">Deactivate Date:</div>
            <div class="info-value subscription-day">
              {{sub.deactivatedDatetimeTimestamp}}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { IsProd } from '../ts/env';

@Component({
  components: {},
})
export default class SubscriberSummary extends Vue {
  @Prop()
  public sub!: Types.SubscriberExtended;
  public logs!: Common.Log[];

  get computedServings() {
    let v = '';
    if (this.sub.servings_non_vegetarian > 0) {
      v += `${this.sub.servings_non_vegetarian} non-veg 🍖`;
      if (this.sub.servings_vegetarian > 0) {
        v += ' & ';
      }
    }
    if (this.sub.servings_vegetarian > 0) {
      v += `${this.sub.servings_vegetarian} vegetarian 🌱`;
    }
    return v;
  }

  get computedSubStatus() {
    if (this.sub.active) {
      return '• Active';
    }
    return '• Deactived';
  }

  get computedDatastoreLink() {
    if (!this.sub || !this.sub.email_prefs) {
      return '';
    }
    let project = 'gigamunch-omninexus';
    if (!IsProd()) {
      project += '-dev';
    }
    return `https://console.cloud.google.com/datastore/entities;kind=Subscriber;ns=__$DEFAULT$__/query/kind;filter=%5B%2216%2FEmailPrefs.Email%7CSTR%7CEQ%7C26%2F${
      this.sub.email_prefs[0].email
    }%22%5D?project=${project}`;
  }
}
</script>

<style scoped lang="scss">
.summary {
  background: white;
  padding-bottom: 24px;
}

.summary > div {
  margin: auto;
  max-width: 1000px;
}

.subscriber-table-info {
  display: flex;
  flex-direction: column;
}

.subscriber-name {
  font-size: 3em;
  font-weight: 600;
}

.subscriber-status {
  margin-bottom: 15px;
}

.sub-deactive {
  color: #a1a1a1;
}

.sub-active {
  color: #26cc6e;
}

.edit-button {
  font-size: 3em;
  font-weight: 600;
}

.subscriber-email {
  flex: 1;
  font-weight: 300;
}

.subscriber-phone-number {
  flex: 1;
  font-weight: 300;
}

.subscriber-address {
  flex: 1;
  font-weight: 300;
}

.subscriber-delivery-tip {
  flex: 1;
  font-weight: 300;
}

.name-container {
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-between;
}

.info-row {
  display: flex;
  flex-direction: row;
  background-color: rgb(243, 243, 243);
  border: 1px solid #d6d6d6;
  padding: 6px 10px;
}

.info-label {
  min-width: 200px;
  font-weight: 600;
}

.info-value {
  min-width: 200px;
  font-weight: 400;
}
</style>