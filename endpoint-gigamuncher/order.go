package gigamuncher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/core/gigamuncher"
	"github.com/atishpatel/Gigamunch-Backend/core/like"
	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/core/order"
	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/core/post"
	"github.com/atishpatel/Gigamunch-Backend/core/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
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

// OrderGigachef is a gigachef info
type OrderGigachef struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PhoneNumber     string `json:"phone_number"`
	PhotoURL        string `json:"photo_url"`
	NumOrders       int    `json:"num_orders"`
	gigachef.Rating        // embedded
}

// Order is an order
type Order struct {
	ID                       string           `json:"id,omitempty"`
	CreatedDateTime          int              `json:"created_datetime"`
	ExpectedExchangeDateTime int              `json:"expected_exchange_datetime"`
	State                    string           `json:"state"`
	ZendeskIssueID           string           `json:"zendesk_issue_id"`
	GigachefCanceled         bool             `json:"gigachef_canceled"`
	GigamuncherCanceled      bool             `json:"gigamuncher_canceled"`
	Gigachef                 OrderGigachef    `json:"gigachef"`
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
	NumLikes                 int              `json:"num_likes"`
	HasLiked                 bool             `json:"has_liked"`
	Status                   string           `json:"status"`
}

func (o *Order) set(orderID int64, order *order.Order, numLikes int, hasLiked bool, chefName, chefPhotoURL string, chefRatings gigachef.Rating, chefPhoneNumber string) {
	o.ID = itos(orderID)
	o.CreatedDateTime = ttoi(order.CreatedDateTime)
	o.ExpectedExchangeDateTime = ttoi(order.ExpectedExchangeDateTime)
	o.State = order.State
	o.ZendeskIssueID = itos(order.ZendeskIssueID)
	o.GigachefCanceled = order.GigachefCanceled
	o.GigamuncherCanceled = order.GigamuncherCanceled
	o.Gigachef.ID = order.GigachefID
	o.Gigachef.Name = chefName
	o.Gigachef.PhotoURL = chefPhotoURL
	o.Gigachef.PhoneNumber = chefPhoneNumber
	o.Gigachef.Rating = chefRatings
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
	o.NumLikes = numLikes
	o.HasLiked = hasLiked

	preparingStartTime := order.ExpectedExchangeDateTime.Add(-30 * time.Minute)
	if preparingStartTime.After(order.PostCloseDateTime) {
		preparingStartTime = order.PostCloseDateTime.Add(-30 * time.Minute)
	}

	if o.GigachefCanceled || o.GigamuncherCanceled {
		o.Status = "canceled"
	} else if time.Now().After(order.ExpectedExchangeDateTime.Add(1 * time.Hour)) {
		o.Status = "closed"
	} else if time.Now().After(order.ExpectedExchangeDateTime) {
		o.Status = "open-received"
	} else if preparingStartTime.After(order.PostCloseDateTime) {
		o.Status = "open-preparing"
	} else {
		o.Status = "open-placed"
	}
}

// MakeOrderReq is the request for MakeOrder
type MakeOrderReq struct {
	Gigatoken           string         `json:"gigatoken"`
	PostID              json.Number    `json:"post_id,omitempty"`
	PostID64            int64          `json:"-"`
	BraintreeNonce      string         `json:"braintree_nonce"`
	Servings            int32          `json:"servings"`
	ExchangeWindowIndex int32          `json:"exchange_window_index"`
	ExchangeMethods     int32          `json:"exchange_methods"`
	Latitude            float64        `json:"latitude"`  // REMOVE
	Longitude           float64        `json:"longitude"` // REMOVE
	Address             Address        `json:"address"`
	AddressType         *types.Address `json:"-"`
	TotalPrice          float32        `json:"total_price"`
}

func (req *MakeOrderReq) gigatoken() string {
	return req.Gigatoken
}

func (req *MakeOrderReq) valid() error {
	if req.BraintreeNonce == "" {
		return fmt.Errorf("BraintreeNonce is invalid.")
	}
	if req.Servings == 0 {
		return fmt.Errorf("Servings is 0.")
	}
	var err error
	req.PostID64, err = req.PostID.Int64()
	if err != nil {
		return fmt.Errorf("error with PostID: %v", err)
	}
	req.AddressType, err = req.Address.get()
	if err != nil {
		return fmt.Errorf("error with decoding address: %#v", err)
	}
	if req.AddressType.Longitude == 0 && req.AddressType.Latitude == 0 { // REMOVE
		req.AddressType.Longitude = req.Longitude
		req.AddressType.Latitude = req.Latitude
	}

	return nil
}

