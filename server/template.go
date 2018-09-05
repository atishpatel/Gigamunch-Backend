package server

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/core/lead"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"

	subnew "github.com/atishpatel/Gigamunch-Backend/core/sub"
	subold "github.com/atishpatel/Gigamunch-Backend/corenew/sub"
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
	r.GET("/", redirictOldHost(handleHome))
	r.GET("/schedule", handleHome)
	r.GET("/passport", handlePassport)
	r.GET("/login", redirictOldHost(handleLogin))
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

	r.GET("/campaign", handleCampaignHome)
	r.GET("/campaign-checkout", handleCampaignCheckout)
	r.GET("/campaign-thank-you", handleCheckoutThankYou)
	r.NotFound = new(handler404)
}

func redirictOldHost(f httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if strings.Contains(req.URL.Hostname(), "gigamunchapp.com") {
			url := "https://eatgigamunch.com" + req.URL.Path
			http.Redirect(w, req, url, http.StatusPermanentRedirect)
			return
		}
		f(w, req, params)
	}
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
	Email           string
	FirstName       string
	LastName        string
	PhoneNumber     string
	Address         string
	APT             string
	DeliveryNotes   string
	Reference       string
	ThankYouPageURL string
}

func handleCheckout(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	page := &checkoutPage{}
	displayCheckout(w, req, params, page)
}

func handleCampaignCheckout(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	page := &checkoutPage{
		ThankYouPageURL: "/campaign-thank-you",
	}

	displayCheckout(w, req, params, page)
}

func displayCheckout(w http.ResponseWriter, req *http.Request, params httprouter.Params, page *checkoutPage) {
	ctx := appengine.NewContext(req)
	defer display(ctx, w, "checkout", page)
	page.ID = "checkout"
	if page.ThankYouPageURL == "" {
		page.ThankYouPageURL = "/checkout-thank-you"
	}
	email := req.FormValue("email")
	terp := req.FormValue("terp")
	if email == "" {
		email = params.ByName("email")
	}
	email = strings.TrimSpace(strings.ToLower(email))
	// TODO: add referred email address

	if email != "" && terp == "" && strings.Contains(email, "@") {
		logging.Infof(ctx, "email: %s", email)

		// save lead
		log, serverInfo, db, _ := setupLoggingAndServerInfo(ctx, "/checkout")
		leadC, err := lead.NewClient(ctx, log, db, serverInfo)
		if err != nil {
			log.Errorf(ctx, "failed to lead.NewClient: %+v", err)
		}
		err = leadC.Create(email)
		if err != nil {
			log.Errorf(ctx, "failed to lead.Create: %+v", err)
		}
		page.Email = email
		// load page if tried before
		subC := subold.New(ctx)
		s, err := subC.GetSubscriber(email)
		if err == nil {
			page.FirstName = s.FirstName
			page.LastName = s.LastName
			page.PhoneNumber = s.PhoneNumber
			if s.Address.GeoPoint.Valid() {
				page.Address = s.Address.StringNoAPT()
			}
			page.APT = s.Address.APT
			page.DeliveryNotes = s.DeliveryTips
			page.Reference = s.Reference
		}
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

	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		suboldC := subold.New(ctx)
		entry, err := suboldC.GetSubscriber(email)
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

	if strings.Contains(email, "@") {
		page.Email = email
		logging.Infof(ctx, "email: %s", email)

		suboldC := subold.New(ctx)
		s, err := suboldC.GetSubscriber(email)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			return
		}
		page.FirstName = s.FirstName
		if s.FirstName == "" {
			page.FirstName = getFirstName(s.Name)
		}
		// increase page count
		log, serverInfo, db, _ := setupLoggingAndServerInfo(ctx, "/referral")
		subnewC, err := subnew.NewClient(ctx, log, db, nil, serverInfo)
		if err == nil {
			err = subnewC.IncrementPageCount(email, 1, 0)
			if err != nil {
				logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			}
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

	if strings.Contains(email, "@") {
		logging.Infof(ctx, "email: %s", email)
		page.ReferrerName = email

		suboldC := subold.New(ctx)
		entry, err := suboldC.GetSubscriber(email)
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

		// increase page count
		log, serverInfo, db, _ := setupLoggingAndServerInfo(ctx, "/referred")
		subnewC, err := subnew.NewClient(ctx, log, db, nil, serverInfo)
		if err == nil {
			err = subnewC.IncrementPageCount(email, 0, 1)
			if err != nil {
				logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			}
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

	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		page.Email = email

		suboldC := subold.New(ctx)
		entry, err := suboldC.GetSubscriber(email)
		if err != nil {
			logging.Errorf(ctx, "failed to datastore.Get: %+v", err)
			return
		}
		page.FirstName = entry.FirstName
		if entry.FirstName == "" {
			page.FirstName = getFirstName(entry.Name)
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

	if email != "" {
		logging.Infof(ctx, "email: %s", email)
		page.ReferrerName = email

		suboldC := subold.New(ctx)
		entry, err := suboldC.GetSubscriber(email)
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

	}
}

func getFirstName(name string) string {
	nameReplaced := strings.NewReplacer(".", "", "Mr", "", "Ms", "", "Mrs", "", "The", "").Replace(strings.Title(name))
	nameArray := strings.Split(nameReplaced, " ")
	return nameArray[0]
}

type homePage struct {
	Page
	ReferrerName    string
	CheckoutPageURL string
}

func handleHome(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	page := &homePage{}
	displayHome(w, req, params, page)
}

func handleCampaignHome(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	page := &homePage{
		Page: Page{
			CampaignName: "CampaignPage",
		},
		CheckoutPageURL: "/campaign-checkout",
	}
	displayHome(w, req, params, page)
}

func handlePassport(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	page := &homePage{
		Page: Page{
			ID:           "passport",
			CampaignName: "Passport",
		},
	}
	displayHome(w, req, params, page)
}

func displayHome(w http.ResponseWriter, req *http.Request, params httprouter.Params, page *homePage) {
	ctx := appengine.NewContext(req)
	if page.ID == "" {
		page.ID = "home"
	}
	if page.CheckoutPageURL == "" {
		page.CheckoutPageURL = "/checkout"
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
