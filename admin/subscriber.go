package admin

import (
	"context"
	"net/http"
	"strings"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	pbshared "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

// SendCustomerSMS sends an CustomerSMS from Gigamunch to number.
func (s *server) SendCustomerSMS(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.SendCustomerSMSReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	dilm := "{{name}}"
	firstNameDilm := "{{first_name}}"
	msg := req.Message

	messageC := message.New(ctx)
	var emails []string
	for _, email := range req.Emails {
		if email != "" {
			emails = append(emails, email)
		}
	}
	subC := subold.New(ctx)
	subs, err := subC.GetSubscribers(emails)
	if err != nil {
		return errors.Annotate(err, "failed to subold.GetSubscribers")
	}
	for _, s := range subs {
		if s.PhoneNumber == "" {
			continue
		}
		name := s.FirstName
		if name == "" {
			name = s.Name
		}
		name = strings.Title(name)
		msg = strings.Replace(msg, dilm, name, -1)
		msg = strings.Replace(msg, firstNameDilm, s.FirstName, -1)
		err = messageC.SendDeliverySMS(s.PhoneNumber, msg)
		if err != nil {
			return errors.Annotate(err, "failed ot message.SendSMS")
		}
		// log
		if log != nil {
			payload := &logging.MessagePayload{
				Platform: "SMS",
				Body:     msg,
				From:     "Gigamunch",
				To:       s.PhoneNumber,
			}
			log.SubMessage(0, s.Email, payload)
		}
	}
	return nil
}

// GetSubscriber gets all info about a subscriber from their email address
func (s *server) GetSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetSubscriberReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	email := req.Email

	subC := subold.New(ctx)
	subscriber, err := subC.GetSubscriber(email)

	if err != nil {
		return errors.GetErrorWithCode(err).Annotate("failed to get subscriber")
	}

	resp := &pb.GetSubscriberResp{
		Subscriber: pbSubscriber(subscriber),
	}

	return resp
}

// GetHasSubscribed gets all subscribers.
func (s *server) GetHasSubscribed(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.GetHasSubscribedReq)
	var err error

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request

	date := getDatetime(req.Date)

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
