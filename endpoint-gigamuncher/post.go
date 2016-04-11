package gigamuncher

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/core/post"
	"github.com/atishpatel/Gigamunch-Backend/core/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils/maps"
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
}

// Set takes a post.Post package and converts it to a endpoint post
func (p *BasePost) Set(id int, distance float32, post *post.Post) {
	p.ID = id
	p.Distance = distance
	p.Description = post.Description
	p.ItemID = int(post.ItemID)
	p.ReadyDateTime = int(post.ReadyDateTime.Unix())
	p.ClosingDateTime = int(post.ClosingDateTime.Unix())
	p.PostedDateTime = int(post.CreatedDateTime.Unix())
	p.Photos = post.Photos
	p.ServingsOffered = post.ServingsOffered
	p.ServingsLeft = post.ServingsOffered - post.NumOrders
	p.PricePerServing = post.PricePerServing
	p.PickupAvaliable = post.AvaliableExchangeMethods.Pickup()
}

// PostGigachef is the basic version of GigachefDetails for post in live feed
type PostGigachef struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	PhotoURL string  `json:"photo_url"`
	Rating   float32 `json:"rating"`
}

// Post is a meal that is no longer live
type Post struct {
	BasePost              // embedded
	Gigachef PostGigachef `json:"gigachef"`
}

// Set takes a post.Post package and converts it to a endpoint post
func (p *Post) Set(id int, distance float32, gigachefName string, gigachefPhotoURL string, avgRating float32, post *post.Post) {
	p.BasePost.Set(id, distance, post)
	p.Gigachef = PostGigachef{
		ID:       post.GigachefID,
		Name:     gigachefName,
		PhotoURL: gigachefPhotoURL,
		Rating:   avgRating,
	}
}

// GetLivePostsReq is the input required to get a list of live posts
type GetLivePostsReq struct {
	StartLimit    int     `json:"start_limit"`
	EndLimit      int     `json:"end_limit"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
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
	postIDs, chefIDs, distances, err := post.GetLivePostsIDs(ctx, point, limit, req.Radius, readyDatetime, req.Decending)
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
	var ratings []gigachef.GigachefRating
	var chefDetails []types.UserDetail
	chefErrChan := make(chan error, 1)
	go func() {
		var chefErr error
		chefDetails, ratings, chefErr = gigachef.GetRatingsAndInfo(ctx, chefIDs)
		chefErrChan <- chefErr
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
		resp.Posts[i].Set(int(postIDs[i]), distances[i], chefDetails[i].Name, chefDetails[i].PhotoURL, ratings[i].AverageRating, &posts[i])
		if req.UserID != "" {
			for j := range posts[i].Orders {
				if posts[i].Orders[j].GigamuncherID == req.UserID {
					resp.Posts[i].HasBought = true
				}
			}
		}

		// TODO check if user is in delivery range
		// DeliveryAvailable bool     `json:"delivery_avaliable"`
	}
	return resp, nil
}

// GigachefDetailed is the detailed info for a Gigachef
type GigachefDetailed struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	PhotoURL                string `json:"photo_url"`
	gigachef.GigachefRating        // embedded
	NumOrders               int    `json:"num_orders"`
}

// Set takes chef info and saves it to an endpoint GigachefDetails
func (g *GigachefDetailed) Set(id, name, photoURL string, ratings gigachef.GigachefRating, numOrders int) {
	g.ID = id
	g.Name = name
	g.PhotoURL = photoURL
	g.GigachefRating = ratings
	g.NumOrders = numOrders
}

// PostDetailed has detailed information for a Post.
type PostDetailed struct {
	BasePost                  // embedded
	Ingredients      []string `json:"ingredients"`
	DietaryNeedsTags []string `json:"dietary_needs_tags"`
	GeneralTags      []string `json:"general_tags"`
	CuisineTags      []string `json:"cuisine_tags"`
}

// Set takes a post.Post package and converts it to an endpoint PostDetailed
func (p *PostDetailed) Set(id int, distance float32, post *post.Post) {
	p.BasePost.Set(id, distance, post)
	p.Ingredients = post.Ingredients
	p.DietaryNeedsTags = post.DietaryNeedsTags
	p.GeneralTags = post.GeneralTags
	p.CuisineTags = post.CuisineTags
	// TODO add gigachef setting stuff
}

// GetPostReq is the input required to get a post
type GetPostReq struct {
	PostID    int     `json:"post_id,string"`
	UserID    string  `json:"user_id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Radius    int     `json:"radius"`
}

// Valid returns an error if input in invalid
func (req *GetPostReq) Valid() error {
	return nil
}

// GetPostResp is the response for getting a post
//returns: post, error
type GetPostResp struct {
	Post     PostDetailed         `json:"post"`
	Gigachef GigachefDetailed     `json:"gigachef"`
	Reviews  []Review             `json:"reviews"`
	Err      errors.ErrorWithCode `json:"err"`
}

// GetPost gets the details for a post
func (service *Service) GetPost(ctx context.Context, req *GetPostReq) (*GetPostResp, error) {
	resp := new(GetPostResp)
	err := req.Valid()
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	p, err := post.GetPost(ctx, int64(req.PostID))
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
	var reviewIDs []int64
	var reviews []review.Review
	reviewsErrChan := make(chan error, 1)
	go func() {
		limit := &types.Limit{
			Start: 0,
			End:   5,
		}
		var reviewErr error
		reviewIDs, reviews, reviewErr = review.GetReviews(ctx, p.GigachefID, limit, p.ItemID)
		reviewsErrChan <- reviewErr
	}()
	// Get distance
	var distance float32
	distanceErrChan := make(chan error, 1)
	go func() {
		chefPoint := p.Address.GeoPoint
		muncherPoint := types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
		var distanceErr error
		distance, _, distanceErr = maps.GetDistance(ctx, chefPoint, muncherPoint)
		distanceErrChan <- distanceErr
	}()
	// check for errors for get Gigachef details and get reviews
	err = <-chefErrChan
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	err = <-reviewsErrChan
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	err = <-distanceErrChan
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	// save reviews and gigachef details
	for i := range reviewIDs {
		resp.Reviews = append(resp.Reviews, Review{ID: int(reviewIDs[i]), Review: reviews[i]})
	}
	resp.Gigachef.Set(p.GigachefID, chef.Name, chef.PhotoURL, chef.GigachefRating, chef.NumOrders)
	resp.Post.Set(req.PostID, distance, p)
	return resp, nil
}
