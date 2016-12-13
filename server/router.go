package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/julienschmidt/httprouter"
)

func init() {
	var err error
	cookSignupPage, err = ioutil.ReadFile("signedUp.html")
	if err != nil {
		log.Fatalf("Failed to read cookSignup page %#v", err)
	}
	r := httprouter.New()

	r.GET(baseLoginURL, handleLogin)
	r.GET(signOutURL, handleSignout)

	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	r.NotFound = http.HandlerFunc(handle404)
	http.HandleFunc("/get-upload-url", handleGetUploadURL)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/get-feed", handleGetFeed)
	http.HandleFunc("/get-item", handleGetItem)
	http.HandleFunc("/signedup", handleCookSignup)
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

var cookSignupPage []byte

type CookSignup struct {
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
	entry := &CookSignup{}
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
