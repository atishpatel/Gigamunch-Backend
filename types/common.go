package types

// Address represents a location with GeoPoints and address
type Address struct {
	Country         string `json:"country" datastore:",noindex"`
	State           string `json:"state" datastore:",noindex"`
	City            string `json:"city" datastore:",noindex"`
	Zip             string `json:"zip" datastore:",noindex"`
	ApartmentNumber string `json:"apartment_number" datastore:",noindex"`
	GeoPoint               // embedded
}

// GeoPoint represents a location as latitude/longitude in degrees.
type GeoPoint struct {
	Latitude  float64 `json:"latitude" datastore:",noindex"`
	Longitude float64 `json:"longitude" datastore:",noindex"`
}

// Valid returns whether a GeoPoint is within [-90, 90] latitude and [-180, 180] longitude.
func (g GeoPoint) Valid() bool {
	return -90 <= g.Latitude && g.Latitude <= 90 && -180 <= g.Longitude && g.Longitude <= 180
}
