package gigachef

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, userID string, chef *Gigachef) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigachef, userID, 0, nil)
	return datastore.Get(ctx, key, chef)
}

func getMulti(ctx context.Context, ids []string) ([]Gigachef, error) {
	chefs := make([]Gigachef, len(ids))
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindGigachef, ids[i], 0, nil)
	}
	err := datastore.GetMulti(ctx, keys, chefs)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	return chefs, nil
}

func put(ctx context.Context, userID string, chef *Gigachef) error {
	var err error
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigachef, userID, 0, nil)

	_, err = datastore.Put(ctx, key, chef)
	return err
}
