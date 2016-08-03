package gigamuncher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/utils"

	"gitlab.com/atishpatel/Gigamunch-Backend/auth"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/like"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/maps"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/post"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/review"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// BasePost is the basic stuff for a post
type BasePost struct {
	ID                string   `json:"id,omitempty"`
	ID64              int64    `json:"-"`
	ItemID            string   `json:"item_id,omitempty"`
	ItemID64          int64    `json:"-"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	PricePerServing   float32  `json:"price_per_serving"`
	ServingsOffered   int32    `json:"servings_offered"`
	ServingsLeft      int32    `json:"servings_left"`
	Photos            []string `json:"photos,omitempty"`
	PostedDateTime    int      `json:"posted_datetime"`
	ClosingDateTime   int      `json:"closing_datetime"`
	ReadyDateTime     int      `json:"ready_datetime"` // TODO: REMOVE
	PickupAvaliable   bool     `json:"pickup_avaliable"`
	PickupStartTime   int      `json:"pickup_start_time"`
	PickupEndTime     int      `json:"pickup_end_time"`
	DeliveryAvailable bool     `json:"delivery_avaliable"`
	DeliveryStartTime int      `json:"delivery_start_time"`
	DeliveryEndTime   int      `json:"delivery_end_time"`
	Distance          float32  `json:"distance"`
	NumServingsBought int32    `json:"num_servings_bought"`
	NumLikes          int      `json:"num_likes"`
	HasLiked          bool     `json:"has_liked"`
}

// set takes a post.Post package and converts it to a endpoint post
func (p *BasePost) set(userID string, id int64, numLikes int, hasLiked bool, distance float32, post *post.Post) {
	p.ID64 = id
	p.ID = itos(id)
	p.ItemID64 = post.ItemID
	p.ItemID = itos(post.ItemID)
	p.Title = post.Title
	p.Description = post.Description
	p.PricePerServing = post.PricePerServing
	p.ServingsOffered = post.ServingsOffered
	p.ServingsLeft = post.ServingsOffered - post.NumServingsOrdered
	p.Photos = post.Photos
	p.PostedDateTime = ttoi(post.CreatedDateTime)
	p.ClosingDateTime = ttoi(post.ClosingDateTime)
	var pickupStartTime, pickupEndTime, deliveryStartTime, deliveryEndTime time.Time
	for _, exchange := range post.ExchangeTimes {
		if exchange.AvailableExchangeMethods.Pickup() {
			p.PickupAvaliable = true
			if pickupStartTime.IsZero() || exchange.StartDateTime.Before(pickupStartTime) {
				pickupStartTime = exchange.StartDateTime
			}
			if pickupEndTime.IsZero() || exchange.EndDateTime.After(pickupEndTime) {
				pickupEndTime = exchange.EndDateTime
			}
		}
		if exchange.AvailableExchangeMethods.Delivery() {
			p.DeliveryAvailable = true
			if deliveryStartTime.IsZero() || exchange.StartDateTime.Before(deliveryStartTime) {
				deliveryStartTime = exchange.StartDateTime
			}
			if deliveryEndTime.IsZero() || exchange.EndDateTime.After(deliveryEndTime) {
				deliveryEndTime = exchange.EndDateTime
			}
		}
	}
	p.ReadyDateTime = ttoi(pickupStartTime)
	p.PickupStartTime = ttoi(pickupStartTime)
	p.PickupEndTime = ttoi(pickupEndTime)
	p.DeliveryStartTime = ttoi(deliveryStartTime)
	p.DeliveryEndTime = ttoi(deliveryEndTime)
	p.Distance = distance
	if userID != "" {
		for _, o := range post.Orders {
			if o.GigamuncherID == userID {
				p.NumServingsBought = o.Servings
			}
		}
	}
	p.NumLikes = numLikes
	p.HasLiked = hasLiked
}

// PostGigachef is the basic version of GigachefDetails for post in live feed
type PostGigachef struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	PhotoURL   string  `json:"photo_url"`
	Rating     float32 `json:"rating"`
	NumRatings int     `json:"num_ratings"`
}

// Post is a meal that is no longer live
type Post struct {
	BasePost              // embedded
	Gigachef PostGigachef `json:"gigachef"`
}

// set takes a post.Post package and converts it to a endpoint post
func (p *Post) set(userID string, id int64, numLikes int, hasLiked bool, distance float32, rating *gigachef.Rating, gigachefName, gigachefPhotoURL string, post *post.Post) {
	p.BasePost.set(userID, id, numLikes, hasLiked, distance, post)
	p.Gigachef = PostGigachef{
		ID:         post.GigachefID,
		Name:       gigachefName,
		PhotoURL:   gigachefPhotoURL,
		Rating:     rating.AverageRating,
		NumRatings: rating.NumRatings,
	}
}

// GetLivePostsReq is the input required to get a list of live posts
type GetLivePostsReq struct {
	StartLimit    int     `json:"start_limit"`
	EndLimit      int     `json:"end_limit"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Radius        int     `json:"radius"`
	ReadyDateTime int     `json:"ready_datetime"`
	Descending    bool    `json:"descending"`
	Gigatoken     string  `json:"gigatoken"`
}

