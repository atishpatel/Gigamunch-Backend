package cook

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id string) (*Cook, error) {
	cook := new(Cook)
	key := datastore.NewKey(ctx, kindCook, id, 0, nil)
	err := datastore.Get(ctx, key, cook)
	cook.ID = id
	return cook, err
}

func getMulti(ctx context.Context, ids []string) ([]Cook, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids cannot be 0 for getMulti")
	}
	dst := make([]Cook, len(ids))
	var err error
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		if ids[i] == "" {
			return nil, errors.New("ids cannot contain an empty string")
		}
		keys[i] = datastore.NewKey(ctx, kindCook, ids[i], 0, nil)
	}
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	return dst, nil
}

func getBySubmerchantID(ctx context.Context, submerchantID string) (*Cook, error) {
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
	results[0].ID = keys[0].StringID()
	return &results[0], nil
}

func put(ctx context.Context, id string, cook *Cook) error {
	var err error
	cook.ID = id
	key := datastore.NewKey(ctx, kindCook, id, 0, nil)
	_, err = datastore.Put(ctx, key, cook)
	return err
}
