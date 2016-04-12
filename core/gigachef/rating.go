package gigachef

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

var (
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// UpdateAvgRating updates the average rating of a Gigachef.
// If rating is new, use 0 for oldRating.
func UpdateAvgRating(ctx context.Context, userID string, oldRating int, newRating int) error {
	if userID == "" || oldRating < 0 || oldRating > 5 || newRating < 1 || newRating > 5 {
		return errInvalidParameter.WithError(fmt.Errorf("userID(%s) oldRating(%d) newRating(%d)", userID, oldRating, newRating))
	}
	if oldRating == newRating {
		return nil
	}
	chef := new(Gigachef)
	// TODO should be transaction?
	err := get(ctx, userID, chef)
	if err != nil {
		return err
	}
	if oldRating != 0 {
		chef.RemoveRating(oldRating)
	}
	chef.AddRating(newRating)
	err = put(ctx, userID, chef)
	if err != nil {
		return errDatastore.WithError(err)
	}
	return nil
}

// GetRatingsAndInfo returns a list of GigachefRatings for the passed in array of ids
func GetRatingsAndInfo(ctx context.Context, chefIDs []string) ([]types.UserDetail, []GigachefRating, error) {
	// TODO add not querying same ids
	chefs, err := getMulti(ctx, chefIDs)
	if err != nil {
		return nil, nil, errDatastore.WithError(err)
	}
	chefUserDetails := make([]types.UserDetail, len(chefs))
	ratings := make([]GigachefRating, len(chefs))
	for i := range chefs {
		ratings[i] = chefs[i].GigachefRating
		chefUserDetails[i] = chefs[i].UserDetail
	}
	return chefUserDetails, ratings, nil
}
