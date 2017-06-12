package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"html/template"

	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
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

	r.POST("/process-subscription", handelProcessSubscription)
	http.HandleFunc("/process-subscribers", handelProcessSubscribers)

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
	emailAddress := strings.Replace(strings.ToLower(req.FormValue("email")), " ", "", -1)
	name := req.FormValue("name")
	terp := req.FormValue("terp")
	utils.Infof(c, "email or phone: %s, name: %s, terp: %s ", emailAddress, name, terp)
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
		err = messageC.SendSMS("6153975516", fmt.Sprintf("Cook %s just signed up using becomecook page. Get on that booty. \nEmail: %s", name, emailAddress))
		if err != nil {
			utils.Criticalf(c, "failed to send sms to Enis. Err: %+v", err)
		}
		_ = messageC.SendSMS("6155454989", fmt.Sprintf("Cook %s just signed up using becomecook page. Get on that booty. \nEmail: %s", name, emailAddress))
		_, _ = w.Write(cookSignupPage)
		return
	}
	utils.Errorf(c, "Error email already registered CookSignUp: emailaddress - %s, err - %#v", emailAddress, err)
	_, _ = w.Write(cookSignupPage)
}

func handleScheduleSignup(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write(scheduleSignupPage)
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
	VegetarianServings string        `json:"vegetarian_servings"`
	Terp               string        `json:"terp"`
	DeliveryTime       int8          `json:"delivery_time"`
	Reference          string        `json:"reference"`
	PhoneNumber        string        `json:"phone_number"`
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
	sReq.Email = strings.Replace(strings.ToLower(sReq.Email), " ", "", -1)
	sReq.PhoneNumber = strings.Replace(sReq.PhoneNumber, " ", "", -1)
	utils.Infof(ctx, "Request struct: %+v", sReq)
	err = sReq.valid()
	if err != nil {
		resp.Err = errors.Wrap("failed to validate request", err)
		return
	}
	if sReq.Address.GreatCircleDistance(types.GeoPoint{Latitude: 36.045565, Longitude: -86.784328}) > 25 {
		// out of delivery range
		if sReq.Address.Street == "" {
			resp.Err = errInvalidParameter.WithMessage("Please select an address from the list as you type your address!")
			return
		}
		// TODO add to some datastore to save address and stuff
		resp.Err = errInvalidParameter.WithMessage("Sorry, you are outside our delivery range! We'll let you know soon as we are in your area!")
		return
	}

	key := datastore.NewKey(ctx, "ScheduleSignUp", sReq.Email, 0, nil)
	entry := &sub.SubscriptionSignUp{}
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
	// var planID string
	var servings int8
	var vegetarianServings int8
	var weeklyAmount float32
	switch sReq.Servings {
	case "":
		fallthrough
	case "0":
		servings = 0
	case "1":
		servings = 1
		weeklyAmount += 17
	case "2":
		servings = 2
		weeklyAmount += float32(servings*15) + 2.93
	default:
		servings = 4
		weeklyAmount += float32(servings*14) + 5.46
	}
	switch sReq.VegetarianServings {
	case "":
		fallthrough
	case "0":
		vegetarianServings = 0
	case "1":
		vegetarianServings = 1
		weeklyAmount += 17
	case "2":
		vegetarianServings = 2
		weeklyAmount += float32(vegetarianServings * 15)
	default:
		vegetarianServings = 4
		weeklyAmount += float32(vegetarianServings * 14)
	}
	customerID := payment.GetIDFromEmail(sReq.Email)
	firstBoxDate := time.Now().Add(48 * time.Hour)
	for firstBoxDate.Weekday() != time.Monday {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24)
	}
	paymentC := payment.New(ctx)
	paymentTokenReq := &payment.GetDefaultPaymentTokenReq{
		CustomerID: customerID,
	}
	paymenttkn, err := paymentC.GetDefaultPaymentToken(paymentTokenReq)
	if err != nil {
		resp.Err = errors.Wrap("failed to payment.GetDefaultPaymentToken", err)
		return
	}
	entry.Email = sReq.Email
	entry.Name = sReq.Name
	entry.Address = sReq.Address
	if entry.Date.IsZero() {
		entry.Date = time.Now()
	}
	// entry.SubscriptionIDs = append(entry.SubscriptionIDs, subID)
	entry.IsSubscribed = true
	entry.CustomerID = customerID
	entry.SubscriptionDate = time.Now()
	// entry.FirstPaymentDate = paymentDate
	entry.FirstBoxDate = firstBoxDate
	entry.DeliveryTime = sReq.DeliveryTime
	entry.Servings = servings
	entry.VegetarianServings = vegetarianServings
	entry.PhoneNumber = sReq.PhoneNumber
	entry.SubscriptionDay = time.Monday.String()
	entry.PaymentMethodToken = paymenttkn
	entry.WeeklyAmount = weeklyAmount
	entry.Reference = sReq.Reference
	_, err = datastore.Put(ctx, key, entry)
	if err != nil {
		resp.Err = errInternal.WithMessage("Woops! Something went wrong. Try again in a few minutes.").WithError(err).Wrapf("failed to put ScheduleSignUp email(%s) into datastore", sReq.Email)
		return
	}
	if !appengine.IsDevAppServer() {
		messageC := message.New(ctx)
		err = messageC.SendSMS("6153975516", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		_ = messageC.SendSMS("9316446755", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		_ = messageC.SendSMS("6155454989", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		_ = messageC.SendSMS("8607485603", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		_ = messageC.SendSMS("9316445311", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
	}
	subC := sub.New(ctx)
	err = subC.Free(firstBoxDate, sReq.Email)
	if err != nil {
		utils.Criticalf(ctx, "Failed to setup free sub box for new sign up(%s) for date(%v). Err:%v", sReq.Email, firstBoxDate, err)
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
	entry := &sub.SubscriptionSignUp{}
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
			_ = messageC.SendSMS("9316446755", fmt.Sprintf("New sign up using schedule page. Get on that booty. \nEmail: %s", email))
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

func handelProcessSubscription(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	parms, err := tasks.ParseProcessSubscriptionRequest(req)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseProcessSubscriptionRequest. Err:%+v", err)
		return
	}
	subC := sub.New(ctx)
	err = subC.Process(parms.Date, parms.SubEmail)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.Process(Date:%s SubEmail:%s). Err:%+v", parms.Date, parms.SubEmail, err)
		// TODO schedule for later?
		return
	}
}

func handelProcessSubscribers(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	in2days := time.Now().Add(48 * time.Hour)
	subC := sub.New(ctx)
	err := subC.SetupSubLogs(in2days)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupSubLogs(Date:%v). Err:%+v", in2days, err)
		return
	}
}
