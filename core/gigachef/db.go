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

func put(ctx context.Context, userID string, gigachef *Gigachef) error {
	var err error
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindGigachef, userID, 0, nil)

	_, err = datastore.Put(ctx, key, gigachef)
	return err
}
