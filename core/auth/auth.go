package auth

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	fb "firebase.google.com/go"
	fba "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

const (
	userIDClaim = "userID"
	delim       = ";%~&~%;"
)

var (
	once   sync.Once
	fbAuth *fba.Client
)

var (
	errDatastore       = errors.InternalServerError
	errInternal        = errors.InternalServerError
	errInvalidArgument = errors.BadRequestError
	errInvalidToken    = errors.SignOutError.WithMessage("Invalid token.")
	errInvalidFBToken  = errors.SignOutError.WithMessage("Invalid Firebase token.")
)

// Client is a client for manipulating auth.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	db         common.DB
	sqlDB      *sqlx.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	var err error
	if serverInfo.IsStandardAppEngine {
		httpClient := urlfetch.Client(ctx)
		err = setupFBApp(ctx, httpClient, serverInfo.ProjectID)
		if err != nil {
			return nil, err
		}
	} else {
		once.Do(func() {
			if !serverInfo.IsStandardAppEngine {
				httpClient := http.DefaultClient
				err = setupFBApp(ctx, httpClient, serverInfo.ProjectID)
			}
		})
		if err != nil {
			return nil, err
		}
	}
	if fbAuth == nil {
		return nil, errInternal.Annotate("setup not called")
	}
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx:        ctx,
		log:        log,
		db:         dbC,
		sqlDB:      sqlC,
		serverInfo: serverInfo,
	}, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}

// Verify verifies the token.
func (c *Client) Verify(token string) (*common.User, error) {
	if token == "" {
		return nil, errInvalidFBToken.Annotate("token is empty")
	}
	tkn, err := fbAuth.VerifyIDToken(c.ctx, token)
	if err != nil {
		return nil, errInvalidFBToken.WithError(err).Annotate("failed to fbauth.VerifyIDToken")
	}
	claims := tkn.Claims
	picture := claims["picture"].(string)
	name := claims["name"].(string)
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errInvalidFBToken.WithMessage("User must have email.").Annotate("firebase token does not have email")
	}
	var userID int64
	userIDTmp, ok := claims[userIDClaim]
	if ok {
		userID, _ = strconv.ParseInt(userIDTmp.(string), 2, 64)
	}
	admin := claims["admin"].(bool)
	nameSplit := strings.Split(name, delim)
	var firstName, lastName string
	if len(nameSplit) >= 1 {
		firstName = nameSplit[0]
	}
	if len(nameSplit) >= 2 {
		lastName = nameSplit[1]
	}
	user := &common.User{
		ID:        userID,
		AuthID:    tkn.UID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		PhotoURL:  picture,
		Admin:     admin,
	}
	return user, nil
}

// MakeAdmin makes a user an admin
func (c *Client) MakeAdmin(userEmail string) error {
	return c.AddCustomClaim(userEmail, "admin", true)
}

// UpdateUserID sets the user id for a token.
func (c *Client) UpdateUserID(authID string, userID int64) error {
	return c.AddCustomClaim(authID, userIDClaim, userID)
}

// UpdateUserName updates the user name in token.
func (c *Client) UpdateUserName(authID, firstName, lastName string) error {
	var userRecord *fba.UserRecord
	var err error
	if strings.Contains(authID, "@") {
		userRecord, err = fbAuth.GetUserByEmail(c.ctx, authID)
	} else {
		userRecord, err = fbAuth.GetUser(c.ctx, authID)
	}
	if err != nil {
		return errInvalidArgument.WithMessage("Invalid email")
	}
	userToUpdate := new(fba.UserToUpdate)
	userToUpdate.DisplayName(firstName + delim + lastName)
	_, err = fbAuth.UpdateUser(c.ctx, userRecord.UID, userToUpdate)
	if err != nil {
		return errInternal.WithError(err).Annotate("failed to fbAuth.UpdateUser")
	}
	return nil
}

// AddCustomClaim adds a custom claim to a user token.
func (c *Client) AddCustomClaim(userIDOrEmail, key string, value interface{}) error {
	var userRecord *fba.UserRecord
	var err error
	if strings.Contains(userIDOrEmail, "@") {
		userRecord, err = fbAuth.GetUserByEmail(c.ctx, userIDOrEmail)
	} else {
		userRecord, err = fbAuth.GetUser(c.ctx, userIDOrEmail)
	}
	if err != nil {
		return errInvalidArgument.WithMessage("Invalid email")
	}
	claims := userRecord.CustomClaims
	claims[key] = value
	userToUpdate := new(fba.UserToUpdate)
	userToUpdate.CustomClaims(claims)
	_, err = fbAuth.UpdateUser(c.ctx, userRecord.UID, userToUpdate)
	if err != nil {
		return errInternal.WithError(err).Annotate("failed to fbAuth.UpdateUser")
	}
	return nil
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

func setupFBApp(ctx context.Context, httpClient *http.Client, projectID string) error {
	var ops []option.ClientOption
	if httpClient != nil {
		ops = append(ops, option.WithHTTPClient(httpClient))
	}
	if appengine.IsDevAppServer() {
		ops = append(ops, option.WithCredentialsFile("../private/gitkit_cert.json"))
	}
	var err error
	fbApp, err := fb.NewApp(ctx, &fb.Config{ProjectID: projectID}, ops...)
	if err != nil {
		return err
	}
	fbAuth, err = fbApp.Auth(ctx)
	if err != nil {
		return err
	}
	return nil
}
