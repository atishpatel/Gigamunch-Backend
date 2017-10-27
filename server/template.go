package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/corenew/mail"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/sub"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var (
	templateFiles = []string{
		"templates/head.html",
		"templates/theme.html",
		"templates/footer.html",
		"templates/checkout.html",
		"templates/checkout-thank-you.html",
	}
	templates = template.Must(template.New("all").Delims("[[", "]]").ParseFiles(templateFiles...))
)

// Page is the basic info required for a template page.
type Page struct {
	Title string
}

func addTemplateRoutes() {
	http.HandleFunc("/checkout", handleCheckout)
	http.HandleFunc("/checkout-thank-you", handleCheckoutThankYou)
}

// display the named template
func display(ctx context.Context, w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		utils.Errorf(ctx, "failed to ExecuteTemplate: %+v", err)
	}
}

type checkoutPage struct {
	Page
	Email string
}

func handleCheckout(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	page := &checkoutPage{
		Page: Page{
			Title: "Checkout | Gigamunch",
		},
	}
	defer display(ctx, w, "checkout", page)
	email := req.FormValue("email")
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

type checkoutThankYouPage struct {
	Page
	FirstName         string
	FirstDeliveryDate string
}

func handleCheckoutThankYou(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	page := &checkoutThankYouPage{
		Page: Page{
			Title: "Thank you | Gigamunch",
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