// valid returns an error if input in invalid
func (req *GetLivePostsReq) valid() error {
	if req.StartLimit < 0 || req.EndLimit < 0 {
		return fmt.Errorf("Limit is out of range.")
	}
	if req.EndLimit <= req.StartLimit {
		return fmt.Errorf("EndLimit cannot be less than or equal to StartLimit.")
	}
	point := types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	if !point.Valid() {
		return fmt.Errorf("Location inputed is not valid")
	}
	return nil
}

// GetLivePostsResp is the response for getting live posts
// returns: posts, error
type GetLivePostsResp struct {
	Posts []Post               `json:"posts,omitempty"`
	Err   errors.ErrorWithCode `json:"err"`
}

// GetLivePosts is an endpoint that returns a list of live posts
func (service *Service) GetLivePosts(ctx context.Context, req *GetLivePostsReq) (*GetLivePostsResp, error) {
	resp := new(GetLivePostsResp)
	defer handleResp(ctx, "GetLivePost", resp.Err)
	var err error
	err = req.valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	// get user
	var userID string
	if req.Gigatoken != "" {
		user, _ := auth.GetUserFromToken(ctx, req.Gigatoken)
		if user != nil {
			userID = user.ID
		}
	}
	// get the live posts
	point := &types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	limit := &types.Limit{Start: req.StartLimit, End: req.EndLimit}
	var readyDatetime time.Time
	if req.ReadyDateTime != 0 {
		readyDatetime = time.Unix(int64(req.ReadyDateTime), 0)
	} else {
		readyDatetime = time.Now()
	}
	postIDs, itemIDs, chefIDs, distances, err := post.GetLivePostsIDs(ctx, point, limit, req.Radius, readyDatetime, req.Descending)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// get posts
	var posts []post.Post
	postErrChan := make(chan error, 1)
	go func() {
		var postErr error
		posts, postErr = post.GetPosts(ctx, postIDs)
		postErrChan <- postErr
	}()
	// get chef ratings
	var ratings []gigachef.Rating
	var chefDetails []types.UserDetail
	chefErrChan := make(chan error, 1)
	go func() {
		var chefErr error
		chefDetails, ratings, chefErr = gigachef.GetRatingsAndInfo(ctx, chefIDs)
		chefErrChan <- chefErr
	}()
	// get likes
	var likes []bool
	var numLikes []int
	likeErrChan := make(chan error, 1)
	go func() {
		var likeErr error
		likeC := like.New(ctx)
		likes, numLikes, likeErr = likeC.LikesItems(userID, itemIDs)
		likeErrChan <- likeErr
	}()
	// handle errors
	err = <-chefErrChan
	if err == nil {
		err = <-postErrChan
	}
	if err == nil {
		err = <-likeErrChan
	}
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}

	// TODO: check if user is in delivery range
	resp.Posts = make([]Post, len(postIDs))
	for i := range postIDs {
		resp.Posts[i].set(userID, postIDs[i], numLikes[i], likes[i], distances[i], &ratings[i], chefDetails[i].Name, chefDetails[i].PhotoURL, &posts[i])
	}
	return resp, nil
}

