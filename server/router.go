package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"html/template"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"
	"github.com/atishpatel/Gigamunch-Backend/corenew/mail"
	"github.com/atishpatel/Gigamunch-Backend/corenew/maps"
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
var projID string

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

	r.GET(signOutURL, handleSignout)

	// r.GET("/scheduleform/:email", handleScheduleForm)
	// r.GET("/scheduleform", handleScheduleForm)

	r.POST("/process-subscription", handelProcessSubscription)
	http.HandleFunc("/process-subscribers", handelProcessSubscribers)
	http.HandleFunc(tasks.SendEmailURL, handleSendEmail)

	http.HandleFunc("/startsubscription", handleScheduleSubscription)
	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	http.HandleFunc("/get-upload-url", handleGetUploadURL)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/get-feed", handleGetFeed)
	http.HandleFunc("/get-item", handleGetItem)
	http.HandleFunc("/signedup", handleCookSignup)
	http.HandleFunc("/schedulesignedup", handleScheduleSignup)
	addTemplateRoutes(r)
	addAPIRoutes(r)
	http.Handle("/", r)
}

func getProjID() string {
	projID = os.Getenv("PROJECTID")
	return projID
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

// inNashvilleZone checks if an address is in Nashville zone.
func inNashvilleZone(ctx context.Context, addr *types.Address) (bool, error) {
	var err error
	if !addr.GeoPoint.Valid() {
		// TODO get geopoint form address
		err = maps.GetGeopointFromAddress(ctx, addr)
		if err != nil {
			return false, errors.Annotate(err, "failed to GetGeopointFromAddress")
		}
	}
	fence := new(geofence.Geofence)
	key := datastore.NewKey(ctx, "Geofence", common.Nashville.String(), 0, nil)
	err = datastore.Get(ctx, key, fence)
	if err != nil {
		return false, errInternal.WithError(err).Annotate("failed to db.Get")
	}
	polygon := geofence.NewPolygon(fence.Points)
	pnt := geofence.Point{
		GeoPoint: common.GeoPoint{
			Latitude:  addr.Latitude,
			Longitude: addr.Longitude,
		},
	}
	contains := polygon.Contains(pnt)
	return contains, nil
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
	inZone, err := inNashvilleZone(ctx, &sReq.Address)
	if err != nil {
		resp.Err = errInternal.WithMessage("Woops! something went wrong").WithError(err).Annotate("failed inNashvilleZone")
		return
	}
	if !inZone {
		utils.Infof(ctx, "failed address zone zip(%s). Address: %s", sReq.Address.Zip, sReq.Address.String())
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
	case "2":
		servings = 2
	default:
		servings = 4
	}
	switch sReq.VegetarianServings {
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
	customerID := payment.GetIDFromEmail(sReq.Email)
	firstBoxDate := time.Now().Add(72 * time.Hour)
	for firstBoxDate.Weekday() != time.Monday {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24)
	}
	// TODO remove after Aug 21, 2017
	if firstBoxDate.Month() == time.August && firstBoxDate.Day() == 21 {
		firstBoxDate = firstBoxDate.Add(time.Hour * 24 * 7)
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
		err = messageC.SendSMS("6155454989", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Chris. Err: %+v", err)
		}
		err = messageC.SendSMS("6153975516", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Enis. Err: %+v", err)
		}
		err = messageC.SendSMS("9316446755", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Piyush. Err: %+v", err)
		}
		_ = messageC.SendSMS("9316445311", fmt.Sprintf("$$$ New subscriber schedule page. Email that booty. \nName: %s\nEmail: %s\nReference: %s", entry.Name, entry.Email, entry.Reference))
		if err != nil {
			utils.Criticalf(ctx, "failed to send sms to Atish. Err: %+v", err)
		}
	}
	subC := sub.New(ctx)
	err = subC.Free(firstBoxDate, sReq.Email)
	if err != nil {
		utils.Criticalf(ctx, "Failed to setup free sub box for new sign up(%s) for date(%v). Err:%v", sReq.Email, firstBoxDate, err)
	}
	// mailC := mail.New(ctx)
	// err = mailC.SendWelcomeEmail(entry)
	// if err != nil {
	// 	utils.Criticalf(ctx, "Failed to send welcome email for new sign up(%s). Err:%v", sReq.Email, err)
	// }
	mailC := mail.New(ctx)
	mailReq := &mail.UserFields{
		Email:             entry.Email,
		Name:              entry.Name,
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

func handleScheduleForm(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	ctx := appengine.NewContext(req)
	email := param.ByName("email")
	terp := ""
	if email == "" {
		email = req.FormValue("email")
		terp = req.FormValue("terp")
	}
	var tkn string
	var err error
	if email != "" {
		utils.Infof(ctx, "email: %s,  terp: %s ", email, terp)
		if terp != "" {
			return
		}
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
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
			if !entry.IsSubscribed {
				mailC := mail.New(ctx)
				err = mailC.AddTag(entry.Email, mail.LeftWebsiteEmail)
				if err != nil {
					utils.Criticalf(ctx, "Error added mail.AddTag for email(%s). Err: %+v", entry.Email, err)
				}
			}
			// hourFromNow := time.Now().Add(time.Hour)
			// sendEmailReq := &tasks.SendEmailParams{
			// 	Email: entry.Email,
			// 	Type:  "welcome",
			// }
			// tasksC := tasks.New(ctx)
			// err = tasksC.AddSendEmail(hourFromNow, sendEmailReq)
			// if err != nil {
			// 	utils.Criticalf(ctx, "Error added sendemail to queue for email(%s). Err: %+v", entry.Email, err)
			// }
		} else {
			utils.Warningf(ctx, "Warning: email already registered ScheduleSignUp: email - %s, err - %#v", email, err)
		}

		paymentC := payment.New(ctx)
		id := payment.GetIDFromEmail(email)
		utils.Infof(ctx, "id: %s", id)
		tkn, err = paymentC.GenerateToken(id)
		if err != nil {
			// Do something
			utils.Errorf(ctx, "Error payment.GenerateToken. Error: %+v", err)
			return
		}
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

func handleSendEmail(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	parms, err := tasks.ParseSendEmailRequest(req)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseSendEmailRequest. Err:%+v", err)
		return
	}
	subC := sub.New(ctx)
	s, err := subC.GetSubscriber(parms.Email)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.GetSubscriber. Err:%+v", err)
		return
	}
	if !s.IsSubscribed {
		mailC := mail.New(ctx)
		err = mailC.SendIntroEmail(s)
		if err != nil {
			utils.Criticalf(ctx, "failed to mail.SendIntroEmail. Err:%+v", err)
			return
		}
	}
}
