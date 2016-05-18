package order

import (
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func getByTransactionID(ctx context.Context, transactionID string) (int64, *Order, error) {
	query := datastore.NewQuery(kindOrder).Filter("PaymentInfo.BTTransactionID =", transactionID)
	var results []Order
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return 0, nil, err
	}
	if len(keys) != 1 || len(results) != 1 {
		return 0, nil, fmt.Errorf("transactionID(%s) has two orders", transactionID)
	}
	return keys[0].IntID(), &results[0], nil
}

// getSortedOrders returns a list of reviews sorted by CreatedDateTime
func getSortedOrders(ctx context.Context, muncherID string, startLimit int, endLimit int) ([]int64, []Order, error) {
	offset := startLimit
	limit := endLimit - startLimit
	query := datastore.NewQuery(kindOrder).
		Filter("GigamuncherID =", muncherID).
		Order("ExpectedExchangeDateTime").
		Offset(offset).
		Limit(limit)
	var results []Order
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int64, len(keys))
	for i := range keys {
		ids[i] = keys[i].IntID()
	}
	return ids, results, nil
}

func get(ctx context.Context, orderID int64, order *Order) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindOrder, "", orderID, nil)
	return datastore.Get(ctx, key, order)
}

func put(ctx context.Context, orderID int64, order *Order) error {
	var err error
	key := datastore.NewKey(ctx, kindOrder, "", orderID, nil)
	_, err = datastore.Put(ctx, key, order)
	// TODO add cache stuff
	return err
}

var putIncomplete = func(ctx context.Context, order *Order) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindOrder, nil)
	key, err = datastore.Put(ctx, key, order)
	return key.IntID(), err
}
