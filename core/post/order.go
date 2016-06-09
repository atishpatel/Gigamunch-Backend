package post

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/order"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	errNotEnoughServings = errors.ErrorWithCode{Code: errors.CodeNotEnoughServingsLeft, Message: "Not enough servings left."}
	errPostIsClosed      = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Post is already closed."}
	errDelivery          = errors.ErrorWithCode{Code: errors.CodeDeliveryMethodNotAvaliable, Message: "Delivery is no longer avaliable."}
)

type orderClient interface {
	Create(context.Context, *order.CreateReq) (*order.Resp, error)
	Cancel(context.Context, string, int64) (*order.Resp, error)
	GetPostID(int64) (int64, error)
}

// MakeOrderReq is the information needed to make an order
type MakeOrderReq struct {
	PostID              int64
	NumServings         int32
	PaymentNonce        string
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
	if req.PaymentNonce == "" {
		return errInvalidParameter.WithMessage("Invalid payment nonce.")
	}
	if req.ExchangeMethod.Pickup() && req.ExchangeMethod.Delivery() ||
		!req.ExchangeMethod.Pickup() && !req.ExchangeMethod.Delivery() {
		return errInvalidParameter.WithMessage("Invalid exchange method. Must pick either pickup or delivery.")
	}
	if !req.GigamuncherAddress.GeoPoint.Valid() {
		return errInvalidParameter.WithMessage("Invalid location.")
	}
	return nil
}

// MakeOrder makes and order for the post
func (c Client) MakeOrder(req *MakeOrderReq) (*order.Resp, error) {
	err := req.valid()
	if err != nil {
		return nil, errors.Wrap("make order request is invalid", err)
	}
	orderC := order.New(c.ctx)

	return makeOrder(c.ctx, req, orderC)
}

func makeOrder(ctx context.Context, req *MakeOrderReq, orderC orderClient) (*order.Resp, error) {
	// get post stuff
	p := new(Post)
	err := get(ctx, req.PostID, p)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot get post(%d)", req.PostID)
	}
	// check if post stuff is avaliable
	if req.NumServings > p.ServingsOffered-p.NumServingsOrdered {
		return nil, errNotEnoughServings
	}
	if time.Now().After(p.ClosingDateTime) {
		return nil, errPostIsClosed
	}
	// calculate delivery info
	var exchangeMethod types.ExchangeMethods
	var exchangeCost float32
	var expectedExchangeTime time.Time
	distance := req.GigamuncherAddress.GreatCircleDistance(p.GigachefAddress.GeoPoint)
	duration := req.GigamuncherAddress.EstimatedDuration(p.GigachefAddress.GeoPoint)
	if req.ExchangeMethod.Delivery() {
		if req.ExchangeMethod.ChefDelivery() {
			if !p.AvailableExchangeMethods.ChefDelivery() {
				return nil, errDelivery.Wrap("chef delivery is not an avaliable option")
			}
			// TODO reculcate GigachefDelivery.TotalDuration
			// p.GigachefDelivery.TotalDuration = maps.GetTotalTime(origins, destinations)

			exchangeMethod.SetChefDelivery(true)
			expectedExchangeTime = p.ReadyDateTime.Add(time.Duration(duration) * time.Second)
			exchangeCost = p.GigachefDelivery.Price
		} else {
			// we don't support yet
			return nil, errDelivery.Wrap("chosen delivery method is not an option")
		}
	} else { // pickup
		if !p.AvailableExchangeMethods.Pickup() {
			return nil, errDelivery.WithMessage("Pickup is not avaliable.").Wrap("pickup is not an avaliable option")
		}
		exchangeCost = 0
		exchangeMethod.SetPickup(true)
		if p.IsOrderNow {
			expectedExchangeTime = time.Now().Add(time.Duration(p.EstimatedPreperationTime) * time.Second)
		} else {
			expectedExchangeTime = p.ReadyDateTime
		}
	}
	// run in transaction
	opts := &datastore.TransactionOptions{XG: true, Attempts: 1}
	var order *order.Resp
	err = datastore.RunInTransaction(ctx, func(tc context.Context) error {
		// make order
		createOrderReq := createOrderRequest(p, req, exchangeMethod, exchangeCost, expectedExchangeTime, distance, duration)
		order, err = orderC.Create(tc, createOrderReq)
		if err != nil {
			return errors.Wrap("cannot create order", err)
		}
		// add order to post
		p.NumServingsOrdered += req.NumServings
		if req.ExchangeMethod.ChefDelivery() {
			p.GigachefDelivery.TotalDuration += duration
		}
		pOrder := OrderPost{
			OrderID:             order.ID,
			GigamuncherID:       req.GigamuncherID,
			GigamuncherName:     req.GigamuncherName,
			GigamuncherPhotoURL: req.GigamuncherPhotoURL,
			GigamuncherGeopoint: req.GigamuncherAddress.GeoPoint,
			ExchangeTime:        order.ExpectedExchangeDateTime,
			ExchangeMethod:      exchangeMethod,
			Servings:            req.NumServings,
		}
		p.Orders = append(p.Orders, pOrder)
		// put order to post
		err = put(tc, req.PostID, p)
		if err != nil {
			_, cErr := orderC.Cancel(ctx, req.GigamuncherID, order.ID)
			if cErr != nil {
				utils.Criticalf(ctx, "BT Transaction (%s) was not voided! Err: %+v", order.PaymentInfo.BTTransactionID, cErr)
			}
			return errDatastore.WithError(err).Wrap("cannot put post")
		}
		return nil
	}, opts)
	if err != nil {
		return nil, errors.Wrap("cannot run in transaction", err)
	}
	return order, nil
}

