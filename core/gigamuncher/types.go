package gigamuncher

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	// kindGigamuncher is used for the basic Gigamuncher account info
	kindGigamuncher = "Gigamuncher"
)

// Gigamuncher contains the basic Gigamuncher account info
type Gigamuncher struct {
	CreatedDatetime  time.Time   `json:"created_datetime" datastore:",noindex"`
	types.UserDetail             //embedded
	Addresses        []Addresses `json:"addresses" datastore:",noindex"`
	BTCustomerID     string      `json:"bt_customer_id" datastore:",index"`
}

// Addresses contains the address and lasttime the address was used
type Addresses struct {
	types.Address           // embedded
	AddedDateTime time.Time `json:"added_datetime" datastore:",noindex"`
	Selected      bool      `json:"selected" datastore:",noindex"`
}
