package eater

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const kindEater = "eater"

// Eater is someone who uses the eater app for Gigamunch.
type Eater struct {
	ID                 string    `json:"id" datastore:",noindex"`
	CreatedDatetime    time.Time `json:"created_datetime" datastore:",noindex"`
	types.UserDetail             // embedded
	Addresses          []Address `json:"addresses" datastore:",noindex"`
	BTCustomerID       string    `json:"bt_customer_id" datastore:",index"`
	NotificationTokens []string  `json:"notification_tokens" datastore:",index"`
}

// Address contains an address and when it was added.
type Address struct {
	types.Address           // embedded
	AddedDateTime time.Time `json:"added_datetime" datastore:",noindex"`
	Selected      bool      `json:"selected" datastore:",noindex"`
}
