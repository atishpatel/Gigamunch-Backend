package post

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	taxPercentage = 9.25
)

var (
	// errNotVerifiedChef is an error for when unverfied chefs try and unauthorized action
	errNotVerifiedChef = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not a verified chef."}
	errDatastore       = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// PostPost posts a live post if the post is valid
// returns postID, error
func PostPost(ctx context.Context, user *types.User, post *Post) (int64, error) { // TODO: update to pass in itemid for stats stuff
	var err error
	if !user.IsVerifiedChef() {
		return 0, errNotVerifiedChef
	}
	// TODO check if chef has sub merchant info stuff
	post.CreatedDateTime = time.Now().UTC()
	post.TaxPercentage = taxPercentage
	post.GigachefID = user.ID
	// get the gigachef address
	var address *types.Address
	var radius int
	address, radius, err = gigachef.GetDeliveryInfo(ctx, user)
	if err != nil {
		return 0, err
	}
	post.GigachefAddress = *address
	post.GigachefDeliveryRadius = radius
	// IncrementNumPost for gigachef
	errChan := make(chan error, 1)
	go func() {
		errChan <- gigachef.IncrementNumPost(ctx, user)
	}()
	// put in datastore
	postID, err := putIncomplete(ctx, post)
	// insert into sql table
	err = insertLivePost(postID, post)
	if err != nil {
		// TODO update to a transaction
		utils.Criticalf(ctx, "failed to put %d in sql database: +%v", postID, err)
		return 0, errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "mysql insert failed"}.WithError(err)
	}
	<-errChan
	return postID, nil
}

// GetUserPosts gets post from a user sorted by closing time
func GetUserPosts(ctx context.Context, user types.User, limit *types.Limit) ([]int64, []Post, error) {
	postIDs, posts, err := getUserPosts(ctx, user.ID, limit.Start, limit.End)
	if err != nil {
		return nil, nil, errDatastore.WithError(err)
	}
	return postIDs, posts, nil
}

// GetLivePostsIDs gets live posts sorted by ready date
// returns: []postIDs, []gigachefIDs, []distances, error
func GetLivePostsIDs(ctx context.Context, geopoint *types.GeoPoint, limit *types.Limit, radius int, readyDatetime time.Time, descending bool) ([]int64, []string, []float32, error) {
	var err error
	if !limit.Valid() {
		err = fmt.Errorf("%d-%d is not a valid range limit", limit.Start, limit.End)
		return nil, nil, nil, errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Limit range is not valid."}.WithError(err)
	}
	if !geopoint.Valid() {
		err = fmt.Errorf("%v is not a valid geopoint", geopoint)
		return nil, nil, nil, errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Geopoint is not valid."}.WithError(err)
	}
	// get list of sorted livePostIDs
	var postIDs []int64
	var distances []float32
	var gigachefIDs []string
	postIDs, gigachefIDs, distances, err = selectLivePosts(ctx, geopoint, limit, radius, readyDatetime, descending)
	if err != nil {
		utils.Criticalf(ctx, "failed to select live posts in sql database: +%v", err)
		return nil, nil, nil, errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "mysql select failed"}.WithError(err)
	}
	return postIDs, gigachefIDs, distances, nil
}

// GetPosts gets post form IDs
func GetPosts(ctx context.Context, postIDs []int64) ([]Post, error) {
	posts := make([]Post, len(postIDs))
	err := getMultiPost(ctx, postIDs, posts)
	if err != nil {
		return nil, errDatastore.WithError(err)
	}
	return posts, nil
}

// GetPost gets the post from a postID
func GetPost(ctx context.Context, postID int64) (*Post, error) {
	post := new(Post)
	err := get(ctx, postID, post)
	if err != nil {
		return nil, errDatastore.WithError(err)
	}
	return post, nil
}
