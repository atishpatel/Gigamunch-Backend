package review

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

func get(ctx context.Context, reviewID int64, review *Review) error {
	// TODO add cache stuff
	key := datastore.NewKey(ctx, kindReview, "", reviewID, nil)
	return datastore.Get(ctx, key, review)
}

// getSortedReviews returns a list of reviews sorted by CreatedDateTime
func getSortedReviews(ctx context.Context, gigachefID string, startLimit int, endLimit int) ([]int64, []Review, error) {
	offset := startLimit
	limit := endLimit - startLimit
	query := datastore.NewQuery(kindReview).
		Filter("GigachefID =", gigachefID).
		Order("-CreatedDateTime").
		Offset(offset).
		Limit(limit)
	var results []Review
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

func put(ctx context.Context, reviewID int64, review *Review) error {
	var err error
	key := datastore.NewKey(ctx, kindReview, "", reviewID, nil)
	_, err = datastore.Put(ctx, key, review)
	// TODO add cache stuff
	return err
}

func putIncomplete(ctx context.Context, review *Review) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindReview, nil)
	key, err = datastore.Put(ctx, key, review)
	// TODO add cache stuff
	return key.IntID(), err
}
