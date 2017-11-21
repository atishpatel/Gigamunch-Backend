package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/server"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/corenew/mail"
	"github.com/atishpatel/Gigamunch-Backend/corenew/maps"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var (
	errBadRequest = errors.BadRequestError
)

func addAPIRoutes(r *httprouter.Router) {
	http.HandleFunc("/api/v1/Login", handler(Login))
	http.HandleFunc("/api/v1/SubmitCheckout", handler(SubmitCheckout))
	http.HandleFunc("/api/v1/UpdatePayment", handler(UpdatePayment))
}

func validateSubmitCheckoutReq(r *pb.SubmitCheckoutReq) error {
	if r.Email == "" {
		return errInvalidParameter.WithMessage("Email address cannot be empty.").Annotate("no email address")
	}
	if !strings.Contains(r.Email, "@") {
		return errInvalidParameter.WithMessage("Email address must be an email.").Annotate("not email address")
	}
	if r.PaymentMethodNonce == "" {
		return errInvalidParameter.WithMessage("Invalid payment info.").Annotate("no payment nonce")
	}
	if r.FirstName == "" {
		return errInvalidParameter.WithMessage("First name must be provided.").Annotate("no first name")
	}
	return nil
}

// Login updates a user's payment.
func Login(ctx context.Context, r *http.Request) Response {
	req := new(pb.TokenOnlyReq)
	var err error
	// decode request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		return failedToDecode(err)
	}
	defer closeRequestBody(r)
	// end decode request
	resp := &pb.TokenOnlyResp{}

	_, authToken, err := auth.GetSessionWithGToken(ctx, req.Token)
	if err != nil {
		resp.Error = errors.GetSharedError(err)
		return resp
	}
	resp.Token = authToken
	return resp
}

// UpdatePayment updates a user's payment.
func UpdatePayment(ctx context.Context, r *http.Request) Response {
	req := new(pb.UpdatePaymentReq)
	var err error
	// decode request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		return failedToDecode(err)
	}
	defer closeRequestBody(r)
	// end decode request
	resp := &pb.ErrorOnlyResp{}

	key := datastore.NewKey(ctx, "ScheduleSignUp", req.Email, 0, nil)
	entry := &sub.SubscriptionSignUp{}
	err = datastore.Get(ctx, key, entry)
	if err == datastore.ErrNoSuchEntity {
		resp.Error = errBadRequest.WithMessage(fmt.Sprintf("Cannot find user with email: %s", req.Email)).Wrapf("failed to get ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		return resp
	}
	if err != nil {
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to get ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		return resp
	}
	paymentC := payment.New(ctx)
	paymentReq := &payment.CreateCustomerReq{
		CustomerID: entry.CustomerID,
		FirstName:  entry.FirstName,
		LastName:   entry.LastName,
		Email:      req.Email,
		Nonce:      req.PaymentMethodNonce,
	}

	paymenttkn, err := paymentC.CreateCustomer(paymentReq)
	if err != nil {
		resp.Error = errors.Wrap("failed to payment.CreateCustomer", err).SharedError()
		return resp
	}

	subC := sub.New(ctx)
	err = subC.UpdatePaymentToken(req.Email, paymenttkn)
	if err != nil {
		resp.Error = errors.Wrap("failed to sub.UpdatePaymentToken", err).SharedError()
		return resp
	}
	messageC := message.New(ctx)
	err = messageC.SendSMS("6155454989", fmt.Sprintf("Credit card updated. $$$ \nName: %s\nEmail: %s", entry.Name, entry.Email))
	if err != nil {
		utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
	}
	return resp
}

