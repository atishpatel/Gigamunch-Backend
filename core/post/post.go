package post

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/notification"

	"appengine"

	"golang.org/x/net/context"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	taxPercentage = 7.25
)

var (
	errUnauthorized     = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access."}
	errNotVerifiedChef  = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User is not a verified chef."}
	errNoSubMerchantID  = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have sub merchant id."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errMySQL            = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "There was a database error with the server."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
)

// Client is a post client
type Client struct {
	ctx context.Context
}

// New is used to create a new client for posts
func New(ctx context.Context) Client {
	return Client{
		ctx: ctx,
	}
}

// Resp is a post response
type Resp struct {
	ID int64
	Post
}

// PostPost posts a live post if the post is valid
// returns postID, error
func PostPost(ctx context.Context, user *types.User, post *Post) (int64, error) {
	var err error
	if !user.IsVerifiedChef() {
		return 0, errNotVerifiedChef
	}
	post.CreatedDateTime = time.Now().UTC()
	post.TaxPercentage = taxPercentage
	post.GigachefID = user.ID
	post.PricePerServing = post.ChefPricePerServing * 1.2
	if post.IsOrderNow {
		post.ReadyDateTime = post.ClosingDateTime
	}
	// get the gigachef post info
	postInfo, err := gigachef.GetPostInfo(ctx, user)
	if err != nil {
		return 0, err
	}
	if !postInfo.Address.Valid() {
		return 0, errUnauthorized.WithMessage("User does not have an address.")
	}
	if postInfo.BTSubMerchantStatus == "" {
		return 0, errNoSubMerchantID
	}
	post.GigachefAddress = postInfo.Address
	post.GigachefDelivery.Radius = postInfo.DeliveryRange
	post.BTSubMerchantID = postInfo.BTSubMerchantID
	post.GigachefDelivery.MaxDuration = 45 * 60 // 45 minutes
	// IncrementNumPost for gigachef
	postErrChan := make(chan error, 1)
	go func() {
		postErrChan <- gigachef.IncrementNumPost(ctx, user)
	}()
	// TODO add IncrementNumCreated post for Item
	// put in datastore
	postID, err := putIncomplete(ctx, post)
	// insert into sql table
	err = insertLivePost(ctx, postID, post)
	if err != nil {
		// TODO update to a transaction
		utils.Criticalf(ctx, "failed to put %d in sql database: +%v", postID, err)
		return 0, errMySQL.WithError(err).Wrap("mysql insert failed")
	}
	<-postErrChan
	if !appengine.IsDevAppServer() {
		// notify enis
		nC := notification.New(ctx)
		var photo string
		if len(post.Photos) != 0 {
			photo = post.Photos[0]
		}
		err = nC.SendSMS("6153975516", fmt.Sprintf("A new post was made by %s. \n Title: %s \n Desc: %s \n Image: %s", user.Name, post.Title, post.Description, photo))
		if err != nil {
			utils.Criticalf(ctx, "failed to notify enis about chef(%s) making a post(%s)", user.ID, postID)
		}
	}
	return postID, nil
}

// GetUserPosts gets post from a user sorted by ready time
func (c *Client) GetUserPosts(chefID string, start, end int) ([]int64, []Post, error) {
	postIDs, posts, err := getUserPosts(c.ctx, chefID, start, end)
	if err != nil {
		return nil, nil, errDatastore.WithError(err).Wrapf("failed to getUserPosts for chef(%s)", chefID)
	}
	return postIDs, posts, nil
}

// GetLivePostsIDs gets live posts sorted by ready date
// returns: []postIDs, []gigachefIDs, []distances, error
func GetLivePostsIDs(ctx context.Context, geopoint *types.GeoPoint, limit *types.Limit, radius int, readyDatetime time.Time, descending bool) ([]int64, []int64, []string, []float32, error) {
	var err error
	if !limit.Valid() {
		return nil, nil, nil, nil, errInvalidParameter.WithMessage("Limit range is not valid.").Wrapf("%d-%d is not a valid range limit", limit.Start, limit.End)
	}
	if !geopoint.Valid() {
		return nil, nil, nil, nil, errInvalidParameter.WithMessage("Geopoint is not valid.").Wrapf("%v is not a valid geopoint", geopoint)
	}
	// get list of sorted livePostIDs
	postIDs, itemIDs, gigachefIDs, distances, err := selectLivePosts(ctx, geopoint, limit, radius, readyDatetime, descending)
	if err != nil {
		utils.Criticalf(ctx, "failed to select live posts in sql database: +%v", err)
		return nil, nil, nil, nil, errors.Wrap("mysql select live post failed", err)
	}
	return postIDs, itemIDs, gigachefIDs, distances, nil
}

// GetPosts gets post form IDs
func GetPosts(ctx context.Context, postIDs []int64) ([]Post, error) {
	posts := make([]Post, len(postIDs))
	err := getMultiPost(ctx, postIDs, posts)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get multi posts(%v)", postIDs)
	}
	return posts, nil
}

// GetPost gets the post from a postID
func GetPost(ctx context.Context, postID int64) (*Post, error) {
	post := new(Post)
	err := get(ctx, postID, post)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get post(%d)", postID)
	}
	return post, nil
}

// GetClosedPosts returns posts that are closing in a minute or closed
func (c *Client) GetClosedPosts() ([]int64, []string, error) {
	return getClosedPosts(c.ctx)
}

// ClosePost removes the post from live posts
func (c *Client) ClosePost(postID int64, chefID string) (*Resp, error) {
	p := new(Post)
	err := get(c.ctx, postID, p)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get post")
	}
	err = removeLivePost(c.ctx, postID)
	if err != nil {
		return nil, errors.Wrap("failed to remove live post", err)
	}
	resp := &Resp{
		ID:   postID,
		Post: *p,
	}
	return resp, nil
}
