package gigachef

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// func SaveChefApplication(authToken )
// create gigachef type
// update user Permission

func get(ctx context.Context, userID string, gigachef *Gigachef) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigachef, userID, 0, nil)
	return datastore.Get(ctx, key, gigachef)
}

func getMulti(ctx context.Context, ids []string) ([]Gigachef, error) {
	gigachefs := make([]Gigachef, len(ids))
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindGigachef, ids[i], 0, nil)
	}
	err := datastore.GetMulti(ctx, keys, gigachefs)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	return gigachefs, nil
}

func put(ctx context.Context, userID string, gigachef *Gigachef) error {
	var err error
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigachef, userID, 0, nil)

	_, err = datastore.Put(ctx, key, gigachef)
	return err
}
