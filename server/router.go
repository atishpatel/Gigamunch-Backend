package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"html/template"

	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/julienschmidt/httprouter"
)

var cookSignupPage []byte
var scheduleSignupPage []byte
var scheduleFormPage *template.Template

func init() {
	var err error
	cookSignupPage, err = ioutil.ReadFile("signedUp.html")
	if err != nil {
		log.Fatalf("Failed to read cookSignup page %#v", err)
	}
	scheduleSignupPage, err = ioutil.ReadFile("scheduleSignedUp.html")
	if err != nil {
		log.Fatalf("Failed to read scheduleSignup page %#v", err)
	}
	scheduleFormPage = template.New("scheduleForm")
	scheduleFormPage, err = scheduleFormPage.ParseFiles("scheduleForm.html")
	if err != nil {
		log.Fatalf("Failed to read scheduleForm page %#v", err)
	}
	r := httprouter.New()

	r.GET(baseLoginURL, handleLogin)
	r.GET(signOutURL, handleSignout)

	r.GET("/scheduleform/:email", handleScheduleForm)
	r.GET("/scheduleform", handleScheduleForm)

	http.HandleFunc("/startsubscription", handleScheduleSubscription)
	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	r.NotFound = http.HandlerFunc(handle404)
	http.HandleFunc("/get-upload-url", handleGetUploadURL)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/get-feed", handleGetFeed)
	http.HandleFunc("/get-item", handleGetItem)
	http.HandleFunc("/signedup", handleCookSignup)
	http.HandleFunc("/schedulesignedup", handleScheduleSignup)
	http.Handle("/", r)
}

func handle404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("GIGA 404 page. :()"))
}

func handleLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := CurrentUser(w, req)
	if user != nil {
		http.Redirect(w, req, baseCookURL, http.StatusTemporaryRedirect)
		return
	}
	removeCookies(w)
	http.Redirect(w, req, "/becomechef", http.StatusTemporaryRedirect)
}

func handleSignout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	removeCookies(w)
	http.Redirect(w, req, homeURL, http.StatusTemporaryRedirect)
}

type SignUp struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
	Name  string    `json:"name"`
}

func handleCookSignup(w http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	emailAddress := req.FormValue("email")
	name := req.FormValue("name")
	terp := req.FormValue("terp")
	utils.Infof(c, "email: %s, name: %s, terp: %s ", emailAddress, name, terp)
	if terp != "" {
		return
	}
	if emailAddress == "" {
		utils.Infof(c, "No email address. ")
		_, _ = w.Write(cookSignupPage)
		return
	}
	key := datastore.NewKey(c, "CookSignUp", emailAddress, 0, nil)
	entry := &SignUp{}
	err := datastore.Get(c, key, entry)
	if err == datastore.ErrNoSuchEntity {
		entry.Date = time.Now()
		entry.Email = emailAddress
		entry.Name = name
		_, err = datastore.Put(c, key, entry)
		if err != nil {
			utils.Criticalf(c, "Error putting CookSignupEmail in datastore ", err)
		}
		messageC := message.New(c)
		err = messageC.SendSMS("6153975516", fmt.Sprintf("%s just signed up using becomecook page. Get on that booty. \nEmail: %s", name, emailAddress))
		if err != nil {
			utils.Criticalf(c, "failed to send sms to Enis. Err: %+v", err)
		}
		_, _ = w.Write(cookSignupPage)
		return
	}
	utils.Errorf(c, "Error email already registered CookSignUp: emailaddress - %s, err - %#v", emailAddress, err)
	_, _ = w.Write(cookSignupPage)
}

