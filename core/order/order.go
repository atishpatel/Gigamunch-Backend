package order

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

var (
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errUnauthorized     = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access."}
)

// Client is a client for orders
type Client struct {
	ctx context.Context
}

// New returns a new Client
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Resp contains the order id and order itself
type Resp struct {
	ID    int64
	Order // embedded
}

type paymentClient interface {
	MakeSale(string, string, float32, float32) (string, error)
	RefundSale(string) (string, error)
	ReleaseSale(string) (string, error)
	CancelRelease(string) (string, error)
}

// CreateReq is the request needed to create an order
type CreateReq struct {
	PostID                int64
	NumServings           int32
	PaymentNonce          string
	GigamuncherAddress    types.Address
	GigamuncherID         string
	GigamuncherName       string
	GigamuncherPhotoURL   string
	GigachefID            string
	GigachefSubMerchantID string
	ItemID                int64
	PostTitle             string
	PostPhotoURL          string
	ChefPricePerServing   float32
	PricePerServing       float32
	TaxPercentage         float32
	ExchangeMethod        types.ExchangeMethods
	ExchangePrice         float32
	ExpectedExchangeTime  time.Time
	CloseDateTime         time.Time
	ReadyDateTime         time.Time
	GigachefAddress       types.Address
	Distance              float32
	Duration              int64
}

func (req *CreateReq) valid() error {
	if req.PostID == 0 {
		return errInvalidParameter.WithMessage("Invalid post id.")
	}
	if req.NumServings <= 0 {
		return errInvalidParameter.WithMessage("Invalid number of servings.")
	}
	if req.PaymentNonce == "" {
		return errInvalidParameter.WithMessage("Invalid payment nonce.")
	}
	if req.ExchangeMethod.Pickup() && req.ExchangeMethod.Delivery() ||
		!req.ExchangeMethod.Pickup() && req.ExchangeMethod.Delivery() {
		return errInvalidParameter.WithMessage("Invalid exchange method. Must pick either pickup or delivery.")
	}
	if !req.GigamuncherAddress.GeoPoint.Valid() {
		return errInvalidParameter.WithMessage("Invalid muncher location.")
	}
	if !req.GigachefAddress.GeoPoint.Valid() {
		return errInvalidParameter.WithMessage("Invalid chef location.")
	}
	if req.PricePerServing == 0 {
		return errInvalidParameter.Wrap("price per serving cannot be 0.")
	}
	if req.ChefPricePerServing == 0 {
		return errInvalidParameter.Wrap("chef price per serving cannot be 0.")
	}
	return nil
}

// Create is used to save order information.
func (c *Client) Create(ctx context.Context, req *CreateReq) (*Resp, error) {
	// make and order
	err := req.valid()
	if err != nil {
		return nil, errors.Wrap("create order request is invalid", err)
	}
	paymentC := payment.New(ctx)
	return create(ctx, req, paymentC)
}

