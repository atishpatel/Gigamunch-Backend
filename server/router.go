package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	subnew "github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/julienschmidt/httprouter"
)

var cookSignupPage []byte

var projID string

func init() {
	var err error
	cookSignupPage, err = ioutil.ReadFile("signedUp.html")
	if err != nil {
		log.Fatalf("Failed to read cookSignup page %#v", err)
	}

	projID = os.Getenv("PROJECTID")

	r := httprouter.New()

	r.GET(signOutURL, handleSignout)

	// r.GET("/scheduleform/:email", handleScheduleForm)
	// r.GET("/scheduleform", handleScheduleForm)

	r.POST("/process-subscription", handelProcessSubscription)
	http.HandleFunc("/process-subscribers", handelProcessSubscribers)

	// http.HandleFunc("/startsubscription", handleScheduleSubscription)
	// // admin stuff
	// adminChain := alice.New(middlewareAdmin)
	// r.Handler("GET", adminHomeURL, adminChain.ThenFunc(handleAdminHome))
	// http.HandleFunc("/get-upload-url", handleGetUploadURL)
	// http.HandleFunc("/upload", handleUpload)
	// http.HandleFunc("/get-feed", handleGetFeed)
	// http.HandleFunc("/get-item", handleGetItem)
	http.HandleFunc("/signedup", handleCookSignup)
	// http.HandleFunc("/schedulesignedup", handleScheduleSignup)
	addTemplateRoutes(r)
	addAPIRoutes(r)
	http.Handle("/", r)

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

func handelProcessSubscription(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	parms, err := tasks.ParseProcessSubscriptionRequest(req)
	if err != nil {
		utils.Criticalf(ctx, "failed to tasks.ParseProcessSubscriptionRequest. Err:%+v", err)
		return
	}
	log, serverInfo, db, _ := setupLoggingAndServerInfo(ctx, "/process-subscription")
	activityC, _ := activity.NewClient(ctx, log, db, nil, serverInfo)
	err = activityC.Process(parms.Date, parms.SubEmail)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.Process(Date:%s SubEmail:%s). Err:%+v", parms.Date, parms.SubEmail, err)
		// TODO schedule for later?
		return
	}
}

func handelProcessSubscribers(w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	in2days := time.Now().Add(48 * time.Hour)
	log, serverInfo, db, _ := setupLoggingAndServerInfo(ctx, "/process-subscribers")
	subC, _ := subnew.NewClient(ctx, log, db, nil, serverInfo)
	err := subC.SetupActivities(in2days)
	if err != nil {
		utils.Criticalf(ctx, "failed to sub.SetupSubLogs(Date:%v). Err:%+v", in2days, err)
		return
	}
}
