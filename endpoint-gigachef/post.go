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
type Post struct {
	ID                  string      `json:"id"`
	BaseItem                        // embedded
	ItemID              string      `json:"item_id"`
	Title               string      `json:"title"`
	GigachefCanceled    bool        `json:"gigachef_canceled"`
	ClosingDateTime     int         `json:"closing_datetime"`
	ServingsOffered     string      `json:"servings_offered"`
	ChefPricePerServing string      `json:"chef_price_per_serving"`
	PricePerServing     string      `json:"price_per_serving"`
	Orders              []PostOrder `json:"orders"`
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
	p.ServingsOffered = itos(int64(post.ServingsOffered))
	p.ChefPricePerServing = ftos(float64(post.ChefPricePerServing))
	p.PricePerServing = ftos(float64(post.PricePerServing))
	p.Orders = make([]PostOrder, len(post.Orders))
	for i := range post.Orders {
		p.Orders[i].set(&post.Orders[i])
	}
}

// PublishPostReq is the input request needed for PublishPost.
type PublishPostReq struct {
	Gigatoken                 string  `json:"gigatoken"`
	BaseItem                          // embedded
	ItemID                    string  `json:"item_id"`
	ItemID64                  int64   `json:"-"`
	Title                     string  `json:"title"`
	ClosingDateTime           int     `json:"closing_datetime"`
	StartPickupDateTime       int     `json:"start_pickup_datetime"`
	EndPickupDateTime         int     `json:"end_pickup_datetime"`
	ChefDelivery              bool    `json:"chef_delivery"`
	StartChefDeliveryDateTime int     `json:"start_chef_delivery_datetime"`
	EndChefDeliveryDateTime   int     `json:"end_chef_delivery_datetime"`
	ChefDeliveryRadius        int32   `json:"chef_delivery_radius"`
	ChefDeliveryBasePrice     string  `json:"chef_delivery_base_price"`
	ChefDeliveryBasePrice32   float32 `json:"-"`
	ServingsOffered           string  `json:"servings_offered"`
	ServingsOffered32         int32   `json:"-"`
	ChefPricePerServing       string  `json:"chef_price_per_serving"`
	ChefPricePerServing32     float32 `json:"-"`
}

// Get creates a post.PublishPostReq
func (req *PublishPostReq) getPublishPostReq(user *types.User) *post.PublishPostReq {
	return &post.PublishPostReq{
		User: user,
		BaseItem: types.BaseItem{
			GigachefID:       user.ID,
			CreatedDateTime:  time.Now(),
			Description:      req.BaseItem.Description,
			GeneralTags:      req.BaseItem.GeneralTags,
			DietaryNeedsTags: req.BaseItem.DietaryNeedsTags,
			CuisineTags:      req.BaseItem.CuisineTags,
			Ingredients:      req.BaseItem.Ingredients,
			Photos:           req.BaseItem.Photos,
		},
		ItemID:                    req.ItemID64,
		Title:                     req.Title,
		ClosingDateTime:           itot(req.ClosingDateTime),
		StartPickupDateTime:       itot(req.StartPickupDateTime),
		EndPickupDateTime:         itot(req.EndPickupDateTime),
		ChefDelivery:              req.ChefDelivery,
		ChefDeliveryRadius:        int32(req.ChefDeliveryRadius),
		ChefDeliveryBasePrice:     req.ChefDeliveryBasePrice32,
		StartChefDeliveryDateTime: itot(req.StartChefDeliveryDateTime),
		EndChefDeliveryDateTime:   itot(req.EndChefDeliveryDateTime),
		ServingsOffered:           req.ServingsOffered32,
		ChefPricePerServing:       req.ChefPricePerServing32,
	}
}

// gigatoken returns the Gigatoken string
func (req *PublishPostReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *PublishPostReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	now := int(time.Now().Unix())
	if req.ClosingDateTime < now {
		return fmt.Errorf("Closing DateTime cannot be before now.")
	}
	var err error
	req.ItemID64, err = stoi(req.ItemID)
	if err != nil {
		return fmt.Errorf("Error with item_id: %v", err)
	}
	req.ServingsOffered32, err = stoi32(req.ServingsOffered)
	if err != nil {
		return fmt.Errorf("Error with servings_offered: %v", err)
	}
	req.ChefPricePerServing32, err = stof(req.ChefPricePerServing)
	if err != nil {
		return fmt.Errorf("Error with chef_price_per_serving: %v", err)
	}
	req.ChefDeliveryBasePrice32, err = stof(req.ChefDeliveryBasePrice)
	if err != nil {
		return fmt.Errorf("Error with chef delivery base price: %v", err)
	}
	return nil
}

// PublishPostResp is the output response for PublishPost.
type PublishPostResp struct {
	Post Post                 `json:"post"`
	Err  errors.ErrorWithCode `json:"err"`
}

// PublishPost is an endpoint that post a post form a Gigachef
func (service *Service) PublishPost(ctx context.Context, req *PublishPostReq) (*PublishPostResp, error) {
	resp := new(PublishPostResp)
	defer handleResp(ctx, "PublishPost", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	postC := post.New(ctx)
	publishPostReq := req.getPublishPostReq(user)
	postID, p, err := postC.PublishPost(publishPostReq)
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
