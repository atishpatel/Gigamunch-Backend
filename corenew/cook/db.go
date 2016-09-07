package cook

import (
	"context"
	"fmt"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id string) (*Cook, error) {
	cook := new(Cook)
	key := datastore.NewKey(ctx, kindCook, id, 0, nil)
	err := datastore.Get(ctx, key, cook)
	return cook, err
}

func getBySubmerchantID(ctx context.Context, submerchantID string) (string, *Cook, error) {
	query := datastore.NewQuery(kindCook).
		Filter("BTSubMerchantID =", submerchantID)
	var results []Cook
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return "", nil, err
	}
	if len(results) != 1 {
		return "", nil, fmt.Errorf("failed to find 1 cook by submerchantID(%s): found: %v", submerchantID, results)
	}
	return keys[0].StringID(), &results[0], nil
}

func put(ctx context.Context, id string, cook *Cook) error {
	var err error
	key := datastore.NewKey(ctx, kindCook, id, 0, nil)
	_, err = datastore.Put(ctx, key, cook)
	return err
}
