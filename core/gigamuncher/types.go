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
}

// Addresses contains the address and lasttime the address was used
type Addresses struct {
	types.Address              // embedded
	LastUsedDataTime time.Time `json:"lastused_datatime" datastore:",noindex"`
}