func handleScheduleSignup(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write(scheduleSignupPage)
	// c := appengine.NewContext(req)
	// emailAddress := req.FormValue("email")
	// terp := req.FormValue("terp")
	// utils.Infof(c, "email: %s,  terp: %s ", emailAddress, terp)
	// if terp != "" {
	// 	return
	// }
	// if emailAddress == "" {
	// 	utils.Infof(c, "No email address. ")
	// 	// _, _ = w.Write(cookSignupPage)
	// 	// TODO redirect?
	// 	return
	// }
	// key := datastore.NewKey(c, "ScheduleSignUp", emailAddress, 0, nil)
	// entry := &SubscriptionSignUp{}
	// err := datastore.Get(c, key, entry)
	// if err == datastore.ErrNoSuchEntity {
	// 	entry.Date = time.Now()
	// 	entry.Email = emailAddress
	// 	_, err = datastore.Put(c, key, entry)
	// 	if err != nil {
	// 		utils.Criticalf(c, "Error putting ScheduleSignupEmail in datastore ", err)
	// 	}
	// 	if !appengine.IsDevAppServer() {
	// 		messageC := message.New(c)
	// 		err = messageC.SendSMS("6153975516", fmt.Sprintf("New sign up using schedule page. Get on that booty. \nEmail: %s", emailAddress))
	// 		if err != nil {
	// 			utils.Criticalf(c, "failed to send sms to Enis. Err: %+v", err)
	// 		}
	// 	}
	// } else {
	// 	utils.Errorf(c, "Error email already registered ScheduleSignUp: emailaddress - %s, err - %#v", emailAddress, err)
	// }
	// param := httprouter.Param{
	// 	Key:   "email",
	// 	Value: emailAddress,
	// }
	// handleScheduleForm(w, req, []httprouter.Param{param})
}

type SubscriptionSignUp struct {
	Email            string        `json:"email"`
	Date             time.Time     `json:"date"`
	Name             string        `json:"name"`
	Address          types.Address `json:"address"`
	CustomerID       string        `json:"customer_id"`
	SubscriptionIDs  []string      `json:"subscription_id"`
	IsSubscribed     bool          `json:"is_subscribed"`
	SubscriptionDate time.Time     `json:"subscription_date"`
	FirstPaymentDate time.Time     `json:"first_payment_date"`
	FirstBoxDate     time.Time     `json:"first_box_date"`
	DeliveryTime     int8          `json:"delivery_time"`
}

type scheduleSubscriptionResp struct {
	Err errors.ErrorWithCode `json:"err"`
}

type scheduleSubscriptionReq struct {
	Email              string        `json:"email"`
	Address            types.Address `json:"address"`
	Name               string        `json:"name"`
	PaymentMethodNonce string        `json:"payment-method-nonce"`
	Servings           string        `json:"servings"`
	Terp               string        `json:"terp"`
	DeliveryTime       int8          `json:"delivery_time"`
}

func (s *scheduleSubscriptionReq) valid() error {
	if s.Terp != "" {
		return errInvalidParameter.WithMessage("Trap was set off D:")
	}
	if s.Email == "" {
		return errInvalidParameter.WithMessage("Email address cannot be empty.").Wrap("No email address.")
	}
	if s.PaymentMethodNonce == "" {
		return errInvalidParameter.WithMessage("Invalid payment info.").Wrap("No payment nonce.")
	}
	if s.Name == "" {
		return errInvalidParameter.WithMessage("Name must be provided.").Wrap("No name.")
	}
	return nil
}

