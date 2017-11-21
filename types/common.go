package types

import (
	"fmt"
	"math"
)

func getKthBit(num int32, k uint32) bool {
	return (uint32(num)>>k)&1 == 1
}

func setKthBit(num int32, k uint32, x bool) int32 {
	if x {
		return int32(uint32(num) ^ ((1<<k)^uint32(num))&(1<<k))
	}
	return int32(uint32(num) ^ ((0<<k)^uint32(num))&(1<<k))
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

func (a Address) String() string {
	return fmt.Sprintf("#%s %s, %s, %s %s, %s", a.APT, a.Street, a.City, a.State, a.Zip, a.Country)
}

// StringNoAPT returns a string of the address without the APT.
func (a Address) StringNoAPT() string {
	return fmt.Sprintf("%s, %s, %s %s, %s", a.Street, a.City, a.State, a.Zip, a.Country)
}

// GeoPoint represents a location as latitude/longitude in degrees.
type GeoPoint struct {
	Latitude  float64 `json:"latitude,string" datastore:",noindex"`
	Longitude float64 `json:"longitude,string" datastore:",noindex"`
}

// String returns Latitude,Longitude
func (g GeoPoint) String() string {
	return fmt.Sprintf("%g,%g", g.Latitude, g.Longitude)
}

// Valid returns whether a GeoPoint is within [-90, 90] latitude and [-180, 180] longitude.
func (g GeoPoint) Valid() bool {
	if g.Latitude == 0 && g.Longitude == 0 {
		return false
	}
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

// Limit is the a range limit for quries
type Limit struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Valid returns whether a LimitRange's EndLimit > StartLimit
func (l Limit) Valid() bool {
	return l.Start >= 0 && l.End > 0 && l.End > l.Start
}
