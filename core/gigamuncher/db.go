package gigamuncher

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, userID string, gigamuncher *Gigamuncher) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigamuncher, userID, 0, nil)
	return datastore.Get(ctx, key, gigamuncher)
}

func put(ctx context.Context, userID string, gigamuncher *Gigamuncher) error {
	var err error
	key := datastore.NewKey(ctx, kindGigamuncher, userID, 0, nil)
	_, err = datastore.Put(ctx, key, gigamuncher)
	// TODO add cache stuff
	return err
}