// SubmitCheckout submits a checkout.
func SubmitCheckout(ctx context.Context, r *http.Request) Response {
	req := new(pb.SubmitCheckoutReq)
	var err error
	// decode request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		return failedToDecode(err)
	}
	defer closeRequestBody(r)
	// end decode request
	resp := &pb.ErrorOnlyResp{}
	req.Email = strings.Replace(strings.ToLower(req.Email), " ", "", -1)
	req.PhoneNumber = strings.Replace(req.PhoneNumber, " ", "", -1)
	utils.Infof(ctx, "Request struct: %+v", req)
	err = validateSubmitCheckoutReq(req)
	if err != nil {
		resp.Error = errors.GetSharedError(err)
		return resp
	}

	key := datastore.NewKey(ctx, "ScheduleSignUp", req.Email, 0, nil)
	entry := &sub.SubscriptionSignUp{}
	err = datastore.Get(ctx, key, entry)
	if err != nil && err != datastore.ErrNoSuchEntity {
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to get ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		return resp
	}
	if entry.IsSubscribed {
		// user is already subscribed
		resp.Error = errInvalidParameter.WithMessage("You already have a subscription! :)").SharedError()
		return resp
	}
	inZone, address, err := InNashvilleZone(ctx, req.Address)
	if err != nil {
		resp.Error = errInternal.WithMessage("Woops! something went wrong").WithError(err).Annotate("failed inNashvilleZone").SharedError()
		return resp
	}
	// var planID string
	var servings int8
	var vegetarianServings int8
	var weeklyAmount float32
	switch req.Servings {
	case "":
		fallthrough
	case "0":
		servings = 0
	case "1":
		servings = 1
	case "2":
		servings = 2
	default:
		servings = 4
	}
	switch req.VegetarianServings {
	case "":
		fallthrough
	case "0":
		vegetarianServings = 0
	case "1":
		vegetarianServings = 1
	case "2":
		vegetarianServings = 2
	default:
		vegetarianServings = 4
	}
	weeklyAmount = sub.DerivePrice(vegetarianServings + servings)
	customerID := payment.GetIDFromEmail(req.Email)
	firstBoxDate := time.Now().Add(72 * time.Hour)
	for firstBoxDate.Weekday() != time.Monday {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24)
	}
	// TODO remove after Aug 21, 2017
	if firstBoxDate.Month() == time.August && firstBoxDate.Day() == 21 {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24 * 7)
	}
	paymentC := payment.New(ctx)
	paymentReq := &payment.CreateCustomerReq{
		CustomerID: customerID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Nonce:      req.PaymentMethodNonce,
	}

	paymenttkn, err := paymentC.CreateCustomer(paymentReq)
	if err != nil {
		resp.Error = errors.Wrap("failed to payment.CreateCustomer", err).SharedError()
		return resp
	}
	entry.Email = req.Email
	entry.Name = req.FirstName + " " + req.LastName
	entry.FirstName = strings.TrimSpace(req.FirstName)
	entry.LastName = strings.TrimSpace(req.LastName)
	entry.Address = *address
	if entry.Date.IsZero() {
		entry.Date = time.Now()
	}
	// entry.SubscriptionIDs = append(entry.SubscriptionIDs, subID)
	if inZone {
		entry.IsSubscribed = true
		entry.SubscriptionDate = time.Now()
		entry.WeeklyAmount = weeklyAmount
		entry.FirstBoxDate = firstBoxDate
		// entry.FirstPaymentDate = paymentDate
		entry.SubscriptionDay = time.Monday.String()
	}
	entry.CustomerID = customerID
	entry.DeliveryTips = req.DeliveryNotes
	entry.Servings = servings
	entry.VegetarianServings = vegetarianServings
	entry.PhoneNumber = req.PhoneNumber
	entry.PaymentMethodToken = paymenttkn
	entry.Reference = req.Reference
	_, err = datastore.Put(ctx, key, entry)
	if err != nil {
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to put ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		return resp
	}
	if !inZone {
		utils.Infof(ctx, "failed address zone zip(%s). Address: %s", address.Zip, address.String())
		// out of delivery range
		if address.Street == "" {
			resp.Error = errInvalidParameter.WithMessage("Please select an address from the list as you type your address!").SharedError()
			return resp
		}
		messageC := message.New(ctx)
		err = messageC.SendSMS("6153975516", fmt.Sprintf("Missed a customer. Out of zone. \nName: %s\nEmail: %s\nAddress: %s", entry.Name, entry.Email, entry.Address.StringNoAPT()))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		_ = messageC.SendSMS("9316445311", fmt.Sprintf("Missed a customer. Out of zone. \nName: %s\nEmail: %s\nAddress: %s", entry.Name, entry.Email, entry.Address.StringNoAPT()))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Atish. Err: %+v", err)
		}
		// TODO: add to some datastore to save address and stuff
		resp.Error = errInvalidParameter.WithMessage("Sorry, you are outside our delivery range! We'll let you know soon as we are in your area!").SharedError()
		return resp
	}
	if !appengine.IsDevAppServer() {
		messageC := message.New(ctx)
		err = messageC.SendSMS("6155454989", fmt.Sprintf("$$$ New subscriber checkout page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
		err = messageC.SendSMS("6153975516", fmt.Sprintf("$$$ New subscriber checkout page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		err = messageC.SendSMS("9316446755", fmt.Sprintf("$$$ New subscriber checkout page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Piyush. Err: %+v", err)
		}
		_ = messageC.SendSMS("9316445311", fmt.Sprintf("$$$ New subscriber checkout page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Atish. Err: %+v", err)
		}
	}
	subC := sub.New(ctx)
	err = subC.Free(firstBoxDate, req.Email)
	if err != nil {
		utils.Criticalf(ctx, "Failed to setup free sub box for new sign up(%s) for date(%v). Err:%v", req.Email, firstBoxDate, err)
	}
	if !strings.Contains(req.Email, "test.com") {
		mailC := mail.New(ctx)
		mailReq := &mail.UserFields{
			Email:             entry.Email,
			Name:              entry.Name,
			FirstName:         entry.FirstName,
			LastName:          entry.LastName,
			FirstDeliveryDate: firstBoxDate,
			AddTags:           []mail.Tag{mail.Subscribed, mail.Customer},
		}
		if vegetarianServings > 0 {
			mailReq.AddTags = append(mailReq.AddTags, mail.Vegetarian)
			mailReq.RemoveTags = append(mailReq.RemoveTags, mail.NonVegetarian)
		} else {
			mailReq.AddTags = append(mailReq.AddTags, mail.NonVegetarian)
			mailReq.RemoveTags = append(mailReq.RemoveTags, mail.Vegetarian)
		}
		err = mailC.UpdateUser(mailReq, getProjID())
		if err != nil {
			utils.Criticalf(ctx, "Failed to mail.UpdateUser email(%s). Err: %+v", entry.Email, err)
		}
	}
	return resp
}

// InNashvilleZone checks if an address is in Nashville zone.
func InNashvilleZone(ctx context.Context, addr *shared.Address) (bool, *types.Address, error) {
	var err error
	address := &types.Address{
		APT: addr.Apt,
	}
	if !(-90 <= addr.Latitude && addr.Latitude <= 90 && -180 <= addr.Longitude && addr.Longitude <= 180) || (addr.Latitude == 0 && addr.Longitude == 0) {
		addrStr := addr.FullAddress
		if addrStr == "" {
			addrStr = fmt.Sprintf(" %s, %s, %s %s, %s", addr.Street, addr.City, addr.State, addr.Zip, addr.Country)
		}
		address, err = maps.GetAddress(ctx, addrStr, addr.Apt)
		if err != nil {
			return false, nil, errors.Annotate(err, "failed to GetAddress")
		}
	} else {
		address.Street = addr.Street
		address.City = addr.City
		address.State = addr.State
		address.Zip = addr.Zip
		address.Country = addr.Country
		address.Latitude = addr.Latitude
		address.Longitude = addr.Longitude
	}
	fence := new(geofence.Geofence)
	key := datastore.NewKey(ctx, "Geofence", "", common.Nashville.ID(), nil)
	err = datastore.Get(ctx, key, fence)
	if err != nil {
		return false, nil, errInternal.WithError(err).Annotate("failed to db.Get")
	}
	polygon := geofence.NewPolygon(fence.Points)
	pnt := geofence.Point{
		GeoPoint: common.GeoPoint{
			Latitude:  address.Latitude,
			Longitude: address.Longitude,
		},
	}
	contains := polygon.Contains(pnt)
	return contains, address, nil
}

func handler(f func(context.Context, *http.Request) Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := appengine.NewContext(r)
		// if !setupDone {
		// 	err = setupWithContext(ctx)
		// 	if err != nil {
		// 		// TODO: Alert but send friendly error back
		// 		log.Fatal("failed to setup: %+v", err)
		// 		return
		// 	}
		// }
		// loggingC, err := logging.NewClient(ctx, r.URL.Path)
		// if err != nil {
		// 	errString := fmt.Sprintf("failed to get new logging client: %+v", err)
		// 	logging.Errorf(ctx, errString)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	_, _ = w.Write([]byte(errString))
		// }
		// ctx = context.WithValue(ctx, common.LoggingKey, loggingC)

		// call function
		resp := f(ctx, r)
		// Log errors
		sharedErr := resp.GetError()
		if sharedErr == nil || sharedErr.Code == shared.Code(0) {
			sharedErr = &shared.Error{
				Code: shared.Code_Success,
			}
		}
		if sharedErr != nil && sharedErr.Code != shared.Code_Success {
			// 	loggingC.LogRequestError(r, errors.GetErrorWithCode(sharedErr))
			logging.Errorf(ctx, "%+v", sharedErr)
		}
		// encode
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(resp)
		if err != nil {
			w.WriteHeader(int(resp.GetError().Code))
			_, _ = w.Write([]byte(fmt.Sprintf("failed to encode response: %+v", err)))
			return
		}
	}
}

// Response is a response to a rpc call. All responses contain an error.
type Response interface {
	GetError() *shared.Error
}

func closeRequestBody(r *http.Request) {
	_ = r.Body.Close()
}

func failedToDecode(err error) *pb.ErrorOnlyResp {
	return &pb.ErrorOnlyResp{
		Error: errBadRequest.WithError(err).Annotate("failed to decode").SharedError(),
	}
}
