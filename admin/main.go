package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/message"

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	"github.com/jmoiron/sqlx"

	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	pbcommon "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon"
	"google.golang.org/appengine"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	serverInfo *common.ServerInfo
	db         *db.Client
	sqlDB      *sqlx.DB
}

var (
	errPermissionDenied = errors.PermissionDeniedError
	// errUnauthenticated  = errors.UnauthenticatedError
	errBadRequest    = errors.BadRequestError
	errInternalError = errors.InternalServerError
)

func main() {
	s := new(server)
	err := s.setup()
	if err != nil {
		log.Fatal("failed to setup", err)
	}
	// **********************
	// Auth
	// **********************
	http.HandleFunc("/admin/api/v1/SetAdmin", s.handler(s.userAdmin(s.SetAdmin)))
	// **********************
	// Subscriber
	// **********************
	http.HandleFunc("/admin/api/v1/ActivateSubscriber", s.handler(s.userAdmin(s.ActivateSubscriber)))
	http.HandleFunc("/admin/api/v1/DeactivateSubscriber", s.handler(s.userAdmin(s.DeactivateSubscriber)))
	http.HandleFunc("/admin/api/v1/ReplaceSubscriberEmail", s.handler(s.userAdmin(s.ReplaceSubscriberEmail)))
	// **********************
	// Activity
	// **********************
	http.HandleFunc("/admin/api/v1/SetupActivites", s.handler(s.SetupActivities))
	http.HandleFunc("/admin/api/v1/SkipActivity", s.handler(s.userAdmin(s.SkipActivity)))
	http.HandleFunc("/admin/api/v1/UnskipActivity", s.handler(s.userAdmin(s.UnskipActivity)))
	http.HandleFunc("/admin/api/v1/RefundActivity", s.handler(s.userAdmin(s.RefundActivity)))
	http.HandleFunc("/admin/api/v1/RefundAndSkipActivity", s.handler(s.userAdmin(s.RefundAndSkipActivity)))
	// **********************
	// Logs
	// **********************
	http.HandleFunc("/admin/api/v1/GetLog", s.handler(s.userAdmin(s.GetLog)))
	http.HandleFunc("/admin/api/v1/GetLogs", s.handler(s.userAdmin(s.GetLogs)))
	http.HandleFunc("/admin/api/v1/GetLogsByEmail", s.handler(s.userAdmin(s.GetLogsByEmail)))
	http.HandleFunc("/admin/api/v1/GetLogsByExecution", s.handler(s.userAdmin(s.GetLogsByExecution)))
	// **********************
	// Sublogs
	// **********************
	http.HandleFunc("/admin/api/v1/GetUnpaidSublogs", s.handler(s.userAdmin(s.GetUnpaidSublogs)))
	http.HandleFunc("/admin/api/v1/ProcessSublog", s.handler(s.userAdmin(s.ProcessSublog)))
	http.HandleFunc("/admin/api/v1/GetSubscriberSublogs", s.handler(s.userAdmin(s.GetSubscriberSublogs)))
	// **********************
	// Subscriber
	// **********************
	http.HandleFunc("/admin/api/v1/GetHasSubscribed", s.handler(s.userAdmin(s.GetHasSubscribed)))
	http.HandleFunc("/admin/api/v1/GetSubscriber", s.handler(s.userAdmin(s.GetSubscriber)))
	http.HandleFunc("/admin/api/v1/SendCustomerSMS", s.handler(s.userAdmin(s.SendCustomerSMS)))
	http.HandleFunc("/admin/api/v1/UpdateDrip", s.handler(s.userAdmin(s.UpdateDrip)))
	// Zone
	http.HandleFunc("/admin/api/v1/UpdateGeofence", s.handler(s.userAdmin(s.UpdateGeofence)))
	// **********************
	// Culture Executions
	// **********************
	http.HandleFunc("/admin/api/v1/GetExecutions", s.handler(s.userAdmin(s.GetExecutions)))
	http.HandleFunc("/admin/api/v1/GetExecution", s.handler(s.userAdmin(s.GetExecution)))
	http.HandleFunc("/admin/api/v1/UpdateExecution", s.handler(s.userAdmin(s.UpdateExecution)))
	// **********************
	// Tasks
	// **********************
	http.HandleFunc("/admin/task/SetupTags", s.handler(s.SetupTags))
	http.HandleFunc("/admin/task/SendPreviewCultureEmail", s.handler(s.SendPreviewCultureEmail))
	http.HandleFunc("/admin/task/SendCultureEmail", s.handler(s.SendCultureEmail))
	http.HandleFunc("/admin/task/CheckPowerSensors", s.handler(s.CheckPowerSensors))
	http.HandleFunc("/admin/task/SendStatsSMS", s.handler(s.SendStatsSMS))
	http.HandleFunc("/admin/task/BackupDatastore", s.handler(s.BackupDatastore))

	http.HandleFunc("/admin/task/ProcessActivity", s.handler(s.ProcessActivity))
	http.HandleFunc("/process-subscription", s.handler(s.ProcessActivity))
	http.HandleFunc("/admin/task/SetupActivites", s.handler(s.SetupActivities))
	// **********************
	// Webhooks
	// **********************
	http.HandleFunc("/admin/webhook/typeform-skip", s.handler(s.TypeformSkip))
	http.HandleFunc("/admin/webhook/twilio-sms", s.handler(s.TwilioSMS))
	http.HandleFunc("/admin/webhook/slack", s.handler(s.Slack))
	// **********************
	// Batch
	// **********************
	http.HandleFunc("/admin/batch/UpdateSubs", s.handler(s.UpdateSubs))
	// http.HandleFunc("/admin/batch/MigrateToNewSubscribersStruct", s.handler(s.MigrateToNewSubscribersStruct))
	//
	http.HandleFunc("/admin/api/v1/Test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("success"))
	})
	appengine.Main()
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
		sqlConnectionString = "server:gigamunch@/gigamunch"
	}
	s.sqlDB, err = sqlx.Connect("mysql", sqlConnectionString+"?collation=utf8mb4_general_ci&parseTime=true")
	if err != nil {
		return fmt.Errorf("failed to get sql database client: %+v", err)
	}
	s.db, err = db.NewClient(context.Background(), projID, nil)
	if err != nil {
		return fmt.Errorf("failed to get database client: %+v", err)
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
			return errors.GetErrorWithCode(err)
		}
		if !user.Admin {
			return errPermissionDenied
		}
		ctx = context.WithValue(ctx, common.ContextUserID, user.ID)
		ctx = context.WithValue(ctx, common.ContextUserEmail, user.Email)
		log.SetContext(ctx)
		return f(ctx, w, r, log)
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
		return nil, errors.Annotate(err, "failed to get auth.Verify")
	}
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
		ctx = context.WithValue(ctx, common.ContextUserID, "")
		ctx = context.WithValue(ctx, common.ContextUserEmail, "")
		defer func() {
			// handle panic, recover
			if r := recover(); r != nil {
				errString := fmt.Sprintf("PANICKING: %+v", r)
				logging.Errorf(ctx, errString)
				messageC := message.New(ctx)
				_ = messageC.SendAdminSMS(message.EmployeeNumbers.OnCallDeveloper(), errString)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf("{\"code\":500,\"message\":\"Woops! Something went wrong. Please try again later.\"}")))
			}
		}()
		// create logging client
		log, err := logging.NewClient(ctx, "admin", r.URL.Path, s.db, s.serverInfo)
		if err != nil {
			errString := fmt.Sprintf("failed to get new logging client: %+v", err)
			logging.Errorf(ctx, errString)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errString))
			return
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
			log.RequestError(r, errors.GetErrorWithCode(sharedErr))
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
	return serverhelper.DecodeRequest(ctx, r, v)
}

func failedToDecode(err error) *pbcommon.ErrorOnlyResp {
	return serverhelper.FailedToDecode(err)
}

// Response is a response to a rpc call. All responses contain an error.
type Response interface {
	GetError() *pbcommon.Error
}

type handle func(context.Context, http.ResponseWriter, *http.Request, *logging.Client) Response

func getDatetime(s string) time.Time {
	return serverhelper.GetDatetime(s)
}
