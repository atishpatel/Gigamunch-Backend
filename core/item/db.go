package item

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id int64, item *Item) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindItem, "", id, nil)
	return datastore.Get(ctx, key, item)
}

// getSortedItems returns a list of reviews sorted by CreatedDataTime
func getSortedItems(ctx context.Context, gigachefID string, startLimit int, endLimit int) ([]int64, []Item, error) {
	offset := startLimit
	limit := endLimit - startLimit
	query := datastore.NewQuery(kindItem).
		Filter("GigachefID =", gigachefID).
		Order("LastUsedDateTime").
		Offset(offset).
		Limit(limit)
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

func put(ctx context.Context, id int64, item *Item) error {
	var err error
	key := datastore.NewKey(ctx, kindItem, "", id, nil)
	_, err = datastore.Put(ctx, key, item)
	// TODO add cache stuff
	return err
}

func putIncomplete(ctx context.Context, item *Item) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindItem, nil)
	key, err = datastore.Put(ctx, key, item)
	// TODO add cache stuff
	return key.IntID(), err
}
