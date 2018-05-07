package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/atishpatel/Gigamunch-Backend/config"

	authold "github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/db"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
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
	http.HandleFunc("/admin/api/v1/GetLogsByEmail", handler(systemsAdmin(GetLogsByEmail)))
	// Sublogs
	http.HandleFunc("/admin/api/v1/GetUnpaidSublogs", handler(userAdmin(GetUnpaidSublogs)))
	http.HandleFunc("/admin/api/v1/ProcessSublog", handler(userAdmin(ProcessSublog)))
	http.HandleFunc("/admin/api/v1/GetHasSubscribed", handler(userAdmin(GetHasSubscribed)))
	// Zone
	// http.HandleFunc("/admin/api/v1/AddGeofence", handler(driverAdmin(AddGeofence)))
	//
	http.HandleFunc("/admin/api/v1/Test", test)
	setupTasksHandlers()
}

func setup() error {
	var err error
	projID = os.Getenv("PROJECT_ID")
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
	config := config.GetConfig(ctx)
	err = auth.Setup(ctx, true, projID, httpClient, dbC, config.JWTSecret)
	if err != nil {
		return fmt.Errorf("failed to setup auth: %+v", err)
	}
	// Setup logging
	err = logging.Setup(ctx, true, projID, "admin", nil, dbC)
	if err != nil {
		return fmt.Errorf("failed to setup logging: %+v", err)
	}
	// Setup mail
	err = mail.Setup(ctx, true, projID, config.DripAPIKey, config.DripAccountID, config.MailgunAPIKey, config.MailgunPublicAPIKey)
	if err != nil {
		return fmt.Errorf("failed to setup mail: %+v", err)
	}
	// Setup Sub
	err = sub.Setup(ctx, true, projID, sqlC, dbC)
	if err != nil {
		return fmt.Errorf("failed to setup sub: %+v", err)
	}
	return nil
}

func userAdmin(f handle) handle {
	return func(ctx context.Context, r *http.Request, log *logging.Client) Response {
		user, err := getUserFromRequest(ctx, r, log)
		if err != nil {
			return err
		}
		if !user.IsUserAdmin() {
			return errPermissionDenied
		}
		ctx = context.WithValue(ctx, common.ContextUserID, user.ID)
		ctx = context.WithValue(ctx, common.ContextUserEmail, user.Email)
		return f(ctx, r, log)
	}
}

func driverAdmin(f handle) handle {
	return func(ctx context.Context, r *http.Request, log *logging.Client) Response {
		user, err := getUserFromRequest(ctx, r, log)
		if err != nil {
			return err
		}
		if !user.IsDriverAdmin() {
			return errPermissionDenied
		}
		ctx = context.WithValue(ctx, common.ContextUserID, user.ID)
		ctx = context.WithValue(ctx, common.ContextUserEmail, user.Email)
		return f(ctx, r, log)
	}
}

func systemsAdmin(f handle) handle {
	return func(ctx context.Context, r *http.Request, log *logging.Client) Response {
		user, err := getUserFromRequest(ctx, r, log)
		if err != nil {
			return err
		}
		if !user.IsSystemsAdmin() {
			return errPermissionDenied
		}
		ctx = context.WithValue(ctx, common.ContextUserID, user.ID)
		ctx = context.WithValue(ctx, common.ContextUserEmail, user.Email)
		return f(ctx, r, log)
	}
}

func getUserFromRequest(ctx context.Context, r *http.Request, log *logging.Client) (*common.User, *errors.ErrorWithCode) {
	token := r.Header.Get("auth-token")
	if token == "" {
		e := errBadRequest.Annotate("auth-token is empty")
		return nil, &e
	}
	// TODO: use new auth client
	// authC, err := auth.NewClient(ctx, log)
	// if err != nil {
	// 	return nil, errors.Annotate(err, "failed to get auth.NewClient")
	// }
	// user, err := authC.GetUser(ctx, token)
	// if err != nil {
	// 	return nil, errors.Annotate(err, "failed to get auth.GetUser")
	// }
	userold, err := authold.GetUserFromToken(ctx, token)
	if err != nil {
		e := errors.Annotate(err, "failed to authold.GetUserFromToken")
		return nil, &e
	}
	first := ""
	last := ""
	name := strings.Title(strings.TrimSpace(userold.Name))
	lastSpace := strings.LastIndex(name, " ")
	if lastSpace == -1 {
		first = name
	} else {
		first = name[:lastSpace]
		last = name[lastSpace:]
	}
	user := &common.User{
		FirstName:   first,
		LastName:    last,
		Email:       userold.Email,
		PhotoURL:    userold.PhotoURL,
		Permissions: userold.Permissions,
	}
	return user, nil
}

func handler(f func(context.Context, *http.Request, *logging.Client) Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Access-Control-Allow-Headers, Access-Control-Allow-Origin, auth-token")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		// get context
		ctx := appengine.NewContext(r)
		if !setupDone {
			err = setupWithContext(ctx)
			if err != nil {
				// TODO: Alert but send friendly error back
				log.Fatalf("failed to setup: %+v", err)
				return
			}
		}
		// create logging client
		loggingC, err := logging.NewClient(ctx, r.URL.Path)
		if err != nil {
			errString := fmt.Sprintf("failed to get new logging client: %+v", err)
			logging.Errorf(ctx, errString)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errString))
		}
		// call function
		resp := f(ctx, r, loggingC)
		// Log errors
		if resp == nil {
			return
		}
		sharedErr := resp.GetError()
		if sharedErr == nil {
			sharedErr = &shared.Error{
				Code: shared.Code_Success,
			}
		}
		if sharedErr != nil && sharedErr.Code != shared.Code_Success && sharedErr.Code != shared.Code(0) {
			logging.Errorf(ctx, "request error: %+v", errors.GetErrorWithCode(sharedErr))
			// loggingC.LogRequestError(r, errors.GetErrorWithCode(sharedErr))
			w.WriteHeader(int(sharedErr.Code))
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

type handle func(context.Context, *http.Request, *logging.Client) Response

func getTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return t, errBadRequest.Annotatef("failed to decode time: %+v", err)
	}
	return t, nil
}
