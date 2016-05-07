package post

import (
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

// OrderInfoResp contains all the information needed related to an order
type OrderInfoResp struct {
	GigachefID               string
	ItemID                   int64
	PostTitle                string
	PostPhotoURL             string
	BTSubMerchantID          string
	ServingsOffered          int32
	NumServingsOrdered       int32
	ChefPricePerServing      float32
	PricePerServing          float32
	TaxPercentage            float32
	AvailableExchangeMethods types.ExchangeMethods
	ClosingDateTime          time.Time
	ReadyDateTime            time.Time
	GigachefDelivery         GigachefDelivery
	GigachefAddress          types.Address
}

// GetOrderInfo returns information related to an order
func (c Client) GetOrderInfo(postID int64) (*OrderInfoResp, error) {
	if postID == 0 {
		return nil, errInvalidParameter.WithMessage("Invalid post id.")
	}
	p := new(Post)
	err := get(c.ctx, postID, p)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get Post")
	}
	var photoURL string
	if len(p.Photos) > 0 {
		photoURL = p.Photos[0]
	}
	resp := &OrderInfoResp{
		GigachefID:               p.GigachefID,
		ItemID:                   p.ItemID,
		PostTitle:                p.Title,
		PostPhotoURL:             photoURL,
		BTSubMerchantID:          p.BTSubMerchantID,
		ServingsOffered:          p.ServingsOffered,
		NumServingsOrdered:       p.NumServingsOrdered,
		ChefPricePerServing:      p.ChefPricePerServing,
		PricePerServing:          p.PricePerServing,
		TaxPercentage:            p.TaxPercentage,
		AvailableExchangeMethods: p.AvailableExchangeMethods,
		ClosingDateTime:          p.ClosingDateTime,
		ReadyDateTime:            p.ReadyDateTime,
		GigachefDelivery:         p.GigachefDelivery,
		GigachefAddress:          p.GigachefAddress,
	}
	return resp, nil
}

// AddOrderReq adds an order to a post
type AddOrderReq struct {
	PostID              int64
	OrderID             int64
	GigamuncherID       string
	ExchangeMethod      types.ExchangeMethods
	ExchangeDuration    int64
	GigamuncherGeopoint types.GeoPoint
	Servings            int32
}

func (req *AddOrderReq) valid() error {
	if req.PostID <= 0 {
		return errInvalidParameter.WithMessage("Invalid post id.")
	}
	if req.OrderID <= 0 {
		return errInvalidParameter.WithMessage("Invalid order id.")
	}
	if req.Servings <= 0 {
		return errInvalidParameter.WithMessage("Invalid number of servings.")
	}
	if req.GigamuncherID == "" {
		return errInvalidParameter.WithMessage("Invalid gigamuncher ID.")
	}
	if !req.GigamuncherGeopoint.Valid() {
		return errInvalidParameter.WithMessage("Invalid location.")
	}
	return nil
}

// AddOrder adds an order to the list of order for a post
func (c Client) AddOrder(ctx context.Context, req *AddOrderReq) error {
	err := req.valid()
	if err != nil {
		return errors.Wrap("order request is invalid", err)
	}
	p := new(Post)
	err = get(ctx, req.PostID, p)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot get post")
	}
	if req.Servings > (p.ServingsOffered - p.NumServingsOrdered) {
		return errNotEnoughServings
	}
	p.NumServingsOrdered = p.NumServingsOrdered + req.Servings
	if req.ExchangeMethod.ChefDelivery() {

		p.GigachefDelivery.TotalDuration += req.ExchangeDuration
		// TODO reculcate GigachefDelivery.TotalDuration
		// p.GigachefDelivery.TotalDuration = maps.GetTotalTime(origins, destinations)
	}
	pOrder := postOrder{
		OrderID:             req.OrderID,
		GigamuncherID:       req.GigamuncherID,
		GigamuncherGeopoint: req.GigamuncherGeopoint,
		ExchangeMethod:      req.ExchangeMethod,
		Servings:            req.Servings,
	}
	p.Orders = append(p.Orders, pOrder)
	err = put(c.ctx, req.PostID, p)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put post")
	}
	return nil
}

// RemoveOrder removes an order from the post
func (c Client) RemoveOrder(ctx context.Context, id int64) error {
	p := new(Post)
	err := get(ctx, id, p)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot get post")
	}

	var found bool
	for i := range p.Orders {
		if p.Orders[i].OrderID == id {
			found = true
			p.NumServingsOrdered += p.Orders[i].Servings
			if p.Orders[i].ExchangeMethod.ChefDelivery() {
				// remove chef delivery duration
				// TODO reculcate GigachefDelivery.TotalDuration
				// p.GigachefDelivery.TotalDuration = maps.GetTotalTime(origins, destinations)
			}
			if i == 0 {
				p.Orders = p.Orders[i+1:]
			} else {
				p.Orders = append(p.Orders[:i-1], p.Orders[i+1:]...)
			}
			break
		}
	}
	if !found {
		return errInvalidParameter.Wrap("order id not in post orders")
	}
	err = put(ctx, id, p)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put post")
	}
	return nil
}
