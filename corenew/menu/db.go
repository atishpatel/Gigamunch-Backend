package menu

import (
	"context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id string) (*Menu, error) {
	menu := new(Menu)
	key := datastore.NewKey(ctx, kindMenu, id, 0, nil)
	err := datastore.Get(ctx, key, menu)
	return menu, err
}

func put(ctx context.Context, id string, menu *Menu) error {
	var err error
	key := datastore.NewKey(ctx, kindMenu, id, 0, nil)
	_, err = datastore.Put(ctx, key, menu)
	return err
}