func create(ctx context.Context, req *CreateReq, paymentC paymentClient) (*Resp, error) {
	// TODO make two transactions if delivery cost is not 0
	totalPricePerServing := req.PricePerServing * float32(req.NumServings)
	gigaFee := req.PricePerServing - req.ChefPricePerServing
	totalGigaFee := gigaFee * float32(req.NumServings)
	totalExchangePrice := req.ExchangePrice
	totalTax := (req.TaxPercentage / 100) * (totalGigaFee + totalPricePerServing + totalExchangePrice)
	totalFeeAndTax := totalGigaFee + totalTax
	totalPrice := totalPricePerServing + totalFeeAndTax + totalExchangePrice

	transactionID, err := paymentC.MakeSale(req.GigachefSubMerchantID, req.PaymentNonce, totalPrice, totalFeeAndTax)
	if err != nil {
		return nil, errors.Wrap("cannot make sale", err)
	}
	order := &Order{
		CreatedDateTime: time.Now(),
		State:           State.Pending,
		ExpectedExchangeDateTime: req.ExpectedExchangeTime,
		PostCloseDateTime:        req.CloseDateTime,
		PostReadyDateTime:        req.ReadyDateTime,
		BasicOrderIDs: BasicOrderIDs{
			GigachefID:    req.GigachefID,
			GigamuncherID: req.GigamuncherID,
			PostID:        req.PostID,
			ItemID:        req.ItemID,
		},
		PostTitle:           req.PostTitle,
		PostPhotoURL:        req.PostPhotoURL,
		GigamuncherName:     req.GigamuncherName,
		GigamuncherPhotoURL: req.GigamuncherPhotoURL,
		ChefPricePerServing: req.ChefPricePerServing,
		PricePerServing:     req.PricePerServing,
		Servings:            req.NumServings,
		PaymentInfo: PaymentInfo{
			BTTransactionID: transactionID,
			Price:           totalPricePerServing,
			ExchangePrice:   req.ExchangePrice,
			GigaFee:         gigaFee,
			TaxPrice:        totalTax,
			TotalPrice:      totalPrice,
		},
		ExchangeMethod: req.ExchangeMethod,
		ExchangePlanInfo: exchangePlanInfo{
			GigachefAddress:    req.GigachefAddress,
			GigamuncherAddress: req.GigamuncherAddress,
			Distance:           req.Distance,
			Duration:           req.Duration,
		},
	}
	id, err := putIncomplete(ctx, order)
	if err != nil {
		_, pErr := paymentC.RefundSale(transactionID)
		if pErr != nil {
			utils.Criticalf(ctx, "BT Transaction (%s) was not voided! Err: %+v", transactionID, pErr)
		}
		return nil, errDatastore.WithError(err).Wrap("cannot put incomplete order")
	}
	resp := &Resp{
		ID:    id,
		Order: *order,
	}
	return resp, nil
}

// Cancel changes the state of an order to canceled. The userID determinds if
// the cancel is from the chef or muncher.
func (c *Client) Cancel(ctx context.Context, userID string, orderID int64) (*Resp, error) {
	if userID == "" {
		return nil, errInvalidParameter.WithMessage("Invalid user id.")
	}
	if orderID == 0 {
		return nil, errInvalidParameter.WithMessage("Invalid order id.")
	}
	paymenyC := payment.New(ctx)
	return cancel(ctx, userID, orderID, paymenyC)
}

func cancel(ctx context.Context, userID string, orderID int64, paymentC paymentClient) (*Resp, error) {
	// get the order
	order := new(Order)
	err := get(ctx, orderID, order)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get order")
	}
	// check if user is chef or muncher
	if order.GigamuncherID != userID && order.GigachefID != userID {
		return nil, errUnauthorized.Wrap("user is not part of order")
	}
	// refund payment
	var transactionID string
	transactionID, err = paymentC.RefundSale(order.PaymentInfo.BTTransactionID)
	if err != nil {
		return nil, errors.Wrap("cannot refund sale", err)
	}
	// change order to cancel
	order.State = State.Refunded
	order.BTRefundTransactionID = transactionID
	if userID == order.GigamuncherID {
		order.GigamuncherCanceled = true
	} else { // userID == order.GigachefID
		order.GigachefCanceled = true
	}
	err = put(ctx, orderID, order)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot put order")
	}
	resp := &Resp{
		ID:    orderID,
		Order: *order,
	}
	return resp, nil
}

// GetPostID returns the post id for an order
func (c *Client) GetPostID(orderID int64) (int64, error) {
	order := new(Order)
	err := get(c.ctx, orderID, order)
	if err != nil {
		return 0, errDatastore.WithError(err).Wrap("cannot get order")
	}
	return order.PostID, nil
}