// MakeOrderResp is the resposne
type MakeOrderResp struct {
	Order Order                `json:"order"`
	Err   errors.ErrorWithCode `json:"err"`
}

// MakeOrder makes an order
func (service *Service) MakeOrder(ctx context.Context, req *MakeOrderReq) (*MakeOrderResp, error) {
	resp := new(MakeOrderResp)
	defer handleResp(ctx, "MakeOrder", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if req.AddressType.Longitude == 0 && req.AddressType.Latitude == 0 {
		err = maps.GetGeopointFromAddress(ctx, req.AddressType)
		if err != nil {
			resp.Err = errors.Wrap("failed to maps.GetGeopointFromAddress", err)
			return resp, nil
		}
	}
	exchangeMethods := types.ExchangeMethods(req.ExchangeMethods)
	postC := post.New(ctx)
	postReq := &post.MakeOrderReq{
		PostID:              req.PostID64,
		NumServings:         req.Servings,
		PaymentNonce:        req.BraintreeNonce,
		ExchangeMethod:      exchangeMethods,
		ExchangeWindowIndex: req.ExchangeWindowIndex,
		GigamuncherAddress:  *req.AddressType,
		GigamuncherID:       user.ID,
		GigamuncherName:     user.Name,
		GigamuncherPhotoURL: user.PhotoURL,
	}
	orderID, order, err := postC.MakeOrder(postReq)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}

	itemIDs := []int64{order.ItemID}
	likeC := like.New(ctx)
	likes, numLikes, err := likeC.LikesItems(user.ID, itemIDs)
	if err != nil {
		resp.Err = errors.Wrap("failed to get liked items", err)
		return resp, nil
	}

	chef, err := gigachef.GetInfo(ctx, order.GigachefID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("cannot chef.GetInfo")
		return resp, nil
	}

	resp.Order.set(orderID, order, numLikes[0], likes[0], chef.Name, chef.PhotoURL, chef.Rating, chef.PhoneNumber)
	return resp, nil
}

// GetOrderReq is the request for GetOrder
type GetOrderReq struct {
	OrderID   string `json:"order_id"`
	OrderID64 int64  `json:"-"`
	Gigatoken string `json:"gigatoken"`
}

func (req *GetOrderReq) gigatoken() string {
	return req.Gigatoken
}

func (req *GetOrderReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	var err error
	req.OrderID64, err = stoi(req.OrderID)
	if err != nil {
		return fmt.Errorf("error with OrderID: %v", err)
	}
	return nil
}

// GetOrderResp is the response  for GetOrder
type GetOrderResp struct {
	Order  Order                `json:"order"`
	Review Review               `json:"review"`
	Err    errors.ErrorWithCode `json:"err"`
}

// GetOrder gets an order
func (service *Service) GetOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, error) {
	resp := new(GetOrderResp)
	defer handleResp(ctx, "GetOrder", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}

	ordersC := order.New(ctx)
	order, err := ordersC.GetOrder(user.ID, req.OrderID64)
	if err != nil {
		resp.Err = errors.Wrap("cannot get order", err)
		return resp, nil
	}

	itemIDs := []int64{order.ItemID}
	likeC := like.New(ctx)
	likes, numLikes, err := likeC.LikesItems(user.ID, itemIDs)
	if err != nil {
		resp.Err = errors.Wrap("failed to get liked items", err)
		return resp, nil
	}
	chef, err := gigachef.GetInfo(ctx, order.GigachefID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("cannot chef.GetInfo")
		return resp, nil
	}

	resp.Order.set(order.ID, &order.Order, numLikes[0], likes[0], chef.Name, chef.PhotoURL, chef.Rating, chef.PhoneNumber)
	// get review
	reviewC := review.New(ctx)
	review, err := reviewC.GetReview(order.ReviewID)
	if err != nil {
		resp.Err = errors.Wrap("cannot get review", err)
		return resp, nil
	}
	resp.Review.set(review)
	return resp, nil
}

// GetOrdersReq is the request for GetOrder
type GetOrdersReq struct {
	StartLimit int    `json:"start_limit"`
	EndLimit   int    `json:"end_limit"`
	Gigatoken  string `json:"gigatoken"`
}

func (req *GetOrdersReq) gigatoken() string {
	return req.Gigatoken
}

