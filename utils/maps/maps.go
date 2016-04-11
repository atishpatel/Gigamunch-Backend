package maps

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	"google.golang.org/appengine/urlfetch"
)

const (
	metersInMile = 1609.344
)

var (
	mapsClient     *maps.Client
	mapsConnectErr = errors.ErrorWithCode{
		Code:    errors.CodeInternalServerErr,
		Message: "Could not connect to Google Maps.",
	}
	errMaps = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with Google Maps."}
)

// GetDistance returns the distance using roads between two points.
// The points should return string "X,Y" where X and Y are floats.
// returns miles, duration, err
func GetDistance(ctx context.Context, p1, p2 fmt.Stringer) (float32, *time.Duration, error) {
	getMapsClient(ctx)
	if mapsClient == nil {
		return 0, nil, mapsConnectErr
	}
	mapsReq := &maps.DistanceMatrixRequest{
		Origins:      []string{p1.String()},
		Destinations: []string{p2.String()},
		Units:        maps.UnitsImperial,
	}
	mapsResp, err := mapsClient.DistanceMatrix(ctx, mapsReq)
	if err != nil {
		return 0, nil, errMaps.WithError(err)
	}
	element := mapsResp.Rows[0].Elements[0]
	miles := float32(element.Distance.Meters) / metersInMile // convert to miles
	return miles, &element.Duration, nil
}

// GetGeopointFromAddress sets Latitude and Longitude to an address.
func GetGeopointFromAddress(ctx context.Context, address *types.Address) error {
	getMapsClient(ctx)
	if mapsClient == nil {
		return mapsConnectErr
	}
	mapsCompMap := make(map[maps.Component]string, 1)
	mapsCompMap[maps.ComponentCountry] = "US"
	mapsReq := &maps.GeocodingRequest{
		Address:    address.String(),
		Components: mapsCompMap,
	}
	mapsGeocodeResults, err := mapsClient.Geocode(ctx, mapsReq)
	if err != nil {
		returnErr := errors.ErrorWithCode{
			Code:    errors.CodeInternalServerErr,
			Message: "Error getting geopoint from address.",
		}
		return returnErr.WithError(err)
	}
	if mapsGeocodeResults[0].Geometry.LocationType != string(maps.GeocodeAccuracyRooftop) && mapsGeocodeResults[0].Geometry.LocationType != string(maps.GeocodeAccuracyRangeInterpolated) {
		returnErr := errors.ErrorWithCode{
			Code:    errors.CodeInvalidParameter,
			Message: "Address is not valid.",
		}
		return returnErr
	}
	location := mapsGeocodeResults[0].Geometry.Location
	address.GeoPoint = types.GeoPoint{
		Latitude:  float32(location.Lat),
		Longitude: float32(location.Lng),
	}
	return nil
}

func getMapsClient(ctx context.Context) {
	if mapsClient == nil {
		var err error
		mapsClient, err = maps.NewClient(maps.WithAPIKey(config.GetServerKey(ctx)), maps.WithHTTPClient(urlfetch.Client(ctx)))
		if err != nil {
			utils.Errorf(ctx, "failed to get maps client: %+v", err)
		}
	}
}
