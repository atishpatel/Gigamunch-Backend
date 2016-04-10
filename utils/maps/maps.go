package maps

import (
	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	"google.golang.org/appengine/urlfetch"
)

var (
	mapsClient *maps.Client
)

func GetGeopointFromAddress(ctx context.Context, address *types.Address) error {
	getMapsClient(ctx)
	if mapsClient == nil {
		returnErr := errors.ErrorWithCode{
			Code:    errors.CodeInternalServerErr,
			Message: "Could not connect to Google Maps.",
		}
		return returnErr
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
