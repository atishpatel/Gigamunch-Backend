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

	"github.com/atishpatel/Gigamunch-Backend/core/serverhelper"

	"github.com/jmoiron/sqlx"

	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/common"
	"google.golang.org/appengine"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

// server is the Subscriber server service.
type server struct {
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

func main() {
	s := new(server)
	err := s.setup()
	if err != nil {
		log.Fatal("failed to setup", err)
	}
	// Activity
	http.HandleFunc("/sub/api/v1/GetExecutions", s.handler(s.GetExecutions))
	http.HandleFunc("/sub/api/v1/GetExecution", s.handler(s.GetExecution))
	// Test
	http.HandleFunc("/sub/api/v1/Test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("success"))
	})
	appengine.Main()
}

// Setup sets up the server.
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

// isSubscriber checks if user is subscriber.
func (s *server) isSubscriber(f handle) handle {
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

// handler encodes response
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
		ctx := r.Context()
		ctx = context.WithValue(ctx, common.ContextUserID, "")
		ctx = context.WithValue(ctx, common.ContextUserEmail, "")
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
			log.RequestError(r, errors.GetErrorWithCode(sharedErr))
			w.WriteHeader(int(sharedErr.Code))
			// Wrap error in ErrorOnlyResp
			if _, ok := resp.(errors.ErrorWithCode); ok {
				resp = &pbcommon.ErrorOnlyResp{
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

// decodeRequest decodes a request into a struct.
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

// handle is the handle for api request.
type handle func(context.Context, http.ResponseWriter, *http.Request, *logging.Client) Response

func getDatetime(s string) time.Time {
	return serverhelper.GetDatetime(s)
}