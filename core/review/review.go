package review

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/core/order"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"

	"appengine/datastore"
)

var (
	errUnauthorizedAccess = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access to this review."}
	errDatastore          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// PostReview creates or updates a review
// returns: reviewID, error
func PostReview(ctx context.Context, user *types.User, reviewID int64, rating int, ratingText string, orderID int64) (int64, error) {
	review := new(Review)
	var err error
	isNewReview := reviewID == 0
	if isNewReview { // new review
		// check if user was the one who made the order
		orderIDs, postInfo, err := order.GetOrderIDsAndPostInfo(ctx, orderID)
		if err != nil {
			return 0, err
		}
		if user.ID != orderIDs.GigamuncherID {
			return 0, errUnauthorizedAccess
		}
		review.GigachefID = orderIDs.GigachefID
		review.GigamuncherID = user.ID
		review.OrderID = orderID
		review.Post.ID = postInfo.ID
		review.Post.Title = postInfo.Title
		review.Post.PhotoURL = postInfo.PhotoURL
		review.CreatedDataTime = time.Now().UTC()
	} else { // update review
		// check if the user has the right to update the review
		err = get(ctx, reviewID, review)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return 0, errUnauthorizedAccess
			}
			return 0, errDatastore
		}
		if review.GigamuncherID != user.ID {
			return 0, errUnauthorizedAccess
		}
		review.IsEdited = true
		review.EditedDateTime = time.Now().UTC()
	}
	if ratingText != "" {
		review.Text = ratingText
	}
	oldRating := review.Rating
	review.Rating = rating
	errChan := make(chan error, 1)
	go func() {
		// update chef avg rating
		errChan <- gigachef.UpdateAvgRating(ctx, review.GigachefID, oldRating, rating)
	}()
	// TODO update avg item review
	// update review
	if isNewReview {
		reviewID, err = putIncomplete(ctx, review)
	} else {
		err = put(ctx, reviewID, review)
	}
	if err != nil {
		return 0, errDatastore.WithError(err)
	}
	err = <-errChan
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}

// GetReviews gets
func GetReviews(ctx context.Context, gigachefID string, limit *types.Limit, itemID int64) ([]int64, []Review, error) {
	ids, reviews, err := getSortedReviews(ctx, gigachefID, limit.Start, limit.End)
	if err != nil {
		return nil, nil, errDatastore.WithError(err)
	}
	if itemID != 0 {
		// TODO resort reviews
	}
	return ids, reviews, nil
}
