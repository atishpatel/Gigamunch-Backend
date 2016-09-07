package item

import (
	"context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id string) (*Item, error) {
	item := new(Item)
	key := datastore.NewKey(ctx, kindItem, id, 0, nil)
	err := datastore.Get(ctx, key, item)
	return item, err
}

func put(ctx context.Context, id string, item *Item) error {
	var err error
	key := datastore.NewKey(ctx, kindItem, id, 0, nil)
	_, err = datastore.Put(ctx, key, item)
	return err
}
