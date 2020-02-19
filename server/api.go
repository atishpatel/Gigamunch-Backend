package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/discount"

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/jmoiron/sqlx"

	pbcommon "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"
	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbserver"
	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/healthcheck"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

var (
	errBadRequest       = errors.BadRequestError
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
	errInternal         = errors.InternalServerError
)

// Login updates a user's payment.
func (s *server) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.TokenOnlyReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	_, err = subC.VerifyAndUpdateAuth(req.Token)
	if err != nil {
		return errors.Annotate(err, "failed to sub.VerifyAndUpdateAuth")
	}
	return resp
}

// UpdatePayment updates a user's payment.
func UpdatePayment(ctx context.Context, r *http.Request) Response {
	var err error
	req := new(pb.UpdatePaymentReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	subC := subold.New(ctx)
	entry, err := subC.GetSubscriber(email)
	if err != nil {
		ewc := errors.GetErrorWithCode(err)
		if ewc.Code == errors.CodeNotFound {
			resp.Error = errBadRequest.WithMessage(fmt.Sprintf("Cannot find user with email: %s", email)).Wrapf("failed to get email(%s) into datastore", req.Email).SharedError()
			utils.Criticalf(ctx, "failed to update payment because can't find email(%s) tkn(%s): %+v", email, req.PaymentMethodNonce, err)
			return resp
		}
		resp.Error = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to get email(%s) into datastore", req.Email).SharedError()
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
		utils.Criticalf(ctx, "failed to update payment: failed to subold.CreateCustomer: email(%s) %+v", email, err)
		resp.Error = errors.Wrap("failed to payment.CreateCustomer", err).SharedError()
		return resp
	}

	err = subC.UpdatePaymentToken(email, paymenttkn)
	if err != nil {
		utils.Criticalf(ctx, "failed to update payment: failed to subold.UpdatePaymentToken: email(%s) tkn(%s) %+v", email, paymenttkn, err)
		resp.Error = errors.Wrap("failed to subold.UpdatePaymentToken", err).SharedError()
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

func campaingFromPB(c *pbcommon.Campaign) common.Campaign {
	t, _ := time.Parse(time.RFC3339, c.Timestamp)
	return common.Campaign{
		Timestamp: t,
		Source:    c.Source,
		Campaign:  c.Campaign,
		Term:      c.Term,
		Content:   c.Content,
		Medium:    c.Medium,
	}
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
		if sharedErr == nil || sharedErr.Code == pbcommon.Code(0) {
			sharedErr = &pbcommon.Error{
				Code: pbcommon.Code_Success,
			}
		}
		if sharedErr != nil && sharedErr.Code != pbcommon.Code_Success {
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

// DeviceCheckin updates a user's payment.
func DeviceCheckin(ctx context.Context, r *http.Request) Response {
	var err error
	req := new(healthcheck.Device)
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

func setupAll(ctx context.Context, path string) (*logging.Client, *common.ServerInfo, common.DB, *sqlx.DB, error) {
	dbC, err := db.NewClient(ctx, projID, nil)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get database client: %+v", err)
	}
	// Setup logging
	serverInfo := &common.ServerInfo{
		ProjectID:           projID,
		IsStandardAppEngine: true,
	}
	log, err := logging.NewClient(ctx, "server", path, dbC, serverInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	sqlConnectionString := os.Getenv("MYSQL_CONNECTION")
	if sqlConnectionString == "" && !appengine.IsDevAppServer() {
		return nil, nil, nil, nil, fmt.Errorf(`You need to set the environment variable "MYSQL_CONNECTION"`)
	}
	sqlDB, err := sqlx.Connect("mysql", sqlConnectionString+"?collation=utf8mb4_general_ci&parseTime=true")
	if err != nil && !appengine.IsDevAppServer() {
		return nil, nil, nil, nil, fmt.Errorf("failed to get sql database client: %+v", err)
	}
	return log, serverInfo, dbC, sqlDB, nil
}

// SubmitCheckoutv2 submits a checkout.
func (s *server) SubmitCheckoutv2(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.SubmitCheckoutReq)
	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	resp := &pb.ErrorOnlyResp{}
	logging.Infof(ctx, "Request struct: %+v", req)
	req.Email = strings.ToLower(req.Email)
	var servingsNonVegetarian int8
	var servingsVegetarian int8
	switch req.Servings {
	case "":
		fallthrough
	case "0":
		servingsNonVegetarian = 0
	case "1":
		servingsNonVegetarian = 1
	case "2":
		servingsNonVegetarian = 2
	case "4":
		servingsNonVegetarian = 4
	default:
		servingsNonVegetarian = 4
	}
	switch req.VegetarianServings {
	case "":
		fallthrough
	case "0":
		servingsVegetarian = 0
	case "1":
		servingsVegetarian = 1
	case "2":
		servingsVegetarian = 2
	case "4":
		servingsVegetarian = 4
	default:
		servingsVegetarian = 4
	}
	firstBoxDate := time.Now().Add(81 * time.Hour)
	// for firstBoxDate.Weekday() != time.Monday && firstBoxDate.Weekday() != time.Thursday {
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
		// if err != nil || firstBoxDate.Weekday() == time.Friday {
		// 	firstBoxDate = firstBoxDate.Add(-12 * time.Hour)
		// }
		// if err != nil || firstBoxDate.Weekday() == time.Wednesday {
		// 	firstBoxDate = firstBoxDate.Add(12 * time.Hour)
		// }

		if err != nil || (firstBoxDate.Weekday() != time.Monday && firstBoxDate.Weekday() != time.Thursday) {
			resp.Error = errBadRequest.WithMessage("Invalid first delivery day selected.").SharedError()
			utils.Criticalf(ctx, "user selected invalid start date: %+v", req.FirstDeliveryDate)
			return resp
		}
	}
	var campaigns []common.Campaign
	for _, c := range req.Campaigns {
		campaigns = append(campaigns, campaingFromPB(c))
	}

	address, err := serverhelper.AddressFromPB(ctx, req.Address)
	if err != nil {
		return errors.Annotate(err, "failed to decode address")
	}
	promoBreakdown := discount.GetPromoBreakdown(req.Promo)

	subC, err := sub.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to sub.NewClient")
	}
	createReq := &sub.CreateReq{
		Email:                 req.Email,
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		PhoneNumber:           req.PhoneNumber,
		Address:               *address,
		DeliveryNotes:         req.DeliveryNotes,
		Reference:             req.Reference,
		ReferenceEmail:        req.ReferenceEmail,
		PaymentMethodNonce:    req.PaymentMethodNonce,
		ServingsNonVegetarian: servingsNonVegetarian,
		ServingsVegetarian:    servingsVegetarian,
		FirstDeliveryDate:     firstBoxDate,
		Campaigns:             campaigns,
		DiscountPercent:       promoBreakdown.DiscountPercent,
		DiscountAmount:        promoBreakdown.DiscountAmount,
	}
	_, err = subC.Create(createReq)
	if err != nil {
		return errors.Annotate(err, "failed to sub.Create")
	}
	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to auth.NewClient")
	}
	err = authC.CreateUserIfNoExist(req.Email, req.Password, req.FirstName+" "+req.LastName)
	if err != nil {
		return errors.Annotate(err, "failed to auth.CreateUserIfNoExist")
	}
	return resp
}
