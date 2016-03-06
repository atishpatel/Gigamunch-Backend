package order

import (
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

var (
	errDatastore = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// GetOrderIDsAndPostInfo gets the user ids involved with the order
// returns: BasicOrderIDs, PostInfo, error
func GetOrderIDsAndPostInfo(ctx context.Context, orderID int64) (BasicOrderIDs, PostInfo, error) {
	var postInfo PostInfo
	order := new(Order)
	err := get(ctx, orderID, order)
	if err != nil {
		var basicOrderIDs BasicOrderIDs
		return basicOrderIDs, postInfo, errDatastore.WithError(err)
	}
	postInfo.ID = order.PostID
	postInfo.Title = order.PostTitle
	postInfo.PhotoURL = order.PostPhotoURL
	return order.BasicOrderIDs, postInfo, nil
}
