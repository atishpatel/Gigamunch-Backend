package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/sub"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/admin"
	"github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

var (
	projID    string
	dbC       *db.Client
	sqlC      *sqlx.DB
	setupDone = false
)

var (
	errPermissionDenied = errors.PermissionDeniedError
	errUnauthenticated  = errors.UnauthenticatedError
	errBadRequest       = errors.BadRequestError
)

func init() {
	err := setup()
	if err != nil {
		log.Fatal("failed to setup", err)
	}
	// Auth
	http.HandleFunc("/admin/api/v1/Login", handler(Login))
	http.HandleFunc("/admin/api/v1/Refresh", handler(Refresh))
	// Logs
	http.HandleFunc("/admin/api/v1/GetLog", handler(systemsAdmin(GetLog)))
	http.HandleFunc("/admin/api/v1/GetLogs", handler(systemsAdmin(GetLogs)))
	//
	http.HandleFunc("/admin/api/v1/Test", test)
}

func setup() error {
	var err error
	projID = os.Getenv("PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "PROJECT_ID"`)
	}
	sqlConnectionString := os.Getenv("MYSQL_CONNECTION")
	if sqlConnectionString == "" {
		log.Fatal(`You need to set the environment variable "MYSQL_CONNECTION"`)
	}
	// Setup sql db
	sqlC, err = sqlx.Connect("mysql", sqlConnectionString)
	if err != nil {
		return fmt.Errorf("failed to get sql database client: %+v", err)
	}
	return nil
}

// setupWithContext can be called in main for flex but needs to be called with each method on standard.
func setupWithContext(ctx context.Context) error {
	var err error
	dbC, err = db.NewClient(ctx, projID, nil)
	if err != nil {
		return fmt.Errorf("failed to get database client: %+v", err)
	}
	// Setup auth
	httpClient := urlfetch.Client(ctx)
	err = auth.Setup(ctx, projID, httpClient, dbC, "TODO: get config")
	if err != nil {
		return fmt.Errorf("failed to setup auth: %+v", err)
	}
	// Setup logging
	err = logging.Setup(ctx, projID, "admin", nil, dbC)
	if err != nil {
		return fmt.Errorf("failed to setup logging: %+v", err)
	}
	// Setup Sub
	err = sub.Setup(ctx, sqlC, dbC)
	if err != nil {
		return fmt.Errorf("failed to setup sub: %+v", err)
	}
	return nil
}

func userAdmin(f func(context.Context, *http.Request) Response) func(context.Context, *http.Request) Response {
	return func(ctx context.Context, r *http.Request) Response {
		user, err := getUserFromRequest(ctx, r)
		if !err.IsNil() {
			return err
		}
		if !user.IsUserAdmin() {
			return errPermissionDenied
		}
		return f(ctx, r)
	}
}

func driverAdmin(f func(context.Context, *http.Request) Response) func(context.Context, *http.Request) Response {
	return func(ctx context.Context, r *http.Request) Response {
		user, err := getUserFromRequest(ctx, r)
		if !err.IsNil() {
			return err
		}
		if !user.IsDriverAdmin() {
			return errPermissionDenied
		}
		return f(ctx, r)
	}
}

func systemsAdmin(f func(context.Context, *http.Request) Response) func(context.Context, *http.Request) Response {
	return func(ctx context.Context, r *http.Request) Response {
		user, err := getUserFromRequest(ctx, r)
		if !err.IsNil() {
			return err
		}
		if !user.IsSystemsAdmin() {
			return errPermissionDenied
		}
		return f(ctx, r)
	}
}

func getUserFromRequest(ctx context.Context, r *http.Request) (*common.User, errors.ErrorWithCode) {
	token := r.Header.Get("auth-token")
	if token == "" {
		return nil, errBadRequest.Annotate("auth-token is empty")
	}
	authC, err := auth.NewClient(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "failed to get auth.NewClient")
	}
	user, err := authC.GetUser(ctx, token)
	if err != nil {
		return nil, errors.Annotate(err, "failed to get auth.GetUser")
	}
	return user, errors.NoError
}

func handler(f func(context.Context, *http.Request) Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := appengine.NewContext(r)
		if !setupDone {
			err = setupWithContext(ctx)
			if err != nil {
				// TODO: Alert but send friendly error back
				log.Fatal("failed to setup: %+v", err)
				return
			}
		}
		loggingC, err := logging.NewClient(ctx, r.URL.Path)
		if err != nil {
			errString := fmt.Sprintf("failed to get new logging client: %+v", err)
			logging.Errorf(ctx, errString)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errString))
		}
		ctx = context.WithValue(ctx, common.LoggingKey, loggingC)

		// call function
		resp := f(ctx, r)
		// Log errors
		sharedErr := resp.GetError()
		if sharedErr != nil && sharedErr.Code != shared.Code_Success {
			loggingC.LogRequestError(r, errors.GetErrorWithCode(sharedErr))
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

// Response is a response to a rpc call. All responses contain an error.
type Response interface {
	GetError() *shared.Error
}

func failedToDecode(err error) *pb.ErrorOnlyResp {
	return &pb.ErrorOnlyResp{
		Error: errBadRequest.WithError(err).Annotate("failed to decode").SharedError(),
	}
}

func closeRequestBody(r *http.Request) {
	_ = r.Body.Close()
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("success"))
}
