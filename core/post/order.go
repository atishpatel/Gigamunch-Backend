package post

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/item"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/maps"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/order"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	chefExchangeMinutes = 5
)

var (
	errNotEnoughServings = errors.ErrorWithCode{Code: errors.CodeNotEnoughServingsLeft, Message: "Not enough servings left."}
	errPostIsClosed      = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Post is already closed."}
	errDelivery          = errors.ErrorWithCode{Code: errors.CodeDeliveryMethodNotAvaliable, Message: "Delivery is no longer avaliable."}
)

type orderClient interface {
	Create(context.Context, *order.CreateReq) (int64, *order.Order, error)
	Cancel(context.Context, string, int64) (*order.Resp, error)
	GetPostID(int64) (int64, error)
}

// MakeOrderReq is the information needed to make an order
type MakeOrderReq struct {
	PostID              int64
	NumServings         int32
	PaymentNonce        string
	ExchangeWindowIndex int32
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
	if req.ExchangeWindowIndex < 0 {
		return errInvalidParameter.WithMessage("Invalid pickup or delivery method.")
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
	itemC := item.New(c.ctx)
	return makeOrder(c.ctx, req, orderC, itemC)
}

func makeOrder(ctx context.Context, req *MakeOrderReq, orderC orderClient, itemC itemClient) (*order.Resp, error) {
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
	if int(req.ExchangeWindowIndex) >= len(p.ExchangeTimes) {
		return nil, errInvalidParameter.WithMessage("Pickup or Delivery window is out of range.")
	}
	if time.Now().After(p.ClosingDateTime) || time.Now().After(p.ExchangeTimes[req.ExchangeWindowIndex].StartDateTime) {
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
			exchangeMethod.SetChefDelivery(true)
			if !p.ExchangeTimes[req.ExchangeWindowIndex].AvailableExchangeMethods.ChefDelivery() {
				return nil, errDelivery.Wrap("Chef delivery is not an avaliable option")
			}
			chefDeliveryPoints := []types.GeoPoint{req.GigamuncherAddress.GeoPoint}
			// find all the waypoints for the time window with the same exchange method
			for _, order := range p.Orders {
				if order.ExchangeWindowIndex == req.ExchangeWindowIndex && order.ExchangeMethod.Equal(exchangeMethod) {
					chefDeliveryPoints = append(chefDeliveryPoints, order.GigamuncherGeopoint)
				}
			}
			// get chef delivery route info for this order
			arrivalTimes, optimalRoute, err := maps.GetDirections(ctx, p.ExchangeTimes[req.ExchangeWindowIndex].StartDateTime, p.GigachefAddress.GeoPoint, chefDeliveryPoints)
			if err != nil {
				return nil, errors.Wrap("cannot maps.GetDirections", err)
			}
			var currentOrderRouteIndex int
			precedingWaypoints := 0
			for i, v := range optimalRoute {
				if v == 0 {
					currentOrderRouteIndex = i
				}
				precedingWaypoints++
			}
			// update if exchange window is still avaliable for delivery
			numDeliveries := len(chefDeliveryPoints)
			deliveriesTimeBuffer := time.Duration(numDeliveries*chefExchangeMinutes) * time.Minute
			lastDeliveryArrivalTime := arrivalTimes[len(arrivalTimes)-1].Add(deliveriesTimeBuffer)
			if lastDeliveryArrivalTime.After(p.ExchangeTimes[req.ExchangeWindowIndex].EndDateTime) {
				p.ExchangeTimes[req.ExchangeWindowIndex].AvailableExchangeMethods.SetChefDelivery(false)
			}

			deliveryTimeBuffer := time.Duration(chefExchangeMinutes*precedingWaypoints) * time.Minute // chef exchanging the food with person
			expectedExchangeTime = arrivalTimes[currentOrderRouteIndex].Add(deliveryTimeBuffer)       // TODO update all expectedExchangeTimes
			exchangeCost = p.GigachefDelivery.BasePrice
		} else {
			// we don't support yet
			return nil, errDelivery.Wrap("chosen delivery method is not an option")
		}
	} else { // pickup
		if !p.ExchangeTimes[req.ExchangeWindowIndex].AvailableExchangeMethods.Pickup() {
			return nil, errDelivery.WithMessage("Pickup is not avaliable.").Wrap("pickup is not an avaliable option")
		}
		exchangeCost = 0
		exchangeMethod.SetPickup(true)
		expectedExchangeTime = p.ExchangeTimes[req.ExchangeWindowIndex].StartDateTime
	}
	// run in transaction
	opts := &datastore.TransactionOptions{XG: true, Attempts: 1}
	var order *order.Resp
	err = datastore.RunInTransaction(ctx, func(tc context.Context) error {
		// make order
		createOrderReq := createOrderRequest(p, req, exchangeMethod, exchangeCost, expectedExchangeTime, distance, duration)
		orderID, _, err := orderC.Create(tc, createOrderReq)
		if err != nil {
			return errors.Wrap("cannot create order", err)
		}
		// add order to post
		p.NumServingsOrdered += req.NumServings
		if req.ExchangeMethod.ChefDelivery() {
			// TODO
		}
		pOrder := OrderPost{
			OrderID:             orderID,
			GigamuncherID:       req.GigamuncherID,
			GigamuncherName:     req.GigamuncherName,
			GigamuncherPhotoURL: req.GigamuncherPhotoURL,
			GigamuncherGeopoint: req.GigamuncherAddress.GeoPoint,
			ExchangeWindowIndex: req.ExchangeWindowIndex,
			ExchangeTime:        createOrderReq.ExpectedExchangeTime,
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
	err = itemC.AddNumTotalOrders(p.ItemID, req.NumServings)
	if err != nil {
		utils.Errorf(ctx, "failed to itemC.AddNumTotalOrders for itemID(%d) by %d : %s", p.ItemID, req.NumServings, err)
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
		ReadyDateTime:         expectedExchangeTime,
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
	itemC := item.New(c.ctx)
	return cancelOrder(c.ctx, userID, orderID, orderC, itemC)
}

func cancelOrder(ctx context.Context, userID string, orderID int64, orderC orderClient, itemC itemClient) (*order.Resp, error) {
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

	var orderIndex int
	orderIndex, err = findOrderIndex(orderID, p)
	if err != nil {
		return nil, err
	}
	// check if it's too late to cancel the order
	if time.Now().After(p.ClosingDateTime) || time.Now().After(p.ExchangeTimes[p.Orders[orderIndex].ExchangeWindowIndex].StartDateTime) {
		return nil, errPostIsClosed
	}
	if p.GigachefID != userID && p.Orders[orderIndex].GigamuncherID != userID {
		return nil, errUnauthorized.Wrap("user is not part of order or post.")
	}
	numServings := p.Orders[orderIndex].Servings
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
			p.ExchangeTimes[p.Orders[orderIndex].ExchangeWindowIndex].AvailableExchangeMethods.SetChefDelivery(true)
		}
		p.NumServingsOrdered -= p.Orders[orderIndex].Servings
		if orderIndex == 0 {
			p.Orders = p.Orders[orderIndex+1:]
		} else {
			p.Orders = append(p.Orders[:orderIndex], p.Orders[orderIndex+1:]...)
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
	err = itemC.AddNumTotalOrders(p.ItemID, -numServings)
	if err != nil {
		utils.Errorf(ctx, "failed to itemC.AddNumTotalOrders for itemID(%d) by -%d : %s", p.ItemID, numServings, err)
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
