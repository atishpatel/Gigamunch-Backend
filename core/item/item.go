package item

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

var (
	// errNotVerifiedChef is an error for when unverfied chefs try and unauthorized action
	errNotChef            = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not a chef."}
	errDatastore          = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errInvalidParameter   = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errUnauthorizedAccess = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have permission to update item."}
)

// SaveItem saves a item. If ItemID is 0, a new item is created.
// returns ItemID, error
func SaveItem(ctx context.Context, user *types.User, itemID int64, item *Item) (int64, error) {
	if item == nil {
		return 0, errInvalidParameter.WithError(fmt.Errorf("Item is nil."))
	}
	if !user.IsChef() {
		return 0, errNotChef
	}
	var err error
	// err = item.Valid()
	item.GigachefID = user.ID
	item.LastUsedDateTime = time.Now().UTC()
	if itemID == 0 {
		// create a new item
		item.CreatedDateTime = time.Now().UTC()
		itemID, err = putIncomplete(ctx, item)
	} else {
		// update item
		oldItem := new(Item)
		err = get(ctx, itemID, oldItem)
		if err != nil {
			return 0, errDatastore.WithError(err)
		}
		if oldItem.GigachefID != user.ID {
			return 0, errUnauthorizedAccess
		}
		item.NumPostsCreated = oldItem.NumPostsCreated
		item.NumTotalOrders = oldItem.NumTotalOrders
		item.AverageItemRating = oldItem.AverageItemRating
		item.NumRatings = oldItem.NumRatings
		err = put(ctx, itemID, item)
		if err != nil {
			return 0, errDatastore.WithError(err)
		}
	}
	return itemID, nil
}

// GetItem gets an item if the user has access to it
func GetItem(ctx context.Context, user *types.User, itemID int64) (*Item, error) {
	item := new(Item)
	err := get(ctx, itemID, item)
	if err != nil {
		return nil, errDatastore.WithError(err)
	}
	if item.GigachefID != user.ID {
		return nil, errUnauthorizedAccess
	}
	return item, nil
}

// GetItems returns an array of items sorted by LastUsedDataTime
// returns: []ids, []items, error
func GetItems(ctx context.Context, user *types.User, limit *types.Limit) ([]int64, []Item, error) {
	ids, items, err := getSortedItems(ctx, user.ID, limit.Start, limit.End)
	if err != nil {
		return nil, nil, errDatastore.WithError(err)
	}
	return ids, items, nil
}
