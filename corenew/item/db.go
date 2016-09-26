package item

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id int64) (*Item, error) {
	item := new(Item)
	key := datastore.NewKey(ctx, kindItem, "", id, nil)
	err := datastore.Get(ctx, key, item)
	return item, err
}

func put(ctx context.Context, id int64, item *Item) error {
	var err error
	key := datastore.NewKey(ctx, kindItem, "", id, nil)
	_, err = datastore.Put(ctx, key, item)
	return err
}

func putIncomplete(ctx context.Context, item *Item) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindItem, nil)
	key, err = datastore.Put(ctx, key, item)
	return key.IntID(), err
}

// getCookItems returns a list of Items by ordered by MenuID
func getCookItems(ctx context.Context, cookID string) ([]int64, []Item, error) {
	query := datastore.NewQuery(kindItem).
		Filter("CookID =", cookID).
		Order("MenuID").
		Limit(1000)
	var results []Item
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int64, len(keys))
	for i := range keys {
		ids[i] = keys[i].IntID()
	}
	return ids, results, nil
}

// getMulti gets a list of Items
func getMulti(ctx context.Context, ids []int64) ([]Item, error) {
	if len(ids) == 0 {
		return nil, errInvalidParameter.Wrap("ids cannot be 0 for getMulti")
	}
	dst := make([]Item, len(ids))
	var err error
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindItem, "", ids[i], nil)
	}
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	return dst, nil
}
