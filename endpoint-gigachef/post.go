package gigachef

import (
	"encoding/json"
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/post"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// Post is a meal that is no longer live
type Post struct {
	BaseItem                        // embedded
	ID                  json.Number `json:"id"`
	ID64                int64       `json:"-"`
	ItemID              json.Number `json:"item_id"`
	ItemID64            int64       `json:"-"`
	Title               string      `json:"title"`
	ClosingDateTime     int         `json:"closing_datetime" endpoints:"req"`
	ReadyDateTime       int         `json:"ready_datetime" endpoints:"req"`
	ServingsOffered     int32       `json:"servings_offered" endpoints:"req"`
	ChefPricePerServing float32     `json:"chef_price_per_serving"`
	Pickup              bool        `json:"pickup"`
	GigachefDelivery    bool        `json:"gigachef_delivery"`
}

// Set takes a post.Post and converts it to a endpoint post
func (p *Post) Set(id int64, post *post.Post) {
	p.ID = itojn(id)
	p.ItemID = itojn(post.ItemID)
	p.Title = post.Title
	p.Description = post.Description
	p.Ingredients = post.Ingredients
	p.GeneralTags = post.GeneralTags
	p.DietaryNeedsTags = post.DietaryNeedsTags
	p.Photos = post.Photos
	p.ClosingDateTime = ttoi(post.ClosingDateTime)
	p.ReadyDateTime = ttoi(post.ReadyDateTime)
	p.ServingsOffered = post.ServingsOffered
	p.ChefPricePerServing = post.ChefPricePerServing
	p.Pickup = post.AvailableExchangeMethods.Pickup()
	p.GigachefDelivery = post.AvailableExchangeMethods.ChefDelivery()
}

// Get creates a post.Post version of the endpoint post
func (p *Post) Get() *post.Post {
	post := new(post.Post)
	post.ItemID = p.ItemID64
	post.Title = p.Title
	post.Description = p.Description
	post.Ingredients = p.Ingredients
	post.GeneralTags = p.GeneralTags
	post.DietaryNeedsTags = p.DietaryNeedsTags
	post.Photos = p.Photos
	post.ClosingDateTime = itot(p.ClosingDateTime)
	post.ReadyDateTime = itot(p.ReadyDateTime)
	post.ServingsOffered = p.ServingsOffered
	post.ChefPricePerServing = p.ChefPricePerServing
	post.AvailableExchangeMethods.SetPickup(p.Pickup)
	post.AvailableExchangeMethods.SetChefDelivery(p.GigachefDelivery)
	return post
}

// PostPostReq is the input request needed for PostPost.
type PostPostReq struct {
	Gigatoken string `json:"gigatoken"`
	Post      Post   `json:"post"`
}

// gigatoken returns the Gigatoken string
func (req *PostPostReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *PostPostReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	var err error
	req.Post.ItemID64, err = req.Post.ItemID.Int64()
	if err != nil {
		return fmt.Errorf("Error with ItemID: %v", err)
	}
	// TODO: check post stuff
	return nil
}

// PostPostResp is the output response for PostPost.
type PostPostResp struct {
	Post Post                 `json:"post"`
	Err  errors.ErrorWithCode `json:"err"`
}

// PostPost is an endpoint that post a post form a Gigachef
func (service *Service) PostPost(ctx context.Context, req *PostPostReq) (*PostPostResp, error) {
	resp := new(PostPostResp)
	defer handleResp(ctx, "PostPost", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	p := req.Post.Get()
	postID, err := post.PostPost(ctx, user, p)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Post.Set(postID, p)
	return resp, nil
}
