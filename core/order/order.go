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
	errOrderIsClosed     = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Order is already closed."}
	errInvalidParameter  = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errUnauthorized      = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access."}
)

// Client is a client for orders
type Client struct {
	ctx context.Context
}

// New returns a new Client
func New(ctx context.Context) Client {
	return Client{ctx: ctx}
}

// Resp contains the order id and order itself
type Resp struct {
	ID    int64
	Order // embedded
}

type postClient interface {
	GetOrderInfo(int64) (*post.OrderInfoResp, error)
	AddOrder(context.Context, *post.AddOrderReq) error
	RemoveOrder(context.Context, int64) error
}

type paymentClient interface {
	MakeSale(string, string, float32, float32) (string, error)
	RefundSale(string) (string, error)
}

// CancelOrder cancels an order and refunds the money
func (c Client) CancelOrder(userID string, orderID int64) (*Resp, error) {
	if userID == "" {
		return nil, errInvalidParameter.WithMessage("Invalid user id.")
	}
	if orderID == 0 {
		return nil, errInvalidParameter.WithMessage("Invalid order id.")
	}
	postC := post.New(c.ctx)
	paymentC := payment.New(c.ctx)
	return cancelOrder(c.ctx, userID, orderID, postC, paymentC)
}

func cancelOrder(ctx context.Context, userID string, orderID int64, postC postClient, paymentC paymentClient) (*Resp, error) {
	// get the order
	order := new(Order)
	err := get(ctx, orderID, order)
	if err != nil {
		return nil, errDatastore.WithError(err).WithMessage("cannot get order")
	}
	// check if user is chef or muncher
	if order.GigamuncherID != userID && order.GigachefID != userID {
		return nil, errUnauthorized.Wrap("user is not part of order")
	}
	err = datastore.RunInTransaction(ctx, func(tc context.Context) error {
		// remove from post
		err = postC.RemoveOrder(tc, orderID)
		if err != nil {
			return errors.Wrap("failed to remove order form post", err)
		}
		// refund payment
		var transactionID string
		transactionID, err = paymentC.RefundSale(order.PaymentInfo.BTTransactionID)
		if err != nil {
			return errors.Wrap("failed to refund sale", err)
		}
		// change order to cancel
		order.State = State.Refunded
		order.BTRefundTransactionID = transactionID
		if userID == order.GigamuncherID {
			order.GigamuncherCanceled = true
		} else { // userID == order.GigachefID
			order.GigachefCanceled = true
		}
		err = put(tc, orderID, order)
		if err != nil {
			return errDatastore.WithError(err).Wrap("cannot put order")
		}
		return nil
	}, nil)
	if err != nil {
		return nil, errors.Wrap("cannot run in transaction", err)
	}
	resp := &Resp{
		ID:    orderID,
		Order: *order,
	}
	return resp, nil
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
	// check if post is closed
	if time.Now().UTC().After(orderInfo.ClosingDateTime) {
		return nil, errOrderIsClosed
	}
	var exchangeMethod types.ExchangeMethods
	var totalExchangePrice float32
	var expectedExchangeTime time.Time
	distance := req.GigamuncherAddress.GreatCircleDistance(orderInfo.GigachefAddress.GeoPoint)
	duration := req.GigamuncherAddress.EstimatedDuration(orderInfo.GigachefAddress.GeoPoint)
	if req.ExchangeMethod.Delivery() {
		// TODO find cheapest delivery method
		if req.ExchangeMethod.ChefDelivery() {
			exchangeMethod.SetChefDelivery(true)
			totalExchangePrice = orderInfo.GigachefDelivery.Price
			expectedExchangeTime = orderInfo.ReadyDateTime
		} else {
			return nil, errInvalidParameter.WithMessage("Invalid delivery method.")
		}

		// TODO schedule delivery stuff
	} else { // pickup
		totalExchangePrice = 0
		exchangeMethod.SetPickup(true)
		expectedExchangeTime = orderInfo.ReadyDateTime
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
		CreatedDateTime:          time.Now(),
		ExpectedExchangeDataTime: expectedExchangeTime,
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

func createOrder(ctx context.Context, order *Order, postC postClient) (int64, error) {
	var orderID int64
	var err error
	err = datastore.RunInTransaction(ctx, func(tc context.Context) error {
		// create order
		orderID, err = putIncomplete(tc, order)
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
		err = postC.AddOrder(tc, addOrderReq)
		if err != nil {
			return errors.Wrap("cannot add order to post", err)
		}
		return nil // success
	}, nil)
	return orderID, err
}
