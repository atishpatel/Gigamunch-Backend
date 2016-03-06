package application

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func getAll(ctx context.Context) ([]*ChefApplication, error) {
	query := datastore.NewQuery(kindChefApplication)
	var dst []*ChefApplication
	var err error
	_, err = query.GetAll(ctx, dst)
	if err != nil {
		return nil, err
	}
	return dst, err
}

func get(ctx context.Context, userID string, chefApplication *ChefApplication) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindChefApplication, userID, 0, nil)
	return datastore.Get(ctx, key, chefApplication)
}

func put(ctx context.Context, userID string, chefApplication *ChefApplication) error {
	var err error
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindChefApplication, userID, 0, nil)

	_, err = datastore.Put(ctx, key, chefApplication)
	return err
}
