package review

import (
	"time"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/order"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"

	"appengine/datastore"
)

var (
	errUnauthorizedAccess = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access to this review."}
	errDatastore          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// Client is a client for Reviews
type Client struct {
	ctx context.Context
}

// New is a new client for reviews
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Resp is a id and review
type Resp struct {
	ID     int64
	Review // embedded
}

// PostReview creates or updates a review
func (c *Client) PostReview(user *types.User, reviewID int64, rating int, ratingText string, orderID int64) (*Resp, error) {
	review := new(Review)
	var err error
	isNewReview := reviewID == 0
	var orderC *order.Client
	if isNewReview { // new review
		// check if user was the one who made the order
		orderC = order.New(c.ctx)
		var order *order.Resp
		order, err = orderC.GetOrder(user.ID, orderID)
		if err != nil {
			return nil, errors.Wrap("cannot get order", err)
		}
		if user.ID != order.GigamuncherID {
			return nil, errUnauthorizedAccess
		}
		review.GigachefID = order.GigachefID
		review.GigamuncherID = user.ID
		review.GigamuncherName = user.Name
		review.GigamuncherPhotoURL = user.PhotoURL
		review.OrderID = orderID
		review.ItemID = order.ItemID
		review.Post.ID = order.PostID
		review.Post.Title = order.PostTitle
		review.Post.PhotoURL = order.PostPhotoURL
		review.CreatedDateTime = time.Now().UTC()
	} else { // update review
		// check if the user has the right to update the review
		err = get(c.ctx, reviewID, review)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return nil, errUnauthorizedAccess
			}
			return nil, errDatastore
		}
		if review.GigamuncherID != user.ID {
			return nil, errUnauthorizedAccess
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
		errChan <- gigachef.UpdateAvgRating(c.ctx, review.GigachefID, oldRating, rating)
	}()
	// TODO update avg item review
	// update review
	if isNewReview {
		reviewID, err = putIncomplete(c.ctx, review)

	} else {
		err = put(c.ctx, reviewID, review)
	}
	if err != nil {
		return nil, errDatastore.WithError(err)
	}
	if isNewReview {
		_, err = orderC.UpdateReviewID(user.ID, orderID, reviewID)
		if err != nil {
			return nil, errors.Wrap("failed to update review id for order", err)
		}
	}
	err = <-errChan
	if err != nil {
		return nil, err
	}

	resp := &Resp{
		ID:     reviewID,
		Review: *review,
	}
	return resp, nil
}

// GetReviews gets the reviews for an item
func (c *Client) GetReviews(gigachefID string, limit *types.Limit, itemID int64) ([]Resp, error) {
	ids, reviews, err := getSortedReviews(c.ctx, gigachefID, limit.Start, limit.End)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get sorted reviews")
	}
	if itemID != 0 {
		// TODO resort reviews
	}
	resps := make([]Resp, len(ids))
	for i := range ids {
		resps[i].ID = ids[i]
		resps[i].Review = reviews[i]
	}
	return resps, nil
}

// GetReview get a review
func (c *Client) GetReview(reviewID int64) (*Resp, error) {
	resp := new(Resp)
	if reviewID == 0 {
		return resp, nil
	}
	r := new(Review)
	err := get(c.ctx, reviewID, r)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get review")
	}
	resp.ID = reviewID
	resp.Review = *r
	return resp, nil
}
