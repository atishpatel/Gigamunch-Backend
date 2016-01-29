package types

import "time"

// Gigamuncher contains the basic Gigamuncher account info
type Gigamuncher struct {
	UserDetail             //embedded
	Addresses  []Addresses `json:"addresses" datastore:",noindex"`
}

// Addresses contains the address and lasttime the address was used
type Addresses struct {
	Address
	LastUsedDataTime time.Time `json:"lastused_datatime" datastore:",noindex"`
}
