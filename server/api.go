package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/server"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
	authold "github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/healthcheck"
	"github.com/atishpatel/Gigamunch-Backend/corenew/maps"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var (
	errBadRequest       = errors.BadRequestError
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
	errInternal         = errors.InternalServerError
)

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

func validateSubmitGiftCheckoutReq(r *giftCheckout) error {
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
	if r.ReferenceEmail == "" {
		return errInvalidParameter.WithMessage("Your email address cannot be empty.").Annotate("no email address")
	}
	return nil
}

// Login updates a user's payment.
func (s *server) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	req := new(pb.TokenOnlyReq)
	var err error
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.TokenOnlyResp{}

	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to get auth.NewClient")
	}

	usr, err := authC.Verify(req.Token)
	if err != nil {
		return errors.Annotate(err, "failed to auth.Verify")
	}
	// TODO: remove log
	log.Infof(ctx, "usr: %+v", usr)
	// TODO: find / create user and update token

	_, authToken, err := authold.GetSessionWithGToken(ctx, req.Token)
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
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
	entry := &sub.SubscriptionSignUp{}
	err = datastore.Get(ctx, key, entry)
	if err == datastore.ErrNoSuchEntity {
		resp.Error = errBadRequest.WithMessage(fmt.Sprintf("Cannot find user with email: %s", email)).Wrapf("failed to get ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		utils.Criticalf(ctx, "failed to update payment because can't find email(%s) tkn(%s): %+v", email, req.PaymentMethodNonce, err)
		return resp
	}
	if err != nil {
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to get ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		utils.Criticalf(ctx, "failed to update payment because can't find email(%s) tkn(%s): %+v", email, req.PaymentMethodNonce, err)
		return resp
	}
	paymentC := payment.New(ctx)
	paymentReq := &payment.CreateCustomerReq{
		CustomerID: entry.CustomerID,
		FirstName:  entry.FirstName,
		LastName:   entry.LastName,
		Email:      email,
		Nonce:      req.PaymentMethodNonce,
	}

	paymenttkn, err := paymentC.CreateCustomer(paymentReq)
	if err != nil {
		utils.Criticalf(ctx, "failed to update payment: failed to sub.CreateCustomer: email(%s) %+v", email, err)
		resp.Error = errors.Wrap("failed to payment.CreateCustomer", err).SharedError()
		return resp
	}
	subC := sub.New(ctx)
	err = subC.UpdatePaymentToken(email, paymenttkn)
	if err != nil {
		utils.Criticalf(ctx, "failed to update payment: failed to sub.UpdatePaymentToken: email(%s) tkn(%s) %+v", email, paymenttkn, err)
		resp.Error = errors.Wrap("failed to sub.UpdatePaymentToken", err).SharedError()
		return resp
	}
	messageC := message.New(ctx)
	err = messageC.SendAdminSMS("6155454989", fmt.Sprintf("Credit card updated. $$$ \nName: %s\nEmail: %s", entry.Name, entry.Email))
	if err != nil {
		utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
	}
	unpaidSublogs, err := subC.GetSubscriberUnpaidSublogs(email)
	if err != nil {
		utils.Errorf(ctx, "failed to GetSubscriberUnpaidSublogs: %+v", err)
		return resp
	}
	tasksC := tasks.New(ctx)
	t := time.Now()
	for _, sublog := range unpaidSublogs {
		if sublog.Date.After(t) {
			continue
		}
		req := &tasks.ProcessSubscriptionParams{
			SubEmail: req.Email,
			Date:     sublog.Date,
		}
		err = tasksC.AddProcessSubscription(t, req)
		if err != nil {
			utils.Errorf(ctx, "failed to AddProcessSubscription: %+v", err)
		}
		t = t.Add(time.Minute * 5)
	}
	return resp
}

