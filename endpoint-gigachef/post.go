package gigachef

import (
	"fmt"
	"time"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/order"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/post"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// OrderPaymentInfo is the payment information related to an order
type OrderPaymentInfo struct {
	Price         float32 `json:"price"`
	ExchangePrice float32 `json:"exchange_price"`
	GigaFee       float32 `json:"giga_fee"`
	TaxPrice      float32 `json:"tax_price"`
	TotalPrice    float32 `json:"total_price"`
}

func (opi *OrderPaymentInfo) set(pi *order.PaymentInfo) {
	opi.Price = pi.Price
	opi.ExchangePrice = pi.ExchangePrice
	opi.GigaFee = pi.GigaFee
	opi.TaxPrice = pi.TaxPrice
	opi.TotalPrice = pi.TotalPrice
}

// ExchangePlanInfo is the plan info
type ExchangePlanInfo struct {
	GigamuncherAddress types.Address `json:"gigamuncher_address"`
	GigachefAddress    types.Address `json:"gigachef_address"`
	Distance           float32       `json:"distance"`
	Duration           int           `json:"duration"`
}

func (epi *ExchangePlanInfo) set(o *order.Order) {
	epi.GigachefAddress = o.ExchangePlanInfo.GigachefAddress
	epi.GigamuncherAddress = o.ExchangePlanInfo.GigamuncherAddress
	epi.Distance = o.ExchangePlanInfo.Distance
	epi.Duration = int(o.ExchangePlanInfo.Duration)
}

// OrderGigamuncher is a gigamuncher info
type OrderGigamuncher struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}

// Order is an order
type Order struct {
	ID                       string           `json:"id"`
	CreatedDateTime          int              `json:"created_datetime"`
	ExpectedExchangeDateTime int              `json:"expected_exchange_datetime"`
	State                    string           `json:"state"`
	ZendeskIssueID           string           `json:"zendesk_issue_id"`
	GigachefCanceled         bool             `json:"gigachef_canceled"`
	GigamuncherCanceled      bool             `json:"gigamuncher_canceled"`
	Gigamuncher              OrderGigamuncher `json:"gigamuncher"`
	ReviewID                 string           `json:"review_id,omitempty"`
	PostID                   string           `json:"post_id,omitempty"`
	ItemID                   string           `json:"item_id,omitempty"`
	PostTitle                string           `json:"post_title"`
	PostPhotoURL             string           `json:"post_photo_url"`
	PricePerServing          float32          `json:"price_per_serving"`
	Servings                 int32            `json:"servings"`
	PaymentInfo              OrderPaymentInfo `json:"payment_info"`
	ExchangeMethod           int32            `json:"exchange_method"`
	ExchangePlanInfo         ExchangePlanInfo `json:"exchange_plan_info"`
}

func (o *Order) set(id int64, order *order.Order) {
	o.ID = itos(id)
	o.CreatedDateTime = ttoi(order.CreatedDateTime)
	o.ExpectedExchangeDateTime = ttoi(order.ExpectedExchangeDateTime)
	o.State = order.State
	o.ZendeskIssueID = itos(order.ZendeskIssueID)
	o.GigachefCanceled = order.GigachefCanceled
	o.GigamuncherCanceled = order.GigamuncherCanceled
	o.Gigamuncher.ID = order.GigamuncherID
	o.Gigamuncher.Name = order.GigamuncherName
	o.Gigamuncher.PhotoURL = order.GigamuncherPhotoURL
	o.ReviewID = itos(order.ReviewID)
	o.PostID = itos(order.PostID)
	o.ItemID = itos(order.ItemID)
	o.PostTitle = order.PostTitle
	o.PostPhotoURL = order.PostPhotoURL
	o.PricePerServing = order.PricePerServing
	o.Servings = order.Servings
	o.PaymentInfo.set(&order.PaymentInfo)
	o.ExchangeMethod = int32(order.ExchangeMethod)
	o.ExchangePlanInfo.set(order)
}

// PostOrder is an order for a post
type PostOrder struct {
	OrderID             string `json:"order_id"`
	GigamuncherID       string `json:"gigamuncher_id"`
	GigamuncherName     string `json:"gigamuncher_name"`
	GigamuncherPhotoURL string `json:"gigamuncher_photo_url"`
	ExchangeTime        int    `json:"exchange_time"`
	ExchangeMethod      int    `json:"exchange_method"`
	Servings            int    `json:"servings"`
}

func (po *PostOrder) set(postOrder *post.OrderPost) {
	po.OrderID = itos(postOrder.OrderID)
	po.GigamuncherID = postOrder.GigamuncherID
	po.GigamuncherName = postOrder.GigamuncherName
	po.GigamuncherPhotoURL = postOrder.GigamuncherPhotoURL
	po.ExchangeTime = ttoi(postOrder.ExchangeTime)
	po.ExchangeMethod = int(postOrder.ExchangeMethod)
	po.Servings = int(postOrder.Servings)
}

