package maps

import (
	"strconv"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	"google.golang.org/appengine/urlfetch"
)

const (
	metersInMile = 1609.344
)

var (
	serverKey           string
	errMapsConnect      = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Could not connect to Google Maps."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errMaps             = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with address."}
)

// GetDirections gets the optimal route for a list of points along with the total duration
// returns: arrival times for each waypoint, optimal route order based on indexes from array, error
func GetDirections(ctx context.Context, depratureTime time.Time, origin types.GeoPoint, points []types.GeoPoint) ([]time.Time, []int, error) {
	if points == nil || len(points) == 0 {
		return nil, nil, errInvalidParameter.WithMessage("Invalid waypoints")
	}
	if len(points) == 1 {
		t := origin.EstimatedDuration(points[0])
		duration := time.Duration(2*t) * time.Second
		return []time.Time{depratureTime.Add(duration)}, []int{0}, nil
	}
	mapsClient, err := getMapsClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	waypoints := []string{"optimize:true"}
	for _, v := range points {
		waypoints = append(waypoints, v.String())
	}
	req := &maps.DirectionsRequest{
		Origin:        origin.String(),
		Destination:   origin.String(),
		DepartureTime: ttos(depratureTime),
		Units:         maps.UnitsImperial,
		Waypoints:     waypoints,
	}
	routes, _, err := mapsClient.Directions(ctx, req)
	if err != nil {
		return nil, nil, errMaps.WithError(err).Wrap("cannot maps.Directions")
	}
	optimalPointsOrder := routes[0].WaypointOrder
	// get arrival times
	numLegs := len(routes[0].Legs)
	var arrivalTimes []time.Time
	for i := 0; i < numLegs-2; i++ {
		arrivalTimes = append(arrivalTimes, routes[0].Legs[i].ArrivalTime)
	}
	return arrivalTimes, optimalPointsOrder, nil
}

// GetDistance returns the distance using roads between two points.
// The points should return string "X,Y" where X and Y are floats.
// returns miles, duration, err
func GetDistance(ctx context.Context, p1, p2 types.GeoPoint) (float32, *time.Duration, error) {
	mapsClient, err := getMapsClient(ctx)
	if err != nil {
		return 0, nil, err
	}
	mapsReq := &maps.DistanceMatrixRequest{
		Origins:      []string{p1.String()},
		Destinations: []string{p2.String()},
		Units:        maps.UnitsImperial,
	}
	mapsResp, err := mapsClient.DistanceMatrix(ctx, mapsReq)
	if err != nil {
		return 0, nil, errMaps.WithError(err).Wrap("cannot get distance martrix")
	}
	element := mapsResp.Rows[0].Elements[0]
	miles := float32(element.Distance.Meters) / metersInMile // convert to miles
	if miles < .01 {
		miles = p1.GreatCircleDistance(p2)
	}
	return miles, &element.Duration, nil
}

// GetGeopointFromAddress sets Latitude and Longitude to an address.
func GetGeopointFromAddress(ctx context.Context, address *types.Address) error {
	mapsClient, err := getMapsClient(ctx)
	if err != nil {
		return err
	}
	mapsCompMap := make(map[maps.Component]string, 1)
	mapsCompMap[maps.ComponentCountry] = "US"
	mapsReq := &maps.GeocodingRequest{
		Address:    address.String(),
		Components: mapsCompMap,
	}
	mapsGeocodeResults, err := mapsClient.Geocode(ctx, mapsReq)
	if err != nil {
		return errMaps.WithError(err).Wrap("cannot get geopoint from address")
	}
	if mapsGeocodeResults[0].Geometry.LocationType != string(maps.GeocodeAccuracyRooftop) && mapsGeocodeResults[0].Geometry.LocationType != string(maps.GeocodeAccuracyRangeInterpolated) {
		return errInvalidParameter.WithMessage("Address is not valid.")
	}
	location := mapsGeocodeResults[0].Geometry.Location
	address.GeoPoint = types.GeoPoint{
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}
	return nil
}

func getMapsClient(ctx context.Context) (*maps.Client, error) {
	if serverKey == "" {
		serverKey = config.GetServerKey(ctx)
	}
	var err error
	mapsClient, err := maps.NewClient(maps.WithAPIKey(serverKey), maps.WithHTTPClient(urlfetch.Client(ctx)))
	if err != nil {
		return nil, errMapsConnect.WithError(err).Wrap("cannot get new maps client")
	}
	return mapsClient, nil
}

func ttos(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
