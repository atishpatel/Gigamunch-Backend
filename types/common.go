package types

import "fmt"

type AvaliableExchangeMethods int32

func (aem AvaliableExchangeMethods) Pickup() bool {
	return getKthBit(int32(aem), 1)
}

// UserDetail is the structure that is stored in the database for a chef's
// or muncher's details
type UserDetail struct {
	Name       string `json:"name" datastore:",noindex"`
	Email      string `json:"email" datastore:",noindex"`
	PhotoURL   string `json:"photo_url" datastore:",noindex"`
	ProviderID string `json:"provider_id" datastore:",noindex"`
}

// Address represents a location with GeoPoints and address
type Address struct {
	APT      string `json:"apt" datastore:",noindex"`
	Street   string `json:"street" datastore:",noindex"`
	City     string `json:"city" datastore:",noindex"`
	State    string `json:"state" datastore:",noindex"`
	Zip      string `json:"zip" datastore:",noindex"`
	Country  string `json:"country" datastore:",noindex"`
	GeoPoint        // embedded
}

func (a *Address) String() string {
	return fmt.Sprintf("%s #%s, %s, %s %s, %s", a.Street, a.APT, a.City, a.State, a.Zip, a.Country)
}

// GeoPoint represents a location as latitude/longitude in degrees.
type GeoPoint struct {
	Latitude  float32 `json:"latitude" datastore:",noindex"`
	Longitude float32 `json:"longitude" datastore:",noindex"`
}

// String returns Latitude,Longitude
func (g GeoPoint) String() string {
	return fmt.Sprintf("%g,%g", g.Latitude, g.Longitude)
}

// Valid returns whether a GeoPoint is within [-90, 90] latitude and [-180, 180] longitude.
func (g GeoPoint) Valid() bool {
	return -90 <= g.Latitude && g.Latitude <= 90 && -180 <= g.Longitude && g.Longitude <= 180
}

// Limit is the a range limit for quries
type Limit struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Valid returns whether a LimitRange's EndLimit > StartLimit
func (l Limit) Valid() bool {
	return l.Start >= 0 && l.End > 0 && l.End > l.Start
}
