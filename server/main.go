package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/auth"

	pbcommon "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/server"
	"github.com/atishpatel/Gigamunch-Backend/core/activity"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	subnew "github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

var cookSignupPage []byte

var projID string

func init() {
	// TODO: Remove
	var err error
	cookSignupPage, err = ioutil.ReadFile("signedUp.html")
	if err != nil {
		log.Fatalf("Failed to read cookSignup page %#v", err)
	}
	projID = os.Getenv("PROJECT_ID")

	// Setup Server
	s := new(server)
	err = s.setup()
	if err != nil {
		log.Fatal("failed to setup", err)
	}

	r := httprouter.New()
	r.GET("/signout", handleSignout)
	r.POST("/process-subscription", handelProcessSubscription)
	http.HandleFunc("/process-subscribers", handelProcessSubscribers)
	http.HandleFunc("/signedup", handleCookSignup)

	// route templates
	addTemplateRoutes(r)
	// route api
	http.HandleFunc("/api/v1/Login", s.handler(s.Login))
	http.HandleFunc("/api/v1/SubmitCheckout", handler(SubmitCheckout))
	http.HandleFunc("/api/v1/SubmitGiftCheckout", handler(SubmitGiftCheckout))
	http.HandleFunc("/api/v1/UpdatePayment", handler(UpdatePayment))
	http.HandleFunc("/api/v1/DeviceCheckin", handler(DeviceCheckin))

	http.Handle("/", r)

}

func handleSignout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:   "",
		MaxAge: -1,
		Secure: true,
	})
	http.SetCookie(w, &http.Cookie{Name: "AUTHTKN", MaxAge: -1, Secure: true})

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
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

// TODO: Remove
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

// TODO: Remove
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

func (s *server) getUserFromRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) (*common.User, error) {
	token := r.Header.Get("auth-token")
	if token == "" {
		return nil, errBadRequest.Annotate("auth-token is empty")
	}
	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return nil, errors.Annotate(err, "failed to get auth.NewClient")
	}
	user, err := authC.Verify(token)
	if err != nil {
		return nil, errors.Annotate(err, "failed to get auth.GetUser")
	}

	return user, nil
}

// server
type server struct {
	once       sync.Once
	serverInfo *common.ServerInfo
	db         *db.Client
	sqlDB      *sqlx.DB
	log        *logging.Client
}

func (s *server) setup() error {
	var err error
	projID := os.Getenv("PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "PROJECT_ID"`)
	}
	// Setup sql db
	sqlConnectionString := os.Getenv("MYSQL_CONNECTION")
	if sqlConnectionString == "" {
		log.Fatal(`You need to set the environment variable "MYSQL_CONNECTION"`)
	}
	if appengine.IsDevAppServer() {
		sqlConnectionString = "root@/gigamunch"
	}
	s.sqlDB, err = sqlx.Connect("mysql", sqlConnectionString)
	if err != nil {
		if !appengine.IsDevAppServer() { // TODO: remove
			return fmt.Errorf("failed to get sql database client: %+v", err)
		}
	}
	s.serverInfo = &common.ServerInfo{
		ProjectID:           projID,
		IsStandardAppEngine: true,
	}
	return nil
}

func (s *server) handler(f handle) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Access-Control-Allow-Headers, Access-Control-Allow-Origin, auth-token")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.Contains(r.URL.Hostname(), "gigamunchapp.com") {
			url := "https://eatgigamunch.com" + r.URL.Path
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			return
		}
		// get context
		ctx := appengine.NewContext(r)
		ctx = context.WithValue(ctx, common.ContextUserID, "")
		ctx = context.WithValue(ctx, common.ContextUserEmail, "")
		s.db, err = db.NewClient(ctx, s.serverInfo.ProjectID, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO:
			w.Write([]byte(fmt.Sprintf("failed to get database client: %+v", err)))
			return
		}
		// create logging client
		log, err := logging.NewClient(ctx, "admin", r.URL.Path, s.db, s.serverInfo)
		if err != nil {
			errString := fmt.Sprintf("failed to get new logging client: %+v", err)
			logging.Errorf(ctx, errString)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errString))
		}
		// call function
		resp := f(ctx, w, r, log)
		if resp == nil {
			return
		}
		// Log errors
		sharedErr := resp.GetError()
		if sharedErr == nil {
			sharedErr = &pbcommon.Error{
				Code: pbcommon.Code_Success,
			}
		}
		if sharedErr != nil && sharedErr.Code != pbcommon.Code_Success && sharedErr.Code != pbcommon.Code(0) {
			logging.Errorf(ctx, "request error: %+v", errors.GetErrorWithCode(sharedErr))
			// log.RequestError((r, errors.GetErrorWithCode(sharedErr), )
			w.WriteHeader(int(sharedErr.Code))
			// Wrap error in ErrorOnlyResp
			if _, ok := resp.(errors.ErrorWithCode); ok {
				resp = &pb.ErrorOnlyResp{
					Error: sharedErr,
				}
			}
		}
		// encode
		encoder := json.NewEncoder(w)
		err = encoder.Encode(resp)
		if err != nil {
			w.WriteHeader(int(resp.GetError().Code))
			_, _ = w.Write([]byte(fmt.Sprintf("failed to encode response: %+v", err)))
			return
		}
	}
}

// Request helpers
func decodeRequest(ctx context.Context, r *http.Request, v interface{}) error {
	if r.Method == "GET" {
		decoder := schema.NewDecoder()
		err := decoder.Decode(v, r.URL.Query())
		logging.Infof(ctx, "Query: %+v", r.URL.Query())
		if err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		logging.Infof(ctx, "Body: %s", body)
		err = json.Unmarshal(body, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func failedToDecode(err error) *pb.ErrorOnlyResp {
	return &pb.ErrorOnlyResp{
		Error: errBadRequest.WithError(err).Annotate("failed to decode").SharedError(),
	}
}

// Response is a response to a rpc call. All responses contain an error.
type Response interface {
	GetError() *pbcommon.Error
}

type handle func(context.Context, http.ResponseWriter, *http.Request, *logging.Client) Response
