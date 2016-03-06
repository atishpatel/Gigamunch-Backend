package post

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func putIncomplete(ctx context.Context, post *Post) (int64, error) {
	var err error
	postKey := datastore.NewIncompleteKey(ctx, kindPost, nil)
	postKey, err = datastore.Put(ctx, postKey, post)
	return postKey.IntID(), err
}

// getMultiPost gets a list of live posts from a list of postIDs
func getMultiPost(ctx context.Context, postIDs []int64, dst []Post) error {
	if postIDs == nil || len(postIDs) == 0 {
		return nil
	}
	if len(postIDs) != len(dst) {
		return fmt.Errorf("postIDs and dst slices have different length")
	}
	var err error
	// TODO add cache stuff
	postKeys := make([]*datastore.Key, len(postIDs))
	for i := range postIDs {
		postKeys[i] = datastore.NewKey(ctx, kindPost, "", postIDs[i], nil)
	}
	//TODO post in dst might be nil
	err = datastore.GetMulti(ctx, postKeys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return err
	}
	return nil
}
