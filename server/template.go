package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/corenew/mail"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	templates = template.Must(template.New("all").Delims("[[", "]]").ParseGlob("templates/*"))
)

// Page is the basic info required for a template page.
type Page struct {
	Title string
}

func addTemplateRoutes(r *httprouter.Router) {
	r.GET("/checkout", handleCheckout)
	r.GET("/checkout-thank-you", handleCheckoutThankYou)
	r.GET("/update-payment", handleUpdatePayment)
	r.GET("/thank-you", handleThankYou)
	r.GET("/gift", handleGift)
	r.GET("/gift/:email", handleGift)
	r.GET("/referred", handleReferred)
	r.GET("/referred/:email", handleReferred)
	r.NotFound = new(handler404)
}

// display the named template
func display(ctx context.Context, w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl := templates.Lookup(tmplName)
	if tmpl == nil {
		utils.Errorf(ctx, "failed to Lookup: %s", tmplName)
		w.WriteHeader(http.StatusNotFound)
		_ = templates.ExecuteTemplate(w, "404", nil)
		return
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		utils.Errorf(ctx, "failed to ExecuteTemplate: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = templates.ExecuteTemplate(w, "500", nil)
	}
}

type checkoutPage struct {
	Page
	Email string
}

func handleCheckout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &checkoutPage{
		Page: Page{
			Title: "Checkout",
		},
	}
	defer display(ctx, w, "checkout", page)
	email := req.FormValue("email")
	// TODO: add referred email address
	var err error
	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		page.Email = email
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
			mailC := mail.New(ctx)
			err = mailC.AddTag(entry.Email, mail.LeftWebsiteEmail)
			if err != nil {
				utils.Criticalf(ctx, "Error added mail.AddTag for email(%s). Err: %+v", entry.Email, err)
			}
		} else if err != nil {
			utils.Criticalf(ctx, "failed to add email(%s) to ScheduleSignUp: err - %+v", email, err)
		} else {
			utils.Infof(ctx, "email already registered ScheduleSignUp: email - %s, err - %#v", email, err)
		}
	}
}

func handleUpdatePayment(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &checkoutPage{
		Page: Page{
			Title: "Update Payment",
		},
	}
	defer display(ctx, w, "update-payment", page)
	email := req.FormValue("email")
	if email != "" {
		page.Email = email
	}
}

type checkoutThankYouPage struct {
	Page
	FirstName         string
	FirstDeliveryDate string
}

func handleCheckoutThankYou(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &checkoutThankYouPage{
		Page: Page{
			Title: "Thank you",
		},
	}
	defer display(ctx, w, "checkout-thank-you", page)
	email := req.FormValue("email")
	var err error
	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			return
		}
		page.FirstName = entry.FirstName
		page.FirstDeliveryDate = mail.DateString(entry.FirstBoxDate)
	}
}

func handleThankYou(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &checkoutThankYouPage{
		Page: Page{
			Title: "Thank you",
		},
	}
	defer display(ctx, w, "thank-you", page)
}

type giftPage struct {
	Page
	Email     string
	FirstName string
}

func handleGift(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &giftPage{
		Page: Page{
			Title: "Gift A Piece of the World",
		},
	}
	defer display(ctx, w, "gift", page)
	email := req.FormValue("email")
	if email == "" {
		email = params.ByName("email")
	}

	var err error
	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			return
		}
		page.Email = email
		if entry.FirstName != "" {
			page.FirstName = entry.FirstName
		} else {
			page.FirstName, _ = splitName(entry.Name)
		}

	}
}

type referredPage struct {
	Page
	ReferrerName string
}

func handleReferred(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &referredPage{
		Page: Page{
			Title: "A Piece of the World",
		},
	}
	defer display(ctx, w, "referred", page)
	email := req.FormValue("email")
	if email == "" {
		email = params.ByName("email")
	}
	var err error
	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			return
		}
		page.ReferrerName = entry.FirstName + " " + entry.LastName
		if page.ReferrerName == "" {
			page.ReferredName = entry.Name
		}
	}
}

type handler404 struct {
}

func (h *handler404) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	w.WriteHeader(http.StatusNotFound)
	defer display(ctx, w, "404", nil)
}

func handle500(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	w.WriteHeader(http.StatusInternalServerError)
	defer display(ctx, w, "500", nil)
}

func splitName(name string) (string, string) {
	first := ""
	last := ""
	name = strings.Title(strings.TrimSpace(name))
	lastSpace := strings.LastIndex(name, " ")
	if lastSpace == -1 {
		first = name
	} else {
		first = name[:lastSpace]
		last = name[lastSpace:]
	}
	return first, last
}