// SubmitCheckout submits a checkout.
func SubmitCheckout(ctx context.Context, r *http.Request) Response {
	req := new(pb.SubmitCheckoutReq)
	var err error
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}
	req.Email = strings.Replace(strings.ToLower(req.Email), " ", "", -1)
	req.PhoneNumber = strings.Replace(req.PhoneNumber, " ", "", -1)
	logging.Infof(ctx, "Request struct: %+v", req)
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
	firstBoxDate := time.Now().Add(81 * time.Hour)
	for firstBoxDate.Weekday() != time.Monday {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24)
	}
	if req.FirstDeliveryDate != "" {
		firstBoxDate, err = time.Parse(time.RFC3339, req.FirstDeliveryDate)
		if err != nil || firstBoxDate.Weekday() == time.Tuesday {
			firstBoxDate = firstBoxDate.Add(-12 * time.Hour)
		}
		if err != nil || firstBoxDate.Weekday() == time.Sunday {
			firstBoxDate = firstBoxDate.Add(12 * time.Hour)
		}
		if err != nil || firstBoxDate.Weekday() != time.Monday {
			resp.Error = errBadRequest.WithMessage("Invalid first delivery day selected.").SharedError()
			utils.Criticalf(ctx, "user selected invalid start date: %+v", req.FirstDeliveryDate)
			return resp
		}
	}
	// TODO: remove
	logging.Infof(ctx, "firstBoxDate would change from %s to %s to %s", firstBoxDate, firstBoxDate.UTC(), firstBoxDate.UTC().Truncate(24*time.Hour))
	tmpNow := time.Now()
	logging.Infof(ctx, "now would change from %s to %s to %s", tmpNow, tmpNow.UTC(), tmpNow.UTC().Truncate(24*time.Hour))

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
	entry.UpdatePhoneNumber(req.PhoneNumber)
	entry.PaymentMethodToken = paymenttkn
	entry.Reference = req.Reference
	entry.ReferenceEmail = req.ReferenceEmail
	for _, c := range req.Campaigns {
		found := false
		var timeStamp time.Time
		timeStamp, _ = time.Parse(time.RFC3339, c.Timestamp)
		for _, loggedC := range entry.Campaigns {
			if loggedC.Campaign != c.Campaign {
				continue
			}
			diff := timeStamp.Sub(loggedC.Timestamp)
			if diff < 0 {
				diff *= -1
			}
			if diff < time.Hour {
				found = true
			}
		}
		if !found {
			entry.Campaigns = append(entry.Campaigns, campaingFromPB(c))
		}
	}
	_, err = datastore.Put(ctx, key, entry)
	if err != nil {
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to put ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		return resp
	}
	if !inZone {
		logging.Infof(ctx, "failed address zone zip(%s). Address: %s", address.Zip, address.String())
		// out of delivery range
		if address.Street == "" {
			resp.Error = errInvalidParameter.WithMessage("Please select an address from the list as you type your address!").SharedError()
			return resp
		}
		messageC := message.New(ctx)
		err = messageC.SendAdminSMS("6153975516", fmt.Sprintf("Missed a customer. Out of zone. \nName: %s\nEmail: %s\nAddress: %s", entry.Name, entry.Email, entry.Address.StringNoAPT()))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		_ = messageC.SendAdminSMS("9316445311", fmt.Sprintf("Missed a customer. Out of zone. \nName: %s\nEmail: %s\nAddress: %s", entry.Name, entry.Email, entry.Address.StringNoAPT()))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Atish. Err: %+v", err)
		}
		// TODO: add to some datastore to save address and stuff
		resp.Error = errInvalidParameter.WithMessage("Sorry, you are outside our delivery range! We'll let you know soon as we are in your area!").SharedError()
		return resp
	}
	if !appengine.IsDevAppServer() && !strings.Contains(entry.Email, "@test.com") {
		messageC := message.New(ctx)
		err = messageC.SendAdminSMS("6155454989", fmt.Sprintf("$$$ New subscriber checkout page. \nName: %s\nEmail: %s\nReference: %s\nReference Email: %s", entry.Name, entry.Email, entry.Reference, entry.ReferenceEmail))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
		err = messageC.SendAdminSMS("6153975516", fmt.Sprintf("$$$ New subscriber checkout page. \nName: %s\nEmail: %s\nReference: %s\nReference Email: %s", entry.Name, entry.Email, entry.Reference, entry.ReferenceEmail))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		err = messageC.SendAdminSMS("9316446755", fmt.Sprintf("$$$ New subscriber checkout page. \nName: %s\nEmail: %s\nReference: %s\nReference Email: %s", entry.Name, entry.Email, entry.Reference, entry.ReferenceEmail))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Piyush. Err: %+v", err)
		}
		_ = messageC.SendAdminSMS("9316445311", fmt.Sprintf("$$$ New subscriber checkout page. \nName: %s\nEmail: %s\nReference: %s\nReference Email: %s", entry.Name, entry.Email, entry.Reference, entry.ReferenceEmail))
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
		log, serverInfo, _, err := setupLoggingAndServerInfo(ctx, "/api/SubmitCheckout")
		if err != nil {
			return errors.Wrap("failed to setupLoggingAndServerInfo", err)
		}
		mailC, err := mail.NewClient(ctx, log, serverInfo)
		mailReq := &mail.UserFields{
			Email:             entry.Email,
			FirstName:         entry.FirstName,
			LastName:          entry.LastName,
			FirstDeliveryDate: firstBoxDate,
			VegServings:       entry.VegetarianServings,
			NonVegServings:    entry.Servings,
		}
		durationTillFirstMeal := time.Until(firstBoxDate.UTC().Truncate(24 * time.Hour))
		if durationTillFirstMeal > 0 && durationTillFirstMeal < ((6*24)-12)*time.Hour {
			mailReq.AddTags = append(mailReq.AddTags, mail.GetPreviewEmailTag(firstBoxDate))
		}
		err = mailC.SubActivated(mailReq)
		if err != nil {
			utils.Criticalf(ctx, "Failed to mail.UpdateUser email(%s). Err: %+v", entry.Email, err)
		}
		// add to task queue
		taskC := tasks.New(ctx)
		r := &tasks.ProcessSubscriptionParams{
			SubEmail: entry.Email,
			Date:     firstBoxDate,
		}
		err = taskC.AddProcessSubscription(firstBoxDate.Add(-24*time.Hour), r)
		if err != nil {
			return errors.Wrap("failed to tasks.AddProcessSubscription", err)
		}
	}
	return resp
}

