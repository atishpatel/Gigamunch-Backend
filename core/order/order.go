package order

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/core/post"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

var (
	errDatastore         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errNotEnoughServings = errors.ErrorWithCode{Code: errors.CodeNotEnoughServingsLeft, Message: "Not enough servings left."}
	errInvalidParameter  = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
)

// Client is a client for orders
type Client struct {
	ctx context.Context
}

// New returns a new Client
func New(ctx context.Context) Client {
	return Client{ctx: ctx}
}

// GetOrderIDsAndPostInfo gets the user ids involved with the order
// returns: BasicOrderIDs, PostInfo, error
func GetOrderIDsAndPostInfo(ctx context.Context, orderID int64) (BasicOrderIDs, PostInfo, error) {
	var postInfo PostInfo
	order := new(Order)
	err := get(ctx, orderID, order)
	if err != nil {
		var basicOrderIDs BasicOrderIDs
		return basicOrderIDs, postInfo, errDatastore.WithError(err)
	}
	postInfo.ID = order.PostID
	postInfo.Title = order.PostTitle
	postInfo.PhotoURL = order.PostPhotoURL
	return order.BasicOrderIDs, postInfo, nil
}

// Resp contains the order id and order itself
type Resp struct {
	ID    int64
	Order // embedded
}

type postClient interface {
	GetOrderInfo(int64) (*post.OrderInfoResp, error)
	AddOrder(context.Context, *post.AddOrderReq) error
}

type paymentClient interface {
	MakeSale(string, string, float32, float32) (string, error)
	RefundSale(string) (string, error)
}

// MakeOrderReq is the information needed to make an order
type MakeOrderReq struct {
	PostID              int64
	NumServings         int32
	BTNonce             string
	ExchangeMethod      types.ExchangeMethods
	GigamuncherAddress  types.Address
	GigamuncherID       string
	GigamuncherName     string
	GigamuncherPhotoURL string
}

func (req *MakeOrderReq) valid() error {
	if req.PostID == 0 {
		return errInvalidParameter.WithMessage("Invalid post id.")
	}
	if req.NumServings <= 0 {
		return errInvalidParameter.WithMessage("Invalid number of servings.")
	}
	if req.BTNonce == "" {
		return errInvalidParameter.WithMessage("Invalid payment nonce.")
	}
	if req.ExchangeMethod.Pickup() && req.ExchangeMethod.Delivery() ||
		!req.ExchangeMethod.Pickup() && req.ExchangeMethod.Delivery() {
		return errInvalidParameter.WithMessage("Invalid exchange method. Must pick either pickup or delivery.")
	}
	if !req.GigamuncherAddress.GeoPoint.Valid() {
		return errInvalidParameter.WithMessage("Invalid location.")
	}
	return nil
}

// MakeOrder makes an order for the user using the nonce as payment
func (c Client) MakeOrder(req *MakeOrderReq) (*Resp, error) {
	err := req.valid()
	if err != nil {
		return nil, errors.Wrap("make order request is invalid", err)
	}
	postC := post.New(c.ctx)
	paymentC := payment.New(c.ctx)
	return makeOrder(c.ctx, req, postC, paymentC)
}