func createOrderRequest(p *Post, req *MakeOrderReq, exchangeMethod types.ExchangeMethods,
	exchangeCost float32, expectedExchangeTime time.Time, distance float32, duration int64) *order.CreateReq {
	var photoURL string
	if len(p.Photos) > 0 {
		photoURL = p.Photos[0]
	}
	return &order.CreateReq{
		PostID:                req.PostID,
		NumServings:           req.NumServings,
		GigamuncherAddress:    req.GigamuncherAddress,
		GigamuncherID:         req.GigamuncherID,
		GigamuncherName:       req.GigamuncherName,
		GigamuncherPhotoURL:   req.GigamuncherPhotoURL,
		GigachefID:            p.GigachefID,
		GigachefSubMerchantID: p.BTSubMerchantID,
		ItemID:                p.ItemID,
		PostTitle:             p.Title,
		PostPhotoURL:          photoURL,
		ChefPricePerServing:   p.ChefPricePerServing,
		PricePerServing:       p.PricePerServing,
		TaxPercentage:         p.TaxPercentage,
		ExchangeMethod:        exchangeMethod,
		ExchangePrice:         exchangeCost,
		ExpectedExchangeTime:  expectedExchangeTime,
		CloseDateTime:         p.ClosingDateTime,
		ReadyDateTime:         p.ReadyDateTime,
		GigachefAddress:       p.GigachefAddress,
		Distance:              distance,
		Duration:              duration,
		PaymentNonce:          req.PaymentNonce,
	}
}

// CancelOrder removes the order from the post and cancels it
func (c Client) CancelOrder(userID string, orderID int64) (*order.Resp, error) {
	if userID == "" {
		return nil, errInvalidParameter.WithMessage("Invalid user id.")
	}
	if orderID == 0 {
		return nil, errInvalidParameter.WithMessage("Invalid order id.")
	}
	orderC := order.New(c.ctx)
	return cancelOrder(c.ctx, userID, orderID, orderC)
}

func cancelOrder(ctx context.Context, userID string, orderID int64, orderC orderClient) (*order.Resp, error) {
	postID, err := orderC.GetPostID(orderID)
	if err != nil {
		return nil, errors.Wrap("cannot get post id", err)
	}
	// get post
	p := new(Post)
	err = get(ctx, postID, p)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get post")
	}
	// check if order is cancelable
	if time.Now().After(p.ClosingDateTime) {
		return nil, errPostIsClosed
	}
	var orderIndex int
	orderIndex, err = findOrderIndex(orderID, p)
	if err != nil {
		return nil, err
	}
	if p.GigachefID != userID && p.Orders[orderIndex].GigamuncherID != userID {
		return nil, errUnauthorized.Wrap("user is not part of order or post.")
	}
	var order *order.Resp
	// run in transaction
	err = datastore.RunInTransaction(ctx, func(tc context.Context) error {
		// cancel order
		order, err = orderC.Cancel(ctx, userID, orderID)
		if err != nil {
			return errors.Wrap("cannot cancel order", err)
		}
		// remove order from post
		if p.Orders[orderIndex].ExchangeMethod.ChefDelivery() {
			// remove chef delivery duration
			// TODO reculcate GigachefDelivery.TotalDuration
			// p.GigachefDelivery.TotalDuration = maps.GetTotalTime(origins, destinations)
		}
		p.NumServingsOrdered -= p.Orders[orderIndex].Servings
		if orderIndex == 0 {
			p.Orders = p.Orders[orderIndex+1:]
		} else {
			p.Orders = append(p.Orders[:orderIndex-1], p.Orders[orderIndex+1:]...)
		}
		// save post
		err = put(ctx, postID, p)
		if err != nil {
			return errDatastore.WithError(err).Wrap("cannot put post")
		}
		return nil
	}, nil)
	if err != nil {
		return nil, errors.Wrap("cannot run in transaction", err)
	}
	return order, nil
}

func findOrderIndex(orderID int64, p *Post) (int, error) {
	for i := range p.Orders {
		if p.Orders[i].OrderID == orderID {
			return i, nil
		}
	}
	return 0, errInvalidParameter.Wrap("order id not in post orders")
}
