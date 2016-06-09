package post

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, id int64, post *Post) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindPost, "", id, nil)
	return datastore.Get(ctx, key, post)
}

func put(ctx context.Context, id int64, post *Post) error {
	var err error
	key := datastore.NewKey(ctx, kindPost, "", id, nil)
	_, err = datastore.Put(ctx, key, post)
	return err
}

func putIncomplete(ctx context.Context, post *Post) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindPost, nil)
	key, err = datastore.Put(ctx, key, post)
	return key.IntID(), err
}

// getMultiPost gets a list of live posts from a list of postIDs
func getMultiPost(ctx context.Context, ids []int64, dst []Post) error {
	if ids == nil || len(ids) == 0 {
		return nil
	}
	if len(ids) != len(dst) {
		return fmt.Errorf("postIDs and dst slices have different length")
	}
	var err error
	// TODO add cache stuff
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindPost, "", ids[i], nil)
	}
	//TODO post in dst might be nil
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return err
	}
	return nil
}

func getUserPosts(ctx context.Context, gigachefID string, startLimit int, endLimit int) ([]int64, []Post, error) {
	offset := startLimit
	limit := endLimit - startLimit
	query := datastore.NewQuery(kindPost).
		Filter("GigachefID =", gigachefID).
		Order("-ReadyDateTime").
		Offset(offset).
		Limit(limit)
	var results []Post
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
