package application

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"googlemaps.github.io/maps"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/core/account"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	mapsClient  *maps.Client
	mapsCompMap map[maps.Component]string
)

func GetApplications(ctx context.Context, user *types.User) ([]*ChefApplication, error) {
	if !user.IsAdmin() {
		utils.Errorf(ctx, "user(%v) attemted to do an admin task.", *user)
		return nil, errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not an admin."}
	}
	chefApplications, err := getAll(ctx)
	return chefApplications, err
}

// GetApplication gets a chef application
func GetApplication(ctx context.Context, user *types.User) (*ChefApplication, error) {
	var err error
	chefApplication := &ChefApplication{}
	err = get(ctx, user.ID, chefApplication)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			chefApplication := new(ChefApplication)
			chefApplication.Name = user.Name
			return chefApplication, nil
		}
		return nil, err
	}
	return chefApplication, nil
}

// SubmitApplication saves a ChefApplication.
// A token should be refreshed if this function is called.
func SubmitApplication(ctx context.Context, user *types.User, chefApplication *ChefApplication) (*ChefApplication, error) {
	var err error
	chefApplicationEntity := &ChefApplication{}
	err = get(ctx, user.ID, chefApplicationEntity)
	if err != nil && err.Error() != datastore.ErrNoSuchEntity.Error() {
		return nil, err
	}

	if err != nil && err.Error() == datastore.ErrNoSuchEntity.Error() {
		chefApplication.ApplicationProgress = 1
		if chefApplication.Address.String() != chefApplicationEntity.Address.String() {
			err = getGeopointFromAddress(ctx, &chefApplication.Address)
			if err != nil {
				return nil, err
			}
		} else {
			chefApplication.Address = chefApplicationEntity.Address
		}
	} else {
		chefApplication.ApplicationProgress = chefApplicationEntity.ApplicationProgress
	}
	chefApplication.UserID = user.ID
	chefApplication.LastUpdatedDateTime = time.Now().UTC()
	if chefApplicationEntity.CreatedDateTime.IsZero() {
		chefApplication.CreatedDateTime = time.Now().UTC()
	}
	err = put(ctx, user.ID, chefApplication)
	if err != nil {
		return nil, err
	}
	err = setAuthUserAndPerm(ctx, user, chefApplication.Name, chefApplication.Email)
	if err != nil {
		return nil, err
	}
	err = account.SaveUserInfo(ctx, user, &chefApplication.Address, chefApplication.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return chefApplication, nil
}

func setAuthUserAndPerm(ctx context.Context, user *types.User, name string, email string) error {
	userChanged := false
	if !user.IsChef() {
		user.SetChef(true)
		userChanged = true
	}
	if !user.HasAddress() {
		user.SetAddress(true)
		userChanged = true
	}
	if user.Name != name {
		user.Name = name
		userChanged = true
	}
	if user.Email != email {
		user.Email = email
		userChanged = true
	}
	if userChanged {
		return auth.SaveUser(ctx, user)
	}
	return nil
}

func getGeopointFromAddress(ctx context.Context, address *types.Address) error {
	getMapsClient(ctx)
	if mapsClient == nil {
		returnErr := errors.ErrorWithCode{
			Code:    errors.CodeInternalServerErr,
			Message: "Could not connect to Google Maps.",
		}
		return returnErr
	}
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

func init() {
	mapsCompMap = make(map[maps.Component]string, 1)
	mapsCompMap[maps.ComponentCountry] = "US"
}
