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
	AvaliableExchangeMethods types.ExchangeMethods
	ClosingDateTime          time.Time
	TotalDeliveryTime        int64
}

// GetOrderInfo returns information related to an order
func GetOrderInfo(ctx context.Context, postID int64) (*OrderInfoResp, error) {
	if postID == 0 {
		return nil, errInvalidParameter.WithMessage("Invalid post id.")
	}
	p := new(Post)
	err := get(ctx, postID, p)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get Post")
	}
	var photoURL string
	if len(p.Photos) > 0 {
		photoURL = p.Photos[0]
	}
	var totalDeliveryTime int64
	for i := range p.Orders {
		totalDeliveryTime += p.Orders[i].DeliveryTime
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
		AvaliableExchangeMethods: p.AvaliableExchangeMethods,
		ClosingDateTime:          p.ClosingDateTime,
		TotalDeliveryTime:        totalDeliveryTime,
	}
	return resp, nil
}

// AddOrderReq adds an order to a post
type AddOrderReq struct {
	PostID         int64
	OrderID        int64
	ExchangeMethod types.ExchangeMethods
	DeliveryTime   int64
	Servings       int32
	GigamuncherID  string
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
	return nil
}

// AddOrder adds an order to the list of order for a post
func AddOrder(ctx context.Context, req *AddOrderReq) error {
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
	pOrder := postOrder{
		OrderID:        req.OrderID,
		GigamuncherID:  req.GigamuncherID,
		ExchangeMethod: req.ExchangeMethod,
		DeliveryTime:   req.DeliveryTime,
		Servings:       req.Servings,
	}
	p.Orders = append(p.Orders, pOrder)
	err = put(ctx, req.PostID, p)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put post")
	}
	return nil
}
