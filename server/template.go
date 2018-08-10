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
	ID             string
	ReferenceEmail string
	CampaignName   string
}

func addTemplateRoutes(r *httprouter.Router) {
	r.GET("/", handleHome)
	r.GET("/schedule", handleHome)
	r.GET("/passport", handlePassport)
	r.GET("/login", handleLogin)
	r.GET("/terms", handleTerms)
	r.GET("/privacy", handlePrivacy)
	r.GET("/checkout", handleCheckout)
	r.GET("/gift-checkout", handleGiftCheckout)
	r.GET("/scheduleform/:email", handleCheckout)
	r.GET("/scheduleform", handleCheckout)
	r.GET("/checkout-thank-you", handleCheckoutThankYou)
	r.GET("/update-payment", handleUpdatePayment)
	r.GET("/thank-you", handleThankYou)
	r.GET("/referral", handleReferral)
	r.GET("/referral/:email", handleReferral)
	r.GET("/referred", handleReferred)
	r.GET("/referred/:email", handleReferred)
	r.GET("/gift", handleGift)
	r.GET("/gift/:email", handleGift)
	r.GET("/gifted", handleGifted)
	r.GET("/gifted/:email", handleGifted)
	r.GET("/becomechef", handleBecomecook)
	r.GET("/becomecook", handleBecomecook)
	r.NotFound = new(handler404)
}

// display the named template
func display(ctx context.Context, w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl := templates.Lookup(tmplName)
	if tmpl == nil {
		utils.Errorf(ctx, "failed to Lookup: %s", tmplName)
		w.WriteHeader(http.StatusNotFound)
		_ = templates.ExecuteTemplate(w, "404", &Page{ID: "404"})
		return
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		utils.Errorf(ctx, "failed to ExecuteTemplate: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = templates.ExecuteTemplate(w, "500", &Page{ID: "500"})
	}
}

type checkoutPage struct {
	Page
	Email         string
	FirstName     string
	LastName      string
	PhoneNumber   string
	Address       string
	APT           string
	DeliveryNotes string
	Reference     string
}

func handleCheckout(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &checkoutPage{
		Page: Page{
			ID: "checkout",
		},
	}
	defer display(ctx, w, "checkout", page)
	email := req.FormValue("email")
	terp := req.FormValue("terp")
	if email == "" {
		email = params.ByName("email")
	}
	email = strings.TrimSpace(strings.ToLower(email))
	// TODO: add referred email address
	var err error
	if email != "" && terp == "" && strings.Contains(email, "@") {
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
			if !appengine.IsDevAppServer() && !strings.Contains(entry.Email, "@test.com") {
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
		page.FirstName = entry.FirstName
		page.LastName = entry.LastName
		page.PhoneNumber = entry.PhoneNumber
		if entry.Address.GeoPoint.Valid() {
			page.Address = entry.Address.StringNoAPT()
		}
		page.APT = entry.Address.APT
		page.DeliveryNotes = entry.DeliveryTips
		page.Reference = entry.Reference
	}
}

func handleGiftCheckout(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &Page{
		ID: "gift-checkout",
	}
	defer display(ctx, w, "gift-checkout", page)
}

func handleUpdatePayment(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &checkoutPage{
		Page: Page{
			ID: "update-payment",
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
			ID: "checkout-thank-you",
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
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
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
			ID: "thank-you",
		},
	}
	defer display(ctx, w, "thank-you", page)
}

type referralPage struct {
	Page
	Email     string
	FirstName string
}

func handleReferral(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &referralPage{
		Page: Page{
			ID: "referral",
		},
	}
	defer display(ctx, w, "referral", page)
	email := req.FormValue("email")
	if email == "" {
		email = params.ByName("email")
	}
	var err error
	if strings.Contains(email, "@") {
		page.Email = email
		logging.Infof(ctx, "email: %s", email)
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			return
		}
		page.FirstName = entry.FirstName
		if entry.FirstName == "" {
			page.FirstName = getFirstName(entry.Name)
		}
		entry.ReferralPageOpens++
		_, err = datastore.Put(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
		}
	}
}

func handleReferred(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &homePage{
		Page: Page{
			ID:           "referred",
			CampaignName: "Referred",
		},
	}
	defer display(ctx, w, "home", page)
	email := req.FormValue("email")
	if email == "" {
		email = params.ByName("email")
	}
	var err error
	if strings.Contains(email, "@") {
		logging.Infof(ctx, "email: %s", email)
		page.ReferrerName = email
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			return
		}
		if entry.Name != "" {
			page.ReferrerName = entry.Name
		}
		if entry.FirstName != "" {
			page.ReferrerName = entry.FirstName + " " + entry.LastName
		}
		page.ReferenceEmail = entry.Email
		entry.ReferredPageOpens++
		_, err = datastore.Put(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
		}
	}
}

func handleGift(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &referralPage{
		Page: Page{
			ID: "gift",
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
		page.Email = email
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			return
		}
		page.FirstName = entry.FirstName
		if entry.FirstName == "" {
			page.FirstName = getFirstName(entry.Name)
		}
		entry.GiftPageOpens++
		_, err = datastore.Put(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
		}
	}
}

func handleGifted(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &homePage{
		Page: Page{
			ID:           "gifted",
			CampaignName: "Gifted",
		},
	}
	defer display(ctx, w, "home", page)
	email := req.FormValue("email")
	if email == "" {
		email = params.ByName("email")
	}
	var err error
	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		page.ReferrerName = email
		key := datastore.NewKey(ctx, "ScheduleSignUp", email, 0, nil)
		entry := &sub.SubscriptionSignUp{}
		err = datastore.Get(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			return
		}
		if entry.Name != "" {
			page.ReferrerName = entry.Name
		}
		if entry.FirstName != "" {
			page.ReferrerName = entry.FirstName + " " + entry.LastName
		}
		page.ReferenceEmail = entry.Email
		entry.GiftedPageOpens++
		_, err = datastore.Put(ctx, key, entry)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
		}
	}
}

func getFirstName(name string) string {
	nameReplaced := strings.NewReplacer(".", "", "Mr", "", "Ms", "", "Mrs", "", "The", "").Replace(strings.Title(name))
	nameArray := strings.Split(nameReplaced, " ")
	return nameArray[0]
}

type homePage struct {
	Page
	ReferrerName string
}

func handleHome(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &homePage{
		Page: Page{
			ID: "home",
		},
	}
	defer display(ctx, w, "home", page)
}

func handlePassport(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &homePage{
		Page: Page{
			ID:           "passport",
			CampaignName: "Passport",
		},
	}
	defer display(ctx, w, "home", page)
}

func handleLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &Page{
		ID: "login",
	}
	defer display(ctx, w, "login", page)
}

func handleBecomecook(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &Page{
		ID: "becomecook",
	}
	defer display(ctx, w, "becomecook", page)
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

func handleTerms(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &Page{
		ID: "terms",
	}
	defer display(ctx, w, "terms", page)
}

func handlePrivacy(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	page := &Page{
		ID: "privacy",
	}
	defer display(ctx, w, "privacy", page)
}
