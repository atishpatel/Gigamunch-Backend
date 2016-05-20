package gigachef

import (
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

var (
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// UpdateAvgRating updates the average rating of a Gigachef.
// If rating is new, use 0 for oldRating.
func UpdateAvgRating(ctx context.Context, id string, oldRating int, newRating int) error {
	if id == "" || oldRating < 0 || oldRating > 5 || newRating < 1 || newRating > 5 {
		return errInvalidParameter.WithError(fmt.Errorf("id(%s) oldRating(%d) newRating(%d)", id, oldRating, newRating))
	}
	if oldRating == newRating {
		return nil
	}
	chef := new(Gigachef)
	// TODO should be transaction?
	err := get(ctx, id, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	if oldRating != 0 {
		chef.removeRating(oldRating)
	}
	chef.addRating(newRating)
	err = put(ctx, id, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put gigachef")
	}
	return nil
}

// GetRatingsAndInfo returns a list of Ratings for the passed in array of ids
func GetRatingsAndInfo(ctx context.Context, ids []string) ([]types.UserDetail, []Rating, error) {
	// TODO add not querying same ids
	chefs, err := getMulti(ctx, ids)
	if err != nil {
		return nil, nil, errDatastore.WithError(err).Wrap("cannot get multi gigachefs")
	}
	chefUserDetails := make([]types.UserDetail, len(chefs))
	ratings := make([]Rating, len(chefs))
	for i := range chefs {
		ratings[i] = chefs[i].Rating
		chefUserDetails[i] = chefs[i].UserDetail
	}
	return chefUserDetails, ratings, nil
}
