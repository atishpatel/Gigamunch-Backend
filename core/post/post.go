package post

import (
	"fmt"
	"time"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/item"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/notification"

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

type notifyClient interface {
	SendSMS(string, string) error
}

type chefClient interface {
	GetPostInfo(string) (*gigachef.PostInfoResp, error)
	IncrementNumPost(string) error
}

type itemClient interface {
	IncrementNumPostsCreated(int64) error
	AddNumTotalOrders(int64, int32) error
}

// PublishPostReq is the req to publish a post
type PublishPostReq struct {
	User                      *types.User    `json:"user"`
	BaseItem                  types.BaseItem `json:"base_item"`
	ItemID                    int64          `json:"item_id"`
	Title                     string         `json:"title"`
	ClosingDateTime           time.Time      `json:"closing_datetime"`
	StartPickupDateTime       time.Time      `json:"start_pickup_datetime"`
	EndPickupDateTime         time.Time      `json:"end_pickup_datetime"`
	ChefDelivery              bool           `json:"chef_delivery"`
	ChefDeliveryRadius        int32          `json:"chef_delivery_radius"`
	ChefDeliveryBasePrice     float32        `json:"chef_delivery_base_price"`
	StartChefDeliveryDateTime time.Time      `json:"start_chef_delivery_datetime"`
	EndChefDeliveryDateTime   time.Time      `json:"end_chef_delivery_datetime"`
	ServingsOffered           int32          `json:"servings_offered"`
	ChefPricePerServing       float32        `json:"chef_price_per_serving"`
}

// PublishPost publishes a post into the food feed
func (c *Client) PublishPost(req *PublishPostReq) (int64, *Post, error) {
	nC := notification.New(c.ctx)
	chefC := gigachef.New(c.ctx)
	itemC := item.New(c.ctx)
	return publishPost(c.ctx, req, chefC, itemC, nC)
}

func publishPost(ctx context.Context, req *PublishPostReq, chefC chefClient, itemC itemClient, nC notifyClient) (int64, *Post, error) {
	var err error
	if !req.User.IsVerifiedChef() {
		return 0, nil, errNotVerifiedChef
	}
	// get the chef post info
	postInfo, err := chefC.GetPostInfo(req.User.ID)
	if err != nil {
		return 0, nil, errors.Wrap("failed to chef.GetPostInfo", err)
	}
	if !postInfo.Address.Valid() {
		return 0, nil, errUnauthorized.WithMessage("User does not have an address.")
	}
	if postInfo.BTSubMerchantStatus == "" {
		return 0, nil, errNoSubMerchantID
	}
	// create post
	post := &Post{
		BaseItem:        req.BaseItem,
		ItemID:          req.ItemID,
		Title:           req.Title,
		ClosingDateTime: req.ClosingDateTime,
		GigachefDelivery: GigachefDelivery{
			Radius:    req.ChefDeliveryRadius,
			BasePrice: req.ChefDeliveryBasePrice,
		},
		ServingsOffered:     req.ServingsOffered,
		ChefPricePerServing: req.ChefPricePerServing,
		TaxPercentage:       taxPercentage,
		PricePerServing:     req.ChefPricePerServing * 1.2,
		GigachefAddress:     postInfo.Address,
		BTSubMerchantID:     postInfo.BTSubMerchantID,
	}
	post.CreatedDateTime = time.Now()
	post.GigachefID = req.User.ID
	// add pickup options
	const pickupBuffer = 30 * time.Minute
	var pickupExchange types.ExchangeMethods
	pickupExchange.SetPickup(true)
	pickupTime := req.StartPickupDateTime
	endTime := req.EndPickupDateTime.Add(-pickupBuffer + (1 * time.Minute))
	for pickupTime.Before(endTime) {
		post.ExchangeTimes = append(post.ExchangeTimes, ExchangeTimeSegment{
			StartDateTime:            pickupTime,
			EndDateTime:              pickupTime.Add(pickupBuffer),
			AvailableExchangeMethods: pickupExchange,
		})
		pickupTime = pickupTime.Add(pickupBuffer)
	}
	// add chef delivery options
	if req.ChefDelivery {
		const chefDeliveryBuffer = 1 * time.Hour
		var chefDeliveryExchange types.ExchangeMethods
		endTime := req.EndChefDeliveryDateTime.Add(-chefDeliveryBuffer + (1 * time.Minute))
		chefDeliveryExchange.SetChefDelivery(true)
		chefDeliveryTime := req.StartChefDeliveryDateTime
		for chefDeliveryTime.Before(endTime) {
			post.ExchangeTimes = append(post.ExchangeTimes, ExchangeTimeSegment{
				StartDateTime:            chefDeliveryTime,
				EndDateTime:              chefDeliveryTime.Add(chefDeliveryBuffer),
				AvailableExchangeMethods: chefDeliveryExchange,
			})
			chefDeliveryTime = chefDeliveryTime.Add(chefDeliveryBuffer)
		}
	}
	// make post
	// put in datastore
	postID, err := putIncomplete(ctx, post)
	// insert into sql table
	err = insertLivePost(ctx, postID, post)
	if err != nil {
		utils.Criticalf(ctx, "failed to put %d in sql database: +%v", postID, err)
		return 0, nil, errMySQL.WithError(err).Wrap("mysql insert failed")
	}
	// IncrementNumPost for gigachef and posts
	const numRoutines = 2
	errChan := make(chan error, numRoutines)
	go func() {
		errChan <- chefC.IncrementNumPost(req.User.ID)
	}()
	go func() {
		errChan <- itemC.IncrementNumPostsCreated(post.ItemID)
	}()
	if !appengine.IsDevAppServer() {
		// notify enis
		var photo string
		if len(post.Photos) != 0 {
			photo = post.Photos[0]
		}
		err = nC.SendSMS("6153975516", fmt.Sprintf("A new post was made. \n Title: %s \n Desc: %s \n Image: %s \n\nPostID: %s", post.Title, post.Description, photo, postID))
		if err != nil {
			utils.Criticalf(ctx, "failed to notify enis about chef(%s) making a post(%s)", req.User.ID, postID)
		}
	}
	// flush goroutines
	for i := 0; i < numRoutines; i++ {
		err = <-errChan
		if err != nil {
			utils.Errorf(ctx, "error while publishingPost: %s", err)
		}
	}
	return postID, post, nil
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
