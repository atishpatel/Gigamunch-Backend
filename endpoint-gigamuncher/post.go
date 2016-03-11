package gigamuncher

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/core/post"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

// BasePost is the basic stuff for a post
type BasePost struct {
	ID                int      `json:"id"`
	ItemID            int      `json:"item_id"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	PricePerServing   float32  `json:"price_per_serving"`
	ServingsOffered   int      `json:"servings_offered"`
	ServingsLeft      int      `json:"servings_left"`
	Photos            []string `json:"photos"`
	PostedDateTime    int      `json:"posted_datetime"`
	ClosingDateTime   int      `json:"closing_datetime"`
	ReadyDateTime     int      `json:"ready_datetime"`
	Distance          float32  `json:"distance"`
	DeliveryAvailable bool     `json:"delivery_avaliable"`
	PickupAvaliable   bool     `json:"pickup_avaliable"`
	HasBought         bool     `json:"has_bought"`
	GigachefID        string   `json:"gigachef_id"`
	GigachefName      string   `json:"gigachef_name"`
	GigachefPhotoURL  string   `json:"gigachef_photo_url"`
}

// Post is a meal that is no longer live
type Post struct {
	BasePost               // embedded
	GigachefRating float32 `json:"gigachef_rating"`
}

type PostDetailed struct {
	BasePost                  // embedded
	Ingredients      []string `json:"ingredients"`
	DietaryNeedsTags []string `json:"dietary_needs_tags"`
	GeneralTags      []string `json:"general_tags"`
	CuisineTags      []string `json:"cuisine_tags"`
}

// GetLivePostsReq is the input required to get a list of live posts
type GetLivePostsReq struct {
	StartLimit    int     `json:"start_limit"`
	EndLimit      int     `json:"end_limit" endpoints:"req"`
	Latitude      float32 `json:"latitude" endpoints:"req"`
	Longitude     float32 `json:"longitude" endpoints:"req"`
	Radius        int     `json:"radius"`
	ReadyDateTime int     `json:"ready_datetime"`
	Decending     bool    `json:"decending"`
	UserID        string  `json:"user_id"`
}

// Valid returns an error if input in invalid
func (req *GetLivePostsReq) Valid() error {
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
//returns: posts, error
type GetLivePostsResp struct {
	Posts []Post               `json:"posts"`
	Err   errors.ErrorWithCode `json:"err"`
}

// GetLivePosts is an endpoint that returns a list of live posts
func (service *Service) GetLivePosts(ctx context.Context, req *GetLivePostsReq) (*GetLivePostsResp, error) {
	resp := new(GetLivePostsResp)
	var err error
	err = req.Valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	point := &types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	limit := &types.Limit{Start: req.StartLimit, End: req.EndLimit}
	readyDatetime := time.Unix(int64(req.ReadyDateTime), 0)
	// get the live posts
	postIDs, gigachefIDs, distances, err := post.GetLivePostsIDs(ctx, point, limit, req.Radius, readyDatetime, req.Decending)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	utils.Debugf(ctx, "postids: ", postIDs)
	// get posts
	var posts []post.Post
	postErrChan := make(chan error, 1)
	go func() {
		var err error
		posts, err = post.GetPosts(ctx, postIDs)
		postErrChan <- err
	}()
	// get chef ratings
	var ratings []float32
	chefErrChan := make(chan error, 1)
	go func() {
		var err error
		ratings, err = gigachef.GetRatings(ctx, gigachefIDs)
		chefErrChan <- err
	}()
	err = <-chefErrChan
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	err = <-postErrChan
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// TODO: check if user is in delivery range
	resp.Posts = make([]Post, len(postIDs))
	for i := range postIDs { // TODO switch to post.Set(post.Post)
		resp.Posts[i].ID = int(postIDs[i])
		resp.Posts[i].Distance = distances[i]
		resp.Posts[i].GigachefRating = ratings[i]
		resp.Posts[i].GigachefID = posts[i].GigachefID
		resp.Posts[i].GigachefName = posts[i].GigachefID
		resp.Posts[i].GigachefPhotoURL = posts[i].GigachefID
		resp.Posts[i].Description = posts[i].Description
		resp.Posts[i].ItemID = int(posts[i].ItemID)
		resp.Posts[i].ReadyDateTime = int(posts[i].ReadyDateTime.Unix())
		resp.Posts[i].ClosingDateTime = int(posts[i].ClosingDateTime.Unix())
		resp.Posts[i].PostedDateTime = int(posts[i].CreatedDateTime.Unix())
		resp.Posts[i].Photos = posts[i].Photos
		resp.Posts[i].ServingsOffered = posts[i].ServingsOffered
		resp.Posts[i].ServingsLeft = posts[i].ServingsOffered - posts[i].NumOrders
		resp.Posts[i].PricePerServing = posts[i].PricePerServing
		if req.UserID != "" {
			for j := range posts[i].Orders {
				if posts[i].Orders[j].GigamuncherID == req.UserID {
					resp.Posts[i].HasBought = true
				}
			}
		}
		resp.Posts[i].PickupAvaliable = posts[i].AvaliableExchangeMethods.Pickup()

		// TODO check if user is in delivery range
		// DeliveryAvailable bool     `json:"delivery_avaliable"`
	}
	return resp, nil
}
