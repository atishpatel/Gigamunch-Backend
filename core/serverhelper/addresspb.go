package serverhelper

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// AddressFromPB Address From PB.
func AddressFromPB(ctx context.Context, in *pbcommon.Address) (*common.Address, error) {
	if in == nil {
		return nil, errors.BadRequestError.WithMessage("Address is not valid.")
	}
	geopoint := common.GeoPoint{Latitude: in.Latitude, Longitude: in.Longitude}
	if geopoint.Valid() && geopoint.Latitude != 0 && geopoint.Longitude != 0 {
		return &common.Address{
			APT:      in.Apt,
			Street:   in.Street,
			City:     in.City,
			State:    in.State,
			Zip:      in.Zip,
			Country:  in.Country,
			GeoPoint: geopoint,
		}, nil
	}
	addressString := in.FullAddress
	if addressString == "" {
		addressString = fmt.Sprintf(" %s, %s, %s %s, %s", in.Street, in.City, in.State, in.Zip, in.Country)
	}
	return maps.GetAddress(ctx, addressString, in.Apt)
}
