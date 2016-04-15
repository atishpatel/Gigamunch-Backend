package order

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

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