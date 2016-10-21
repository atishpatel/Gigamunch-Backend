package menu

import (
	"errors"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id int64) (*Menu, error) {
	menu := new(Menu)
	key := datastore.NewKey(ctx, kindMenu, "", id, nil)
	err := datastore.Get(ctx, key, menu)
	menu.ID = id
	return menu, err
}

func getMulti(ctx context.Context, ids []int64) ([]Menu, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids cannot be 0 for getMulti")
	}
	dst := make([]Menu, len(ids))
	var err error
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindMenu, "", ids[i], nil)
	}
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	for i := range dst {
		dst[i].ID = ids[i]
	}
	return dst, nil
}

// getCookMenus returns a list of Menus
func getCookMenus(ctx context.Context, cookID string) ([]Menu, error) {
	query := datastore.NewQuery(kindMenu).
		Filter("CookID =", cookID).
		Limit(1000)
	var results []Menu
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		results[i].ID = keys[i].IntID()
	}
	return results, nil
}

func put(ctx context.Context, id int64, menu *Menu) error {
	var err error
	menu.ID = id
	key := datastore.NewKey(ctx, kindMenu, "", id, nil)
	_, err = datastore.Put(ctx, key, menu)
	return err
}

func putIncomplete(ctx context.Context, menu *Menu) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindMenu, nil)
	key, err = datastore.Put(ctx, key, menu)
	menu.ID = key.IntID()
	return key.IntID(), err
}