func makeOrder(ctx context.Context, req *MakeOrderReq, postC postClient, paymentC paymentClient) (*Resp, error) {
	// get post
	orderInfo, err := postC.GetOrderInfo(req.PostID)
	if err != nil {
		return nil, errors.Wrap("cannot get order info", err)
	}
	// check if numServings exist
	if req.NumServings > orderInfo.ServingsOffered-orderInfo.NumServingsOrdered {
		return nil, errNotEnoughServings
	}
	var exchangeMethod types.ExchangeMethods
	var totalExchangePrice float32
	distance := req.GigamuncherAddress.GreatCircleDistance(orderInfo.GigachefAddress.GeoPoint)
	duration := req.GigamuncherAddress.EstimatedDuration(orderInfo.GigachefAddress.GeoPoint)
	if req.ExchangeMethod.Delivery() {
		// TODO find cheapest delivery method
		if req.ExchangeMethod.ChefDelivery() {
			exchangeMethod.SetChefDelivery(true)
			totalExchangePrice = orderInfo.GigachefDelivery.Price
		} else {
			return nil, errInvalidParameter.WithMessage("Invalid delivery method.")
		}

		// TODO schedule delivery stuff
	} else { // pickup
		totalExchangePrice = 0
		exchangeMethod.SetPickup(true)
	}

	totalPricePerServing := orderInfo.PricePerServing * float32(req.NumServings)
	gigaFee := orderInfo.PricePerServing - orderInfo.ChefPricePerServing
	totalGigaFee := gigaFee * float32(req.NumServings)
	totalTax := (orderInfo.TaxPercentage / 100) * (totalGigaFee + totalPricePerServing + totalExchangePrice)
	totalFeeAndTax := totalGigaFee + totalTax
	totalPrice := totalPricePerServing + totalFeeAndTax + totalExchangePrice

	// submit braintree info
	var transactionID string
	transactionID, err = paymentC.MakeSale(orderInfo.BTSubMerchantID, req.BTNonce, totalPrice, totalFeeAndTax)
	if err != nil {
		return nil, errors.Wrap("cannot make BT sale", err)
	}
	order := &Order{
		CreatedDateTime: time.Now().UTC(),
		// ExpectedExchangeDataTime: orderInfo TODO
		BasicOrderIDs: BasicOrderIDs{
			GigachefID:    orderInfo.GigachefID,
			GigamuncherID: req.GigamuncherID,
			PostID:        req.PostID,
			ItemID:        orderInfo.ItemID,
		},
		PostTitle:           orderInfo.PostTitle,
		PostPhotoURL:        orderInfo.PostPhotoURL,
		GigamuncherName:     req.GigamuncherName,
		GigamuncherPhotoURL: req.GigamuncherPhotoURL,
		ChefPricePerServing: orderInfo.ChefPricePerServing,
		Servings:            req.NumServings,
		PaymentInfo: PaymentInfo{
			BTTransactionID: transactionID,
			Price:           totalPricePerServing,
			ExchangePrice:   totalExchangePrice,
			GigaFee:         totalGigaFee,
			TaxPrice:        totalTax,
			TotalPrice:      totalPrice,
		},
		ExchangeMethod: exchangeMethod,
		ExchangePlanInfo: exchangePlanInfo{
			GigachefAddress:    orderInfo.GigachefAddress,
			GigamuncherAddress: req.GigamuncherAddress,
			Distance:           distance,
			Duration:           duration,
		},
	}
	// create the order
	orderID, err := createOrder(ctx, order, postC)
	if err != nil {
		_, tErr := paymentC.RefundSale(transactionID)
		if tErr != nil {
			utils.Criticalf(ctx, "BT Transaction (%s) was not voided! Err: %+v", transactionID, tErr)
		}
		return nil, errors.Wrap("cannot run in transaction", err)
	}
	resp := &Resp{
		ID:    orderID,
		Order: *order,
	}
	return resp, nil
}

// func CancleOrder

func createOrder(ctx context.Context, order *Order, postC postClient) (int64, error) {
	var orderID int64
	var err error
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		// create order
		orderID, err = putIncomplete(ctx, order)
		if err != nil {
			return errDatastore.WithError(err).Wrap("cannot put incomplete order")
		}
		// add order to post
		addOrderReq := &post.AddOrderReq{
			PostID:              order.PostID,
			OrderID:             orderID,
			ExchangeMethod:      order.ExchangeMethod,
			ExchangeDuration:    order.ExchangePlanInfo.Duration,
			GigamuncherGeopoint: order.ExchangePlanInfo.GigamuncherAddress.GeoPoint,
			Servings:            order.Servings,
			GigamuncherID:       order.GigamuncherID,
		}
		err = postC.AddOrder(ctx, addOrderReq)
		if err != nil {
			return errors.Wrap("cannot add order to post", err)
		}
		return nil // success
	}, nil)
	return orderID, err
}
