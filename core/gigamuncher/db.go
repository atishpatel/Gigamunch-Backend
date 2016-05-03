package gigamuncher

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id string, muncher *Gigamuncher) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigamuncher, id, 0, nil)
	return datastore.Get(ctx, key, muncher)
}

func put(ctx context.Context, id string, muncher *Gigamuncher) error {
	var err error
	key := datastore.NewKey(ctx, kindGigamuncher, id, 0, nil)
	_, err = datastore.Put(ctx, key, muncher)
	// TODO add cache stuff
	return err
}