func campaingFromPB(c *pb.Campaign) sub.Campaign {
	t, _ := time.Parse(time.RFC3339, c.Timestamp)
	return sub.Campaign{
		Timestamp: t,
		Source:    c.Source,
		Campaign:  c.Campaign,
		Term:      c.Term,
		Content:   c.Content,
		Medium:    c.Medium,
	}
}

type giftCheckout struct {
	pb.SubmitCheckoutReq
	NumGiftDinners int    `json:"num_gift_dinners"`
	ReferenceEmail string `json:"reference_email"`
	GiftRevealDate string `json:"gift_reveal_date"`
}

// SubmitGiftCheckout submits a checkout.
func SubmitGiftCheckout(ctx context.Context, r *http.Request) Response {
	req := new(giftCheckout)
	var err error
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}
	req.Email = strings.Replace(strings.ToLower(req.Email), " ", "", -1)
	req.PhoneNumber = strings.Replace(req.PhoneNumber, " ", "", -1)
	logging.Infof(ctx, "Request struct: %+v", req)
	err = validateSubmitGiftCheckoutReq(req)
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
	customerID := payment.GetIDFromEmail(req.ReferenceEmail)
	firstBoxDate := time.Now().Add(81 * time.Hour)
	for firstBoxDate.Weekday() != time.Monday {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24)
	}
	if req.FirstDeliveryDate != "" {
		firstBoxDate, err = time.Parse(time.RFC3339, req.FirstDeliveryDate)
		if err != nil || firstBoxDate.Weekday() == time.Tuesday {
			firstBoxDate = firstBoxDate.Add(-12 * time.Hour)
		}
		if err != nil || firstBoxDate.Weekday() == time.Sunday {
			firstBoxDate = firstBoxDate.Add(12 * time.Hour)
		}
		if err != nil || firstBoxDate.Weekday() != time.Monday {
			resp.Error = errBadRequest.WithMessage("Invalid first delivery day selected.").SharedError()
			utils.Criticalf(ctx, "user selected invalid first delivery date: %+v", req.FirstDeliveryDate)
			return resp
		}
	}
	giftRevealDate := firstBoxDate
	if req.GiftRevealDate != "" {
		tmpDate, err := time.Parse(time.RFC3339, req.GiftRevealDate)
		if err != nil {
			resp.Error = errBadRequest.WithMessage("Invalid gift reveal date.").SharedError()
			utils.Criticalf(ctx, "user selected invalid gift reaveal date: %+v", req.FirstDeliveryDate)
			return resp
		}
		giftRevealDate = tmpDate
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
	entry.UpdatePhoneNumber(req.PhoneNumber)
	entry.PaymentMethodToken = paymenttkn
	entry.Reference = req.Reference
	entry.ReferenceEmail = req.ReferenceEmail
	entry.GiftRevealDate = giftRevealDate
	if req.NumGiftDinners > 2 {
		req.NumGiftDinners++
	}
	entry.NumGiftDinners = req.NumGiftDinners
	_, err = datastore.Put(ctx, key, entry)
	if err != nil {
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to put ScheduleSignUp email(%s) into datastore", req.Email).SharedError()
		return resp
	}
	if !inZone {
		logging.Infof(ctx, "failed address zone zip(%s). Address: %s", address.Zip, address.String())
		// out of delivery range
		if address.Street == "" {
			resp.Error = errInvalidParameter.WithMessage("Please select an address from the list as you type your address!").SharedError()
			return resp
		}
		messageC := message.New(ctx)
		err = messageC.SendAdminSMS("6153975516", fmt.Sprintf("Missed a customer. Out of zone. \nName: %s\nEmail: %s\nAddress: %s", entry.Name, entry.Email, entry.Address.StringNoAPT()))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		_ = messageC.SendAdminSMS("9316445311", fmt.Sprintf("Missed a customer. Out of zone. \nName: %s\nEmail: %s\nAddress: %s", entry.Name, entry.Email, entry.Address.StringNoAPT()))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Atish. Err: %+v", err)
		}
		// TODO: add to some datastore to save address and stuff
		resp.Error = errInvalidParameter.WithMessage("Sorry, you are outside our delivery range! We'll let you know soon as we are in your area!").SharedError()
		return resp
	}
	if !appengine.IsDevAppServer() && !strings.Contains(entry.Email, "@test.com") {
		messageC := message.New(ctx)
		err = messageC.SendAdminSMS("6155454989", fmt.Sprintf("$$$ Gift checkout page. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
		err = messageC.SendAdminSMS("6153975516", fmt.Sprintf("$$$ Gift checkout page. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		err = messageC.SendAdminSMS("9316446755", fmt.Sprintf("$$$ Gift checkout page. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Piyush. Err: %+v", err)
		}
		_ = messageC.SendAdminSMS("9316445311", fmt.Sprintf("$$$ Gift checkout page. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Atish. Err: %+v", err)
		}
	}
	subC := sub.New(ctx)
	if entry.NumGiftDinners > 2 {
		err = subC.Free(firstBoxDate, entry.Email)
		if err != nil {
			utils.Criticalf(ctx, "error in gifCheckout: Failed to setup free sub box for new sign up(%s) for date(%v). Err:%v", entry.Email, firstBoxDate, err)
		}
	} else {
		err = subC.Setup(firstBoxDate, entry.Email, entry.Servings, entry.VegetarianServings, entry.WeeklyAmount, 6, entry.PaymentMethodToken, entry.CustomerID)
		if err != nil {
			utils.Criticalf(ctx, "error in giftCheckout: Failed to setup sub box for new sign up(%s) for date(%v). Err:%v", entry.Email, firstBoxDate, err)
		}
	}
	if !strings.Contains(entry.Email, "test.com") {
		tasksC := tasks.New(ctx)
		tasksReq := &tasks.UpdateDripParams{
			Email: entry.Email,
		}
		err = tasksC.AddUpdateDrip(entry.GiftRevealDate, tasksReq)
		if err != nil {
			utils.Criticalf(ctx, "error in giftCheckout: failed to AddUpdateDrip: %+v", err)
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
	fence, err := getNasvilleGeopoint(ctx)
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

func getNasvilleGeopoint(ctx context.Context) (*geofence.Geofence, error) {
	fence := new(geofence.Geofence)
	key := datastore.NewKey(ctx, "Geofence", "", common.Nashville.ID(), nil)
	err := datastore.Get(ctx, key, fence)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}
	if err == datastore.ErrNoSuchEntity {
		fence = &geofence.Geofence{
			ID:   "Nashville",
			Type: geofence.ServiceZone,
			Name: "Nashville",
			Points: []geofence.Point{
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3195513, Longitude: -86.5475464}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3347628, Longitude: -86.5248873}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3521861, Longitude: -86.5420532}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3438904, Longitude: -86.7253876}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2719574, Longitude: -86.7576599}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2459345, Longitude: -86.8139649}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2304274, Longitude: -86.8805695}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1971752, Longitude: -86.9011731}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1417563, Longitude: -86.8956756}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.079072, Longitude: -87.0570373}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.9518857, Longitude: -87.0206452}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.8253964, Longitude: -87.0253076}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.8233808, Longitude: -86.8421173}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.8462043, Longitude: -86.6670227}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 35.9526416, Longitude: -86.6629813}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.0080895, Longitude: -86.6210983}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.0835115, Longitude: -86.5626526}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1323292, Longitude: -86.5956116}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.1927545, Longitude: -86.5647125}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.2614385, Longitude: -86.5805054}},
				geofence.Point{GeoPoint: common.GeoPoint{Latitude: 36.3195513, Longitude: -86.5475464}},
			},
		}
	}
	return fence, nil
}

// DeviceCheckin updates a user's payment.
func DeviceCheckin(ctx context.Context, r *http.Request) Response {
	req := new(healthcheck.Device)
	var err error
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}

	healthC := healthcheck.New(ctx)
	err = healthC.Checkin(req)
	if err != nil {
		resp.Error = errors.GetSharedError(err)
		return resp
	}
	return resp
}

func setupLoggingAndServerInfo(ctx context.Context, path string) (*logging.Client, *common.ServerInfo, common.DB, error) {
	dbC, err := db.NewClient(ctx, projID, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get database client: %+v", err)
	}
	// Setup logging
	serverInfo := &common.ServerInfo{
		ProjectID:           projID,
		IsStandardAppEngine: true,
	}
	log, err := logging.NewClient(ctx, "admin", path, dbC, serverInfo)
	if err != nil {
		return nil, nil, nil, err
	}
	return log, serverInfo, dbC, nil
}
