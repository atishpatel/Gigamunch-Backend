package order

import (
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
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

// GetMulti returns orders from the array
func (c *Client) GetMulti(orderIDs []int64) ([]Order, error) {
	var orders []Order
	if len(orderIDs) == 0 {
		return orders, nil
	}
	orders, err := getMultiOrders(c.ctx, orderIDs)
	if err != nil {
		return nil, errors.Wrap("failed to getMultiOrders", err)
	}
	return orders, nil
}
