package gigachef

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/errors"
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
	gigachef := new(Gigachef)
	// TODO should be transaction?
	err := get(ctx, userID, gigachef)
	if err != nil {
		return err
	}
	if oldRating != 0 {
		gigachef.RemoveRating(oldRating)
	}
	gigachef.AddRating(newRating)
	err = put(ctx, userID, gigachef)
	if err != nil {
		return errDatastore.WithError(err)
	}
	return nil
}

// GetRatings returns a list of GigachefRatings for the passed in array of ids
func GetRatings(ctx context.Context, gigachefIDs []string) ([]float32, error) {
	// TODO add not querying same ids
	gigachefs, err := getMulti(ctx, gigachefIDs)
	if err != nil {
		return nil, errDatastore.WithError(err)
	}
	ratings := make([]float32, len(gigachefs))
	for i := range gigachefs {
		ratings[i] = gigachefs[i].AverageRating
	}
	return ratings, nil
}
