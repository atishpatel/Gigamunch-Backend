package gigachef

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/post"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

// Post is a meal that is no longer live
type Post struct {
	BaseItem                // embedded
	ID              int     `json:"id"`
	ItemID          int     `json:"item_id"`
	Title           string  `json:"title"`
	ClosingDateTime int     `json:"closing_datetime" endpoints:"req"`
	ReadyDateTime   int     `json:"ready_datetime" endpoints:"req"`
	ServingsOffered int     `json:"servings_offered" endpoints:"req"`
	PricePerServing float32 `json:"price_per_serving"`
}

// Set takes a post.Post and converts it to a endpoint post
func (p *Post) Set(id int, post *post.Post) {
	p.ID = id
	p.ItemID = int(post.ItemID)
	p.Title = post.Title
	p.Subtitle = post.Subtitle
	p.Description = post.Description
	p.Ingredients = post.Ingredients
	p.GeneralTags = post.GeneralTags
	p.DietaryNeedsTags = post.DietaryNeedsTags
	p.Photos = post.Photos
	p.ClosingDateTime = int(post.ClosingDateTime.Unix())
	p.ReadyDateTime = int(post.ReadyDateTime.Unix())
	p.ServingsOffered = post.ServingsOffered
	p.PricePerServing = post.PricePerServing
}

// Get creates a post.Post version of the endpoint post
func (p *Post) Get() *post.Post {
	post := new(post.Post)
	post.ItemID = int64(p.ItemID)
	post.Title = p.Title
	post.Subtitle = p.Subtitle
	post.Description = p.Description
	post.Ingredients = p.Ingredients
	post.GeneralTags = p.GeneralTags
	post.DietaryNeedsTags = p.DietaryNeedsTags
	post.Photos = p.Photos
	post.ClosingDateTime = time.Unix(int64(p.ClosingDateTime), 0)
	post.ReadyDateTime = time.Unix(int64(p.ReadyDateTime), 0)
	post.ServingsOffered = p.ServingsOffered
	post.PricePerServing = p.PricePerServing
	return post
}

type PostWithOrders struct {
	Post // embedded
	// Orders []
}

// PostPostReq is the input request needed for PostPost.
type PostPostReq struct {
	GigaToken string `json:"gigatoken"`
	Post      Post   `json:"post" endpoints:"req"`
}

// Gigatoken returns the GigaToken string
func (req *PostPostReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *PostPostReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
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
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "PostPost err: ", resp.Err)
		}
	}()
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
	resp.Post.Set(int(postID), p)
	return resp, nil
}