func (req *GetOrdersReq) valid() error {
	if req.StartLimit < 0 || req.EndLimit < 0 {
		return fmt.Errorf("Limit is out of range.")
	}
	if req.EndLimit <= req.StartLimit {
		return fmt.Errorf("EndLimit cannot be less than or equal to StartLimit.")
	}
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// GetOrdersResp is the response for GetOrders
type GetOrdersResp struct {
	Orders []Order              `json:"orders,omitempty"`
	Err    errors.ErrorWithCode `json:"err"`
}

// GetOrders gets the orders for a muncher
func (service *Service) GetOrders(ctx context.Context, req *GetOrdersReq) (*GetOrdersResp, error) {
	resp := new(GetOrdersResp)
	defer handleResp(ctx, "GetOrders", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	limit := &types.Limit{Start: req.StartLimit, End: req.EndLimit}
	ordersC := order.New(ctx)
	orders, err := ordersC.GetOrders(user.ID, limit)
	if err != nil {
		resp.Err = errors.Wrap("cannot get order", err)
		return resp, nil
	}
	itemIDs := make([]int64, len(orders))
	chefIDs := make([]string, len(orders))
	for i := range orders {
		itemIDs[i] = orders[i].ItemID
		chefIDs[i] = orders[i].GigachefID
	}
	likeC := like.New(ctx)
	likes, numLikes, err := likeC.LikesItems(user.ID, itemIDs)
	if err != nil {
		resp.Err = errors.Wrap("failed to get liked items", err)
		return resp, nil
	}
	chefs, err := gigachef.GetMultiInfo(ctx, chefIDs)
	if err != nil {
		resp.Err = errors.Wrap("failed to chef.GetRatingAndInfo", err)
		return resp, nil
	}

	for i := range orders {
		o := Order{}
		o.set(orders[i].ID, &orders[i].Order, numLikes[i], likes[i], chefs[i].Name, chefs[i].PhotoURL, chefs[i].Rating, chefs[i].PhoneNumber)
		resp.Orders = append(resp.Orders, o)
	}
	return resp, nil
}

// GetBraintreeTokenResp is the response for GetBraintreeToken
type GetBraintreeTokenResp struct {
	BraintreeToken string               `json:"braintree_token"`
	Err            errors.ErrorWithCode `json:"err"`
}

// GetBraintreeToken gets a braintree token
func (service *Service) GetBraintreeToken(ctx context.Context, req *GigatokenOnlyReq) (*GetBraintreeTokenResp, error) {
	resp := new(GetBraintreeTokenResp)
	defer handleResp(ctx, "GetBraintree", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.Wrap("failed to validate", err)
		return resp, nil
	}
	customerID, err := gigamuncher.GetBTCustomerID(ctx, user.ID)
	if err != nil {
		resp.Err = errors.Wrap("cannot GetBTCustomerID", err)
		return resp, nil
	}
	paymentC := payment.New(ctx)
	token, err := paymentC.GenerateToken(customerID)
	if err != nil {
		resp.Err = errors.Wrap("cannot GenerateToken", err)
		return resp, nil
	}
	resp.BraintreeToken = token
	return resp, nil
}

// CancelOrderReq is the request for CancelOrder
type CancelOrderReq struct {
	OrderID   string `json:"order_id"`
	OrderID64 int64  `json:"-"`
	Gigatoken string `json:"gigatoken"`
}

func (req *CancelOrderReq) gigatoken() string {
	return req.Gigatoken
}

func (req *CancelOrderReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	var err error
	req.OrderID64, err = stoi(req.OrderID)
	if err != nil {
		return fmt.Errorf("error with OrderID: %v", err)
	}
	return nil
}

// CancelOrderResp is the response for CancelOrder
type CancelOrderResp struct {
	Order  Order                `json:"order"`
	Review Review               `json:"review"`
	Err    errors.ErrorWithCode `json:"err"`
}

// CancelOrder cancels an order
func (service *Service) CancelOrder(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, error) {
	resp := new(CancelOrderResp)
	defer handleResp(ctx, "CancelOrder", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: err.Error()}
		return resp, nil
	}
	postC := post.New(ctx)
	order, err := postC.CancelOrder(user.ID, req.OrderID64)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to cancel order")
		return resp, nil
	}
	// get items
	itemIDs := []int64{order.ItemID}
	likeC := like.New(ctx)
	likes, numLikes, err := likeC.LikesItems(user.ID, itemIDs)
	if err != nil {
		resp.Err = errors.Wrap("failed to get liked items", err)
		return resp, nil
	}
	chef, err := gigachef.GetInfo(ctx, order.GigachefID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("cannot chef.GetInfo")
		return resp, nil
	}
	// get review
	reviewC := review.New(ctx)
	review, err := reviewC.GetReview(order.ReviewID)
	if err != nil {
		resp.Err = errors.Wrap("cannot get review", err)
		return resp, nil
	}
	resp.Order.set(order.ID, &order.Order, numLikes[0], likes[0], chef.Name, chef.PhotoURL, chef.Rating, chef.PhoneNumber)
	resp.Review.set(review)
	return resp, nil
}
