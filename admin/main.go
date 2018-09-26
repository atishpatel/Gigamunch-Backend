package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"

	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
	"google.golang.org/appengine"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	once       sync.Once
	serverInfo *common.ServerInfo
	db         *db.Client
	sqlDB      *sqlx.DB
}

var (
	errPermissionDenied = errors.PermissionDeniedError
	errUnauthenticated  = errors.UnauthenticatedError
	errBadRequest       = errors.BadRequestError
	errInternalError    = errors.InternalServerError
)

func init() {
	s := new(server)
	err := s.setup()
	if err != nil {
		log.Fatal("failed to setup", err)
	}
	// Auth
	http.HandleFunc("/admin/api/v1/SetAdmin", s.handler(s.userAdmin(s.SetAdmin)))
	// Activity
	http.HandleFunc("/admin/api/v1/SetupActivites", s.handler(s.SetupActivities))
	http.HandleFunc("/admin/api/v1/SkipActivity", s.handler(s.userAdmin(s.SkipActivity)))
	http.HandleFunc("/admin/api/v1/UnskipActivity", s.handler(s.userAdmin(s.UnskipActivity)))
	// Logs
	http.HandleFunc("/admin/api/v1/GetLog", s.handler(s.userAdmin(s.GetLog)))
	http.HandleFunc("/admin/api/v1/GetLogs", s.handler(s.userAdmin(s.GetLogs)))
	http.HandleFunc("/admin/api/v1/GetLogsByEmail", s.handler(s.userAdmin(s.GetLogsByEmail)))
	// Sublogs
	http.HandleFunc("/admin/api/v1/GetUnpaidSublogs", s.handler(s.userAdmin(s.GetUnpaidSublogs)))
	http.HandleFunc("/admin/api/v1/ProcessSublog", s.handler(s.userAdmin(s.ProcessSublog)))
	http.HandleFunc("/admin/api/v1/GetSubscriberSublogs", s.handler(s.userAdmin(s.GetSubscriberSublogs)))
	http.HandleFunc("/admin/api/v1/RefundAndSkipSublog", s.handler(s.userAdmin(s.RefundAndSkipSublog)))
	// Subscriber
	http.HandleFunc("/admin/api/v1/GetHasSubscribed", s.handler(s.userAdmin(s.GetHasSubscribed)))
	http.HandleFunc("/admin/api/v1/GetSubscriber", s.handler(s.userAdmin(s.GetSubscriber)))
	http.HandleFunc("/admin/api/v1/SendCustomerSMS", s.handler(s.userAdmin(s.SendCustomerSMS)))
	// Zone
	// http.HandleFunc("/admin/api/v1/AddGeofence", handler(driverAdmin(s.AddGeofence)))
	// Culture Executions
	http.HandleFunc("/admin/api/v1/GetExecutions", s.handler(s.userAdmin(s.GetExecutions)))
	http.HandleFunc("/admin/api/v1/GetExecution", s.handler(s.userAdmin(s.GetExecution)))
	http.HandleFunc("/admin/api/v1/UpdateExecution", s.handler(s.userAdmin(s.UpdateExecution)))
	// Tasks
	http.HandleFunc("/admin/task/SetupTags", s.handler(s.SetupTags))
	http.HandleFunc("/admin/task/SendPreviewCultureEmail", s.handler(s.SendPreviewCultureEmail))
	http.HandleFunc("/admin/task/SendCultureEmail", s.handler(s.SendCultureEmail))
	http.HandleFunc("/admin/task/CheckPowerSensors", s.handler(s.CheckPowerSensors))
	http.HandleFunc("/admin/task/SendStatsSMS", s.handler(s.SendStatsSMS))
	http.HandleFunc("/admin/task/BackupDatastore", s.handler(s.BackupDatastore))

	http.HandleFunc("/admin/task/ProcessActivity", s.handler(s.ProcessActivity))
	http.HandleFunc("/process-subscription", s.handler(s.ProcessActivity))
	http.HandleFunc("/admin/task/SetupActivites", s.handler(s.SetupActivities))
	// Webhooks
	http.HandleFunc("/admin/webhook/typeform-skip", s.handler(s.TypeformSkip))
	http.HandleFunc("/admin/webhook/twilio-sms", s.handler(s.TwilioSMS))
	http.HandleFunc("/admin/webhook/slack", s.handler(s.Slack))
	// Batch
	http.HandleFunc("/admin/batch/UpdatePhoneNumbers", s.handler(s.UpdatePhoneNumbers))
	//
	http.HandleFunc("/admin/api/v1/Test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

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
		return fmt.Errorf("failed to get sql database client: %+v", err)
	}
	s.serverInfo = &common.ServerInfo{
		ProjectID:           projID,
		IsStandardAppEngine: true,
	}
	return nil
}

func (s *server) userAdmin(f handle) handle {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
		user, err := s.getUserFromRequest(ctx, w, r, log)
		if err != nil {
			return err
		}
		if !user.Admin {
			return errPermissionDenied
		}
		ctx = context.WithValue(ctx, common.ContextUserID, user.ID)
		ctx = context.WithValue(ctx, common.ContextUserEmail, user.Email)
		return f(ctx, w, r, log)
	}
}

func (s *server) driverAdmin(f handle) handle {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
		user, err := s.getUserFromRequest(ctx, w, r, log)
		if err != nil {
			return err
		}
		if !user.Admin {
			return errPermissionDenied
		}
		ctx = context.WithValue(ctx, common.ContextUserID, user.ID)
		ctx = context.WithValue(ctx, common.ContextUserEmail, user.Email)
		return f(ctx, w, r, log)
	}
}

func (s *server) getUserFromRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) (*common.User, *errors.ErrorWithCode) {
	token := r.Header.Get("auth-token")
	if token == "" {
		e := errBadRequest.Annotate("auth-token is empty")
		return nil, &e
	}
	// TODO: use new auth client
	authC, err := auth.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		e := errors.Annotate(err, "failed to get auth.NewClient")
		return nil, &e
	}
	user, err := authC.Verify(token)
	if err != nil {
		e := errors.Annotate(err, "failed to get auth.Verify")
		return nil, &e
	}
	// userold, err := authold.GetUserFromToken(ctx, token)
	// if err != nil {
	// 	e := errors.Annotate(err, "failed to authold.GetUserFromToken")
	// 	return nil, &e
	// }
	// first := ""
	// last := ""
	// name := strings.Title(strings.TrimSpace(userold.Name))
	// lastSpace := strings.LastIndex(name, " ")
	// if lastSpace == -1 {
	// 	first = name
	// } else {
	// 	first = name[:lastSpace]
	// 	last = name[lastSpace:]
	// }
	// user := &common.User{
	// 	FirstName: first,
	// 	LastName:  last,
	// 	Email:     userold.Email,
	// 	PhotoURL:  userold.PhotoURL,
	// 	Admin:     userold.IsAdmin(),
	// }
	return user, nil
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
			http.Redirect(w, r, url, http.StatusMovedPermanently)
			return
		}
		// get context
		ctx := appengine.NewContext(r)
		ctx = context.WithValue(ctx, common.ContextUserID, int64(0))
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
			sharedErr = &shared.Error{
				Code: shared.Code_Success,
			}
		}
		if sharedErr != nil && sharedErr.Code != shared.Code_Success && sharedErr.Code != shared.Code(0) {
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
	GetError() *shared.Error
}

type handle func(context.Context, http.ResponseWriter, *http.Request, *logging.Client) Response

func getDatetime(s string) time.Time {
	if len(s) == 10 {
		s += "T12:12:12.000Z"
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