// Post is a meal that is no longer live
type Post struct { // TODO add num likes and stuff
	BaseItem                             // embedded
	ID                       string      `json:"id"`
	ID64                     int64       `json:"-"`
	ItemID                   string      `json:"item_id"`
	ItemID64                 int64       `json:"-"`
	Title                    string      `json:"title"`
	ClosingDateTime          int         `json:"closing_datetime"`
	ReadyDateTime            int         `json:"ready_datetime"`
	ServingsOffered          string      `json:"servings_offered"`
	ServingsOffered32        int32       `json:"-"`
	ChefPricePerServing      string      `json:"chef_price_per_serving"`
	ChefPricePerServing32    float32     `json:"-"`
	EstimatedPreperationTime int         `json:"estimated_preperation_time"`
	Pickup                   bool        `json:"pickup"`
	GigachefDelivery         bool        `json:"gigachef_delivery"`
	IsOrderNow               bool        `json:"is_order_now"`
	Orders                   []PostOrder `json:"orders"`
}

// Set takes a post.Post and converts it to a endpoint post
func (p *Post) set(id int64, post *post.Post) {
	p.ID = itos(id)
	p.ItemID = itos(post.ItemID)
	p.Title = post.Title
	p.Description = post.Description
	p.Ingredients = post.Ingredients
	p.GeneralTags = post.GeneralTags
	p.DietaryNeedsTags = post.DietaryNeedsTags
	p.Photos = post.Photos
	p.ClosingDateTime = ttoi(post.ClosingDateTime)
	p.ReadyDateTime = ttoi(post.ReadyDateTime)
	p.ServingsOffered = itos(int64(post.ServingsOffered))
	p.ChefPricePerServing = ftos(float64(post.ChefPricePerServing))
	p.EstimatedPreperationTime = int(post.EstimatedPreperationTime)
	p.Pickup = post.AvailableExchangeMethods.Pickup()
	p.GigachefDelivery = post.AvailableExchangeMethods.ChefDelivery()
	p.IsOrderNow = post.IsOrderNow
	p.Orders = make([]PostOrder, len(post.Orders))
	for i := range post.Orders {
		p.Orders[i].set(&post.Orders[i])
	}
}

// Get creates a post.Post version of the endpoint post
func (p *Post) get() *post.Post {
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
	post.ServingsOffered = p.ServingsOffered32
	post.ChefPricePerServing = p.ChefPricePerServing32
	post.EstimatedPreperationTime = int64(p.EstimatedPreperationTime)
	post.AvailableExchangeMethods.SetPickup(p.Pickup)
	post.AvailableExchangeMethods.SetChefDelivery(p.GigachefDelivery)
	post.IsOrderNow = p.IsOrderNow
	return post
}

func (p *Post) valid() error {
	var err error
	p.ItemID64, err = stoi(p.ItemID)
	if err != nil {
		return fmt.Errorf("Error with item_id: %v", err)
	}
	p.ServingsOffered32, err = stoi32(p.ServingsOffered)
	if err != nil {
		return fmt.Errorf("Error with servings_offered: %v", err)
	}
	p.ChefPricePerServing32, err = stof(p.ChefPricePerServing)
	if err != nil {
		return fmt.Errorf("Error with chef_price_per_serving: %v", err)
	}
	return nil
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
	now := int(time.Now().Unix())
	if req.Post.ClosingDateTime < now {
		return fmt.Errorf("Closing DateTime cannot be before now.")
	}
	if !req.Post.IsOrderNow && req.Post.ReadyDateTime < now {
		return fmt.Errorf("Ready DateTime cannot be before now.")
	}
	return req.Post.valid()
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
	p := req.Post.get()
	postID, err := post.PostPost(ctx, user, p)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Post.set(postID, p)
	return resp, nil
}

// GetPostsReq is a req for GetPosts
type GetPostsReq struct {
	Gigatoken  string `json:"gigatoken"`
	StartLimit int    `json:"start_limit"`
	EndLimit   int    `json:"end_limit"`
}

func (req *GetPostsReq) gigatoken() string {
	return req.Gigatoken
}

func (req *GetPostsReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	if req.StartLimit >= req.EndLimit {
		return fmt.Errorf("Limit range isn't valid.")
	}
	return nil
}

// GetPostsResp has a list of posts and error
type GetPostsResp struct {
	Posts []Post               `json:"posts"`
	Err   errors.ErrorWithCode `json:"err"`
}

// GetPosts gets a chef's posts
func (service *Service) GetPosts(ctx context.Context, req *GetPostsReq) (*GetPostsResp, error) {
	resp := new(GetPostsResp)
	defer handleResp(ctx, "GetPosts", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	postC := post.New(ctx)
	postIDs, posts, err := postC.GetUserPosts(user.ID, req.StartLimit, req.EndLimit)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to post.GetUserPosts")
		return resp, nil
	}
	resp.Posts = make([]Post, len(postIDs))
	for i := range postIDs {
		resp.Posts[i].set(postIDs[i], &posts[i])
	}
	return resp, nil
}