// called via ajax
func handleScheduleSubscription(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	resp := new(scheduleSubscriptionResp)
	defer handleResp(ctx, w, "ScheduleSubscription", resp.Err, resp)
	// decode request
	dec := json.NewDecoder(req.Body)
	sReq := new(scheduleSubscriptionReq)
	err := dec.Decode(&sReq)
	if err != nil {
		resp.Err = errInvalidParameter.WithError(err).Wrap("failed to decode request.")
		return
	}
	defer req.Body.Close()
	utils.Infof(ctx, "Request struct: %+v", sReq)
	err = sReq.valid()
	if err != nil {
		resp.Err = errors.Wrap("failed to validate request", err)
		return
	}
	if sReq.Address.GreatCircleDistance(types.GeoPoint{Latitude: 36.1513632, Longitude: -86.7255927}) > 35 {
		// out of delivery range
		// TODO add to some datastore to save address and stuff
		resp.Err = errInvalidParameter.WithMessage("Sorry, you are outside our delivery range! We'll let you know soon as we are in your area!")
		return
	}

	key := datastore.NewKey(ctx, "ScheduleSignUp", sReq.Email, 0, nil)
	entry := &SubscriptionSignUp{}
	err = datastore.Get(ctx, key, entry)
	if err != nil && err != datastore.ErrNoSuchEntity {
		resp.Err = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to get ScheduleSignUp email(%s) into datastore", sReq.Email)
		return
	}
	if entry.IsSubscribed {
		// user is already subscribed
		resp.Err = errInvalidParameter.WithMessage("You already have a subscription! :)")
		return
	}
	var planID string
	switch sReq.Servings {
	case "2":
		planID = "basic_2"
	case "4":
		planID = "basic_4"
	default:
		planID = "basic_1"
	}
	paymentC := payment.New(ctx)
	customerID := payment.GetIDFromEmail(sReq.Email)
	firstBoxDate := time.Now().Add(72 * time.Hour)
	for firstBoxDate.Weekday() != time.Monday {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24)
	}
	paymentDate := firstBoxDate.Add(time.Hour * 96)
	paymentSubReq := &payment.StartSubscriptionReq{
		CustomerID: customerID,
		Nonce:      sReq.PaymentMethodNonce,
		PlanID:     planID,
		StartDate:  paymentDate,
	}
	subID, err := paymentC.StartSubscription(paymentSubReq)
	if err != nil {
		resp.Err = errors.Wrap("failed to payment.StartSubscription", err)
		return
	}
	utils.Infof(ctx, "Subscription started for %s subID(%s)", sReq.Email, subID)
	entry.Email = sReq.Email
	entry.Name = sReq.Name
	entry.Address = sReq.Address
	if entry.Date.IsZero() {
		entry.Date = time.Now()
	}
	entry.SubscriptionIDs = append(entry.SubscriptionIDs, subID)
	entry.IsSubscribed = true
	entry.CustomerID = customerID
	entry.SubscriptionDate = time.Now()
	entry.FirstPaymentDate = paymentDate
	entry.FirstBoxDate = firstBoxDate
	entry.DeliveryTime = sReq.DeliveryTime
	_, err = datastore.Put(ctx, key, entry)
	if err != nil {
		resp.Err = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to put ScheduleSignUp email(%s) into datastore", sReq.Email)
		return
	}
	if !appengine.IsDevAppServer() {
		messageC := message.New(ctx)
		err = messageC.SendSMS("6153975516", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nEmail: %s", entry.Email))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
	}
	return
}

func handleScheduleForm(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	ctx := appengine.NewContext(req)
	email := param.ByName("email")
	terp := ""
	if email == "" {
		email = req.FormValue("email")
		terp = req.FormValue("terp")
	}
	if email == "" {
		// TODO redirect if email is empty

		return
	}

	utils.Infof(ctx, "email: %s,  terp: %s ", email, terp)
	if terp != "" {
		return
	}
	key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
	entry := &SubscriptionSignUp{}
	err := datastore.Get(ctx, key, entry)
	if err == datastore.ErrNoSuchEntity {
		entry.Date = time.Now()
		entry.Email = email
		_, err = datastore.Put(ctx, key, entry)
		if err != nil {
			utils.Criticalf(ctx, "Error putting ScheduleSignupEmail in datastore ", err)
		}
		if !appengine.IsDevAppServer() {
			messageC := message.New(ctx)
			err = messageC.SendSMS("6153975516", fmt.Sprintf("New sign up using schedule page. Get on that booty. \nEmail: %s", email))
			if err != nil {
				utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
			}
		}
	} else {
		utils.Warningf(ctx, "Warning: email already registered ScheduleSignUp: email - %s, err - %#v", email, err)
	}

	paymentC := payment.New(ctx)
	id := payment.GetIDFromEmail(email)
	utils.Infof(ctx, "id: %s", id)
	tkn, err := paymentC.GenerateToken(id)
	if err != nil {
		// Do something
		utils.Errorf(ctx, "Error payment.GenerateToken. Error: %+v", err)
		return
	}
	fields := scheduleFormFields{
		BTToken: tkn,
		Email:   email,
	}
	err = scheduleFormPage.ExecuteTemplate(w, "scheduleForm.html", fields)
	if err != nil {
		// do something?
		utils.Errorf(ctx, "Error scheduleFormPage.Execute. Error: %+v", err)
		return
	}
}

type scheduleFormFields struct {
	Email   string
	BTToken string
}
