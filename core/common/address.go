package common

import (
	"fmt"
	"math"
)

// Address represents a location with GeoPoints and address
type Address struct {
	APT      string `json:"apt" datastore:",noindex"`
	Street   string `json:"street" datastore:",index"`
	City     string `json:"city" datastore:",index"`
	State    string `json:"state" datastore:",index"`
	Zip      string `json:"zip" datastore:",index"`
	Country  string `json:"country" datastore:",noindex"`
	GeoPoint        // embedded
}

func (a Address) String() string {
	return fmt.Sprintf("#%s %s, %s, %s %s, %s", a.APT, a.Street, a.City, a.State, a.Zip, a.Country)
}

// StringNoAPT returns a string of the address without the APT.
func (a Address) StringNoAPT() string {
	return fmt.Sprintf("%s, %s, %s %s, %s", a.Street, a.City, a.State, a.Zip, a.Country)
}

// GeoPoint represents a location as latitude/longitude in degrees.
type GeoPoint struct {
	Latitude  float64 `json:"latitude,string" datastore:",index"`
	Longitude float64 `json:"longitude,string" datastore:",index"`
}

// String returns Latitude,Longitude
func (g GeoPoint) String() string {
	return fmt.Sprintf("%g,%g", g.Latitude, g.Longitude)
}

// Valid returns whether a GeoPoint is within [-90, 90] latitude and [-180, 180] longitude.
func (g GeoPoint) Valid() bool {
	return -90 <= g.Latitude && g.Latitude <= 90 && -180 <= g.Longitude && g.Longitude <= 180
}

// GreatCircleDistance calculates the Haversine distance between two points in miles
func (g GeoPoint) GreatCircleDistance(p2 GeoPoint) float32 {
	dLat := (p2.Latitude - g.Latitude) * (math.Pi / 180.0)
	dLon := (p2.Longitude - g.Longitude) * (math.Pi / 180.0)

	lat1 := (g.Latitude) * (math.Pi / 180.0)
	lat2 := (p2.Latitude) * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return float32(3959.0 * c)
}

// EstimatedDuration calculates a guess of how long it will take to get from one
// point to another at 15 miles/hour in seconds.
func (g GeoPoint) EstimatedDuration(p2 GeoPoint) int64 {
	distance := g.GreatCircleDistance(p2)
	return int64(distance * 4.5 * 60)
}
