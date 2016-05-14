package server

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/core/payment"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

func init() {
	r := httprouter.New()

	// chef stuff
	// loggedInChain := alice.New(middlewareLoggedIn)
	// r.Handler("GET", baseGigachefURL+"/*path", loggedInChain.ThenFunc(handleGigachefApp))

	r.GET(baseLoginURL, handleLogin)
	r.GET("/signin", handleSignin)
	r.GET(signOutURL, handleSignout)
	r.POST("/upload", handleUpload)

	r.POST("/sub-merchant-approved", handleSubMerchantApproved)
	r.POST("/sub-merchant-declined", handleSubMerchantDeclined)
	r.POST("/sub-merchant-disbursement-exception", handleDisbursementException)

	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	r.NotFound = http.HandlerFunc(handle404)
	http.Handle("/", r)
}

func handle404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("GIGA 404 page. :()"))
}

func handleSubMerchantApproved(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	paymentC := payment.New(ctx)
	btNotification(ctx, w, req, "SubMerchantApproved", paymentC.SubMerchantApproved)
}

func handleSubMerchantDeclined(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	paymentC := payment.New(ctx)
	btNotification(ctx, w, req, "SubMerchantDeclined", paymentC.SubMerchantDeclined)
}

func handleDisbursementException(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	paymentC := payment.New(ctx)
	btNotification(ctx, w, req, "DisbursementException", paymentC.DisbursementException)
}

func btNotification(ctx context.Context, w http.ResponseWriter, req *http.Request, fnName string, fn func(string, string) error) {
	err := req.ParseForm()
	if err != nil {
		utils.Errorf(ctx, "Error parsing %s request form: %v", fnName, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload := req.FormValue("bt_payload")
	signature := req.FormValue("bt_signature")
	err = fn(signature, payload)
	if err != nil {
		utils.Errorf(ctx, "Error parsing %s request form: %v", "SubMerchantApproved", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	page, err := ioutil.ReadFile("app/login.html")
	if err != nil {
		ctx := appengine.NewContext(req)
		utils.Errorf(ctx, "Error reading login page: %+v", err)
	}
	w.Write(page)
}

func handleSignin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := CurrentUser(w, req)
	if user != nil {
		http.Redirect(w, req, baseGigachefURL, http.StatusTemporaryRedirect)
	}
	http.Redirect(w, req, baseLoginURL, http.StatusTemporaryRedirect)
}

func handleSignout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	removeCookies(w)
	http.Redirect(w, req, homeURL, http.StatusTemporaryRedirect)
}
