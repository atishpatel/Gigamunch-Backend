package menu

import (
	"context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id int64) (*Menu, error) {
	menu := new(Menu)
	key := datastore.NewKey(ctx, kindMenu, "", id, nil)
	err := datastore.Get(ctx, key, menu)
	return menu, err
}

// getCookMenus returns a list of Menus
func getCookMenus(ctx context.Context, cookID string) ([]int64, []Menu, error) {
	query := datastore.NewQuery(kindMenu).
		Filter("CookID =", cookID).
		Limit(1000)
	var results []Menu
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

func put(ctx context.Context, id int64, menu *Menu) error {
	var err error
	key := datastore.NewKey(ctx, kindMenu, "", id, nil)
	_, err = datastore.Put(ctx, key, menu)
	return err
}
