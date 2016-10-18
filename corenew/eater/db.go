package eater

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id string) (*Eater, error) {
	eater := new(Eater)
	key := datastore.NewKey(ctx, kindEater, id, 0, nil)
	err := datastore.Get(ctx, key, eater)
	return eater, err
}

func put(ctx context.Context, id string, eater *Eater) error {
	var err error
	key := datastore.NewKey(ctx, kindEater, id, 0, nil)
	_, err = datastore.Put(ctx, key, eater)
	return err
}