// GigachefDetailed is the detailed info for a Gigachef
type GigachefDetailed struct {
	ID              string  `json:"id,omitempty"`
	Name            string  `json:"name,omitempty"`
	PhotoURL        string  `json:"photo_url,omitempty"`
	NumOrders       int     `json:"num_orders,omitempty"`
	gigachef.Rating         // embedded
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

// set takes chef info and saves it to an endpoint GigachefDetails
func (g *GigachefDetailed) set(id, name, photoURL string, ratings gigachef.Rating, numOrders int, latitude, longitude float64) {
	g.ID = id
	g.Name = name
	g.PhotoURL = photoURL
	g.Rating = ratings
	g.NumOrders = numOrders
	g.Latitude = latitude
	g.Longitude = longitude
}

// ExchangeTimeSegment is the time range where the exchange can be made
type ExchangeTimeSegment struct {
	StartDateTime            int     `json:"start_datetime"`
	EndDateTime              int     `json:"end_datetime"`
	AvailableExchangeMethods int32   `json:"available_exchange_methods"`
	Price                    float32 `json:"price"`
	Index                    int     `json:"index"`
}

// PostDetailed has detailed information for a Post.
type PostDetailed struct {
	BasePost                               // embedded
	ExchangeTimes    []ExchangeTimeSegment `json:"exchange_times,omitempty"`
	Ingredients      []string              `json:"ingredients,omitempty"`
	DietaryNeedsTags []string              `json:"dietary_needs_tags,omitempty"`
	GeneralTags      []string              `json:"general_tags,omitempty"`
	CuisineTags      []string              `json:"cuisine_tags,omitempty"`
}

// set takes a post.Post package and converts it to an endpoint PostDetailed
func (p *PostDetailed) set(userID string, id int64, numLikes int, hasLiked bool, distance float32, post *post.Post) {
	p.BasePost.set(userID, id, numLikes, hasLiked, distance, post)
	p.Ingredients = post.Ingredients
	p.DietaryNeedsTags = post.DietaryNeedsTags
	p.GeneralTags = post.GeneralTags
	p.CuisineTags = post.CuisineTags
	now := time.Now()
	for i, exchangeTime := range post.ExchangeTimes {
		if !exchangeTime.AvailableExchangeMethods.IsZero() && exchangeTime.StartDateTime.After(now) {
			// add pickup option
			if exchangeTime.AvailableExchangeMethods.Pickup() {
				ets := ExchangeTimeSegment{
					StartDateTime:            ttoi(exchangeTime.StartDateTime),
					EndDateTime:              ttoi(exchangeTime.EndDateTime),
					AvailableExchangeMethods: int32(types.PickupOnlyExchangeMethod),
					Index: i,
				}
				// add exchange option
				p.ExchangeTimes = append(p.ExchangeTimes, ets)
			}
			// set chef delivery to false if out of delivery radius
			if distance > float32(post.GigachefDelivery.Radius) {
				exchangeTime.AvailableExchangeMethods.SetChefDelivery(false)
			}
			// add delivery option
			if exchangeTime.AvailableExchangeMethods.Delivery() {
				var price float32
				if exchangeTime.AvailableExchangeMethods.ChefDelivery() {
					price = post.GigachefDelivery.BasePrice
				} // else calculate price
				exchangeTime.AvailableExchangeMethods.SetPickup(false)
				ets := ExchangeTimeSegment{
					StartDateTime:            ttoi(exchangeTime.StartDateTime),
					EndDateTime:              ttoi(exchangeTime.EndDateTime),
					AvailableExchangeMethods: int32(exchangeTime.AvailableExchangeMethods),
					Index: i,
					Price: price,
				}
				// add exchange option
				p.ExchangeTimes = append(p.ExchangeTimes, ets)
			}
		}
	}
}

// PaymentInfo has the payment info
type PaymentInfo struct {
	ChefPricePerServing float32 `json:"chef_price_per_serving"`
	GigaFee             float32 `json:"giga_fee"`
	TaxPercentage       float32 `json:"tax_percentage"`
}

func (p *PaymentInfo) set(post *post.Post) {
	p.ChefPricePerServing = post.ChefPricePerServing
	p.GigaFee = post.PricePerServing - post.ChefPricePerServing
	p.TaxPercentage = post.TaxPercentage
}

// GetPostReq is the input required to get a post
type GetPostReq struct {
	Gigatoken string      `json:"gigatoken"`
	PostID    json.Number `json:"post_id"`
	PostID64  int64       `json:"-"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Radius    int         `json:"radius"`
}

// valid returns an error if input in invalid
func (req *GetPostReq) valid() error {
	var err error
	req.PostID64, err = req.PostID.Int64()
	if err != nil {
		return fmt.Errorf("error with PostID: %v", err)
	}
	return nil
}

// GetPostResp is the response for getting a post
type GetPostResp struct {
	Post        PostDetailed         `json:"post"`
	Gigachef    GigachefDetailed     `json:"gigachef"`
	Reviews     []Review             `json:"reviews,omitempty"`
	PaymentInfo PaymentInfo          `json:"payment_info"`
	Err         errors.ErrorWithCode `json:"err"`
}

// GetPost gets the details for a post
func (service *Service) GetPost(ctx context.Context, req *GetPostReq) (*GetPostResp, error) {
	resp := new(GetPostResp)
	defer handleResp(ctx, "GetPost", resp.Err)
	err := req.valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	p, err := post.GetPost(ctx, req.PostID64)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// Get Gigachef details
	var chef *gigachef.Gigachef
	chefErrChan := make(chan error, 1)
	go func() {
		var chefErr error
		chef, chefErr = gigachef.GetInfo(ctx, p.GigachefID)
		chefErrChan <- chefErr
	}()
	// Get reviews
	var reviews []review.Resp
	reviewsErrChan := make(chan error, 1)
	go func() {
		limit := &types.Limit{
			Start: 0,
			End:   5,
		}
		reviewC := review.New(ctx)
		var reviewErr error
		reviews, reviewErr = reviewC.GetReviews(p.GigachefID, limit, p.ItemID)
		reviewsErrChan <- reviewErr
	}()
	// Get distance
	var distance float32
	distanceErrChan := make(chan error, 1)
	go func() {
		chefPoint := p.GigachefAddress.GeoPoint
		muncherPoint := types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
		var distanceErr error
		distance, _, distanceErr = maps.GetDistance(ctx, chefPoint, muncherPoint)
		distanceErrChan <- distanceErr
	}()
	// check for errors for get Gigachef details and get reviews
	err = <-chefErrChan
	if err == nil {
		err = <-reviewsErrChan
	}
	if err == nil {
		err = <-distanceErrChan
	}
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// save reviews and gigachef details
	for i := range reviews {
		r := Review{}
		r.set(&reviews[i])
		resp.Reviews = append(resp.Reviews, r)
	}
	// get user
	var userID string
	if req.Gigatoken != "" {
		user, _ := auth.GetUserFromToken(ctx, req.Gigatoken)
		if user != nil {
			userID = user.ID
		}
	}
	resp.Gigachef.set(p.GigachefID, chef.Name, chef.PhotoURL, chef.Rating, chef.NumOrders, chef.Address.Latitude, chef.Address.Longitude)

	numLikes := 0
	hasLiked := false
	likeC := like.New(ctx)
	hasLikedList, numLikesList, err := likeC.LikesItems(userID, []int64{p.ItemID})
	if err != nil {
		utils.Errorf(ctx, "failed to get like.LikesItem: %v", err)
	} else {
		hasLiked = hasLikedList[0]
		numLikes = numLikesList[0]
	}
	resp.Post.set(userID, req.PostID64, numLikes, hasLiked, distance, p)
	resp.PaymentInfo.set(p)
	return resp, nil
}
