package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	pbshared "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/gorilla/schema"
)

// GetHasSubscribed gets all subscribers.
func GetHasSubscribed(ctx context.Context, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetHasSubscribedReq)
	var err error
	// decode request
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(req, r.URL.Query())
		if err != nil {
			return failedToDecode(err)
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			return failedToDecode(err)
		}
		defer closeRequestBody(r)
	}
	logging.Infof(ctx, "Request: %+v", req)
	// end decode request
	date, err := getTime(req.Date)
	if err != nil {
		return errors.Annotate(err, "failed to decode date")
	}

	subC := subold.New(ctx)
	subscribers, err := subC.GetHasSubscribed(date)

	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get all subscribers")
	}

	resp := &pb.GetHasSubscribedResp{
		Subscribers: pbSubscribers(subscribers),
	}

	return resp
}

func pbSubscribers(subscribers []subold.SubscriptionSignUp) []*pb.Subscriber {
	sls := make([]*pb.Subscriber, len(subscribers))
	for i := range subscribers {
		sls[i] = pbSubscriber(&subscribers[i])
	}
	return sls
}

func pbSubscriber(subscriber *subold.SubscriptionSignUp) *pb.Subscriber {
	return &pb.Subscriber{
		Email:              subscriber.Email,
		Date:               subscriber.Date.Format(time.RFC3339),
		Name:               subscriber.Name,
		FirstName:          subscriber.FirstName,
		LastName:           subscriber.LastName,
		Address:            pbAddress(&subscriber.Address),
		CustomerId:         subscriber.CustomerID,
		SubscriptionIds:    subscriber.SubscriptionIDs,
		FirstPaymentDate:   subscriber.FirstPaymentDate.Format(time.RFC3339),
		IsSubscribed:       subscriber.IsSubscribed,
		SubscriptionDate:   subscriber.SubscriptionDate.Format(time.RFC3339),
		UnsubscribedDate:   subscriber.UnSubscribedDate.Format(time.RFC3339),
		FirstBoxDate:       subscriber.FirstBoxDate.Format(time.RFC3339),
		Servings:           int32(subscriber.Servings),
		VegetarianServings: int32(subscriber.VegetarianServings),
		DeliveryTime:       int32(subscriber.DeliveryTime),
		SubscriptionDay:    subscriber.SubscriptionDay,
		WeeklyAmount:       subscriber.WeeklyAmount,
		PaymentMethodToken: subscriber.PaymentMethodToken,
		Reference:          subscriber.Reference,
		PhoneNumber:        subscriber.PhoneNumber,
		DeliveryTips:       subscriber.DeliveryTips,
		BagReminderSms:     subscriber.BagReminderSMS,
		// gift
		NumGiftDinners: int32(subscriber.NumGiftDinners),
		ReferenceEmail: subscriber.ReferenceEmail,
		GiftRevealDate: subscriber.GiftRevealDate.Format(time.RFC3339),
		// stats
		ReferralPageOpens: int32(subscriber.ReferralPageOpens),
		ReferredPageOpens: int32(subscriber.ReferredPageOpens),
		GiftPageOpens:     int32(subscriber.GiftPageOpens),
		GiftedPageOpens:   int32(subscriber.GiftedPageOpens),
	}
}

func pbAddress(address *types.Address) *pbshared.Address {
	return &pbshared.Address{
		Apt:         address.APT,
		Street:      address.Street,
		City:        address.City,
		State:       address.State,
		Zip:         address.Zip,
		Country:     address.Country,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
		FullAddress: address.String(),
	}
}
