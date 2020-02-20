package auth

import (
	"context"
	"strings"

	fb "firebase.google.com/go"
	fba "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"google.golang.org/appengine"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

const (
	userIDClaim = "giga_user_id"
	delim       = ";%~&~%;"
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
	fbAuth     *fba.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	var err error
	fbAuth, err := setupFBApp(ctx, serverInfo.ProjectID)
	if err != nil {
		return nil, err
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
		fbAuth:     fbAuth,
	}, nil
}

// CreateUserIfNoExist creates a user if it does not already exist.
func (c *Client) CreateUserIfNoExist(email, password, name string) error {
	_, err := c.fbAuth.GetUserByEmail(c.ctx, email)
	if err != nil {
		utc := &fba.UserToCreate{}
		utc.Email(strings.ToLower(email))
		utc.Password((password))
		utc.DisplayName(name)
		_, err = c.fbAuth.CreateUser(c.ctx, utc)
		if err != nil {
			return errInternal.WithError(err).Annotate("failed to fbauth.CreateUser")
		}
	}
	return nil
}

// Verify verifies the token.
func (c *Client) Verify(token string) (*common.User, error) {
	if token == "" {
		return nil, errInvalidFBToken.Annotate("token is empty")
	}
	tkn, err := c.fbAuth.VerifyIDToken(c.ctx, token)
	if err != nil {
		return nil, errInvalidFBToken.WithError(err).Annotate("failed to fbauth.VerifyIDToken")
	}
	claims := tkn.Claims
	getString := func(c map[string]interface{}, key string) string {
		tmp := c[key]
		if tmp == nil {
			return ""
		}
		return tmp.(string)
	}
	picture := getString(claims, "picture")
	name := getString(claims, "name")
	email := getString(claims, "email")
	if email == "" {
		return nil, errInvalidFBToken.WithMessage("User must have email.").Annotate("firebase token does not have email")
	}
	userID := getString(claims, userIDClaim)
	var admin bool
	adminTmp, ok := claims["admin"]
	if ok {
		admin = adminTmp.(bool)
	}
	// var activeSubscriber bool
	// activeSubscriberTmp, ok := claims["active_subscriber"]
	// if ok {
	// 	activeSubscriber = activeSubscriberTmp.(bool)
	// }
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
		Email:     strings.ToLower(email),
		PhotoURL:  picture,
		Admin:     admin,
		// ActiveSubscriber: activeSubscriber,
	}
	return user, nil
}

// // SetActiveSubscriber makes a user an admin
// func (c *Client) SetActiveSubscriber(authIDOrEmail string, active bool) error {
// 	return c.AddCustomClaim(authIDOrEmail, "active_subscriber", active)
// }

// SetAdmin makes a user an admin
func (c *Client) SetAdmin(authIDOrEmail string, active bool) error {
	return c.AddCustomClaim(authIDOrEmail, "admin", active)
}

// UpdateUserID sets the user id for a token.
func (c *Client) UpdateUserID(authIDOrEmail string, userID string) error {
	return c.AddCustomClaim(authIDOrEmail, userIDClaim, userID)
}

// UpdateUser updates the user's ID and name in token.
func (c *Client) UpdateUser(authID, userID, email, firstName, lastName string) error {
	var userRecord *fba.UserRecord
	var err error
	if strings.Contains(authID, "@") {
		userRecord, err = c.fbAuth.GetUserByEmail(c.ctx, authID)
	} else {
		userRecord, err = c.fbAuth.GetUser(c.ctx, authID)
	}
	if err != nil {
		return errInvalidArgument.WithMessage("Invalid email")
	}
	userToUpdate := new(fba.UserToUpdate)
	if firstName != "" {
		userToUpdate.DisplayName(firstName + delim + lastName)
	}
	if email != "" {
		userToUpdate.Email(email)
	}
	if userID != "" {
		claims := userRecord.CustomClaims
		if claims == nil {
			claims = make(map[string]interface{})
		}
		claims[userIDClaim] = userID
		userToUpdate.CustomClaims(claims)
	}
	_, err = c.fbAuth.UpdateUser(c.ctx, userRecord.UID, userToUpdate)
	if err != nil {
		return errInternal.WithError(err).Annotate("failed to fbAuth.UpdateUser")
	}
	return nil
}

// AddCustomClaim adds a custom claim to a user token.
func (c *Client) AddCustomClaim(authIDOrEmail, key string, value interface{}) error {
	var userRecord *fba.UserRecord
	var err error
	if strings.Contains(authIDOrEmail, "@") {
		userRecord, err = c.fbAuth.GetUserByEmail(c.ctx, authIDOrEmail)
	} else {
		userRecord, err = c.fbAuth.GetUser(c.ctx, authIDOrEmail)
	}
	if err != nil {
		return errInvalidArgument.WithError(err).Annotate("Invalid parameter.")
	}
	claims := userRecord.CustomClaims
	if claims == nil {
		claims = make(map[string]interface{})
	}
	claims[key] = value
	userToUpdate := new(fba.UserToUpdate)
	userToUpdate.CustomClaims(claims)
	_, err = c.fbAuth.UpdateUser(c.ctx, userRecord.UID, userToUpdate)
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

func setupFBApp(ctx context.Context, projectID string) (*fba.Client, error) {
	var ops []option.ClientOption
	var err error
	if appengine.IsDevAppServer() {
		ops = append(ops, option.WithCredentialsFile("../private/gitkit_cert.json"))
	}
	fbApp, err := fb.NewApp(ctx, nil, ops...)
	if err != nil {
		return nil, err
	}
	fbAuth, err := fbApp.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return fbAuth, nil
}