// GetOrder gets an order
func (c *Client) GetOrder(userID string, orderID int64) (*Resp, error) {
	resp := new(Resp)
	if orderID == 0 {
		return resp, nil
	}
	o := new(Order)
	err := get(c.ctx, orderID, o)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get order")
	}
	if userID != o.GigamuncherID && userID != o.GigachefID {
		return nil, errUnauthorized.Wrapf("user(%s) does not have access to order(%d)", userID, orderID)
	}
	resp.ID = orderID
	resp.Order = *o
	return resp, nil
}

// UpdateReviewID updates the review id for an order
func (c *Client) UpdateReviewID(userID string, orderID int64, reviewID int64) (*Resp, error) {
	resp := new(Resp)
	if orderID == 0 {
		return resp, nil
	}
	o := new(Order)
	err := get(c.ctx, orderID, o)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get order(%d)", orderID)
	}
	if userID != o.GigamuncherID && userID != o.GigachefID {
		return nil, errUnauthorized.Wrapf("user(%s) does not have access to order(%d)", userID, orderID)
	}
	o.ReviewID = reviewID
	err = put(c.ctx, orderID, o)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get order(%d)", orderID)
	}
	resp.ID = orderID
	resp.Order = *o
	return resp, nil
}

// GetOrders gets the orders
func (c *Client) GetOrders(muncherID string, limit *types.Limit) ([]Resp, error) {
	ids, orders, err := getSortedOrders(c.ctx, muncherID, limit.Start, limit.End)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get sorted orders")
	}

	resps := make([]Resp, len(ids))
	for i := range ids {
		resps[i].ID = ids[i]
		resps[i].Order = orders[i]
	}

	return resps, nil
}

// Process will process the order
func (c *Client) Process(id int64) (*Resp, error) {
	resp := &Resp{ID: id}
	o := new(Order)
	err := get(c.ctx, id, o)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get order(%d)", id)
	}
	switch o.State {
	case State.Canceled:
	case State.Refunded:
	case State.Paid:
	case State.Issues:
		resp.Order = *o
		return resp, nil
	case State.Pending:
		paymentC := payment.New(c.ctx)
		return payOrder(c.ctx, id, o, paymentC)
	}
	// default
	return nil, errInvalidParameter.WithMessage("order is in an unexpected state").Wrapf("order(%d) is unkown state(%s)", id, o.State)
}

func payOrder(ctx context.Context, id int64, o *Order, paymentC paymentClient) (*Resp, error) {
	if o.State != State.Pending {
		return nil, errInvalidParameter.WithMessage("order isn't in pending state").Wrapf("order(%d) isn't in pending state. state(%s)", id, o.State)
	}
	_, err := paymentC.ReleaseSale(o.PaymentInfo.BTTransactionID)
	if err != nil {
		return nil, errors.Wrap("failed to release sale", err)
	}
	o.State = State.Paid
	err = put(ctx, id, o)
	if err != nil {
		_, btErr := paymentC.CancelRelease(o.PaymentInfo.BTTransactionID)
		if btErr != nil {
			utils.Criticalf(ctx, "failed to cancel release of transaction(%s): err: %#v", o.PaymentInfo.BTTransactionID, btErr)
			return nil, errors.Wrap("failed to cancel release", btErr)
		}
		return nil, errDatastore.WithError(err).Wrapf("cannot put order(%d)", id)
	}
	resp := &Resp{
		ID:    id,
		Order: *o,
	}
	return resp, nil
}

// SetStateToPendingByTransactionID sets the state to pending
func (c *Client) SetStateToPendingByTransactionID(transactionID string) (int64, error) {
	id, o, err := getByTransactionID(c.ctx, transactionID)
	if err != nil {
		return 0, errDatastore.WithError(err).Wrapf("cannot getByTransactionID order(%d)", id)
	}
	o.State = State.Pending
	err = put(c.ctx, id, o)
	if err != nil {
		return 0, errDatastore.WithError(err).Wrapf("cannot put order(%d)", id)
	}
	return id, nil
}
