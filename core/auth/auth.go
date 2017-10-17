package auth

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	fb "firebase.google.com/go"
	fba "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	standAppEngine bool
	projID         string
	fbAuth         *fba.Client
	db             common.DB
	jwtKey         []byte
)

var (
	errDatastore      = errors.InternalServerError
	errInternal       = errors.InternalServerError
	errInvalidToken   = errors.SignOutError.WithMessage("Invalid token.")
	errInvalidFBToken = errors.SignOutError.WithMessage("Invalid Firebase token.")
)

// Client is a client for manipulating auth.
type Client struct {
	ctx context.Context
	log *logging.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context) (*Client, error) {
	var err error
	if standAppEngine {
		httpClient := urlfetch.Client(ctx)
		err = setupFBApp(ctx, httpClient)
		if err != nil {
			return nil, err
		}
	}
	if fbAuth == nil {
		return nil, errInternal.Annotate("setup not called")
	}
	log, ok := ctx.Value(common.LoggingKey).(*logging.Client)
	if !ok {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx: ctx,
		log: log,
	}, nil
}

func createSessionToken(ctx context.Context, fbID, name, email, photoURL, provider, firebase string) (*Token, error) {
	var err error
	firstTime := false
	var multiUserSessions []*UserSessions
	keys, err := db.QueryFilter(ctx, kind, 0, 1, "User.AuthID=", fbID, &multiUserSessions)
	if err != nil {
		if err == db.ErrNoSuchEntity() {
			firstTime = true
		} else {
			return nil, errDatastore.WithError(err).Annotate("failed to QueryFilter")
		}
	}
	if len(keys) > 1 {
		return nil, errDatastore.Annotate("two users with same AuthID")
	}
	var userSessions *UserSessions
	var userID int64
	if len(multiUserSessions) == 0 {
		userSessions = new(UserSessions)
	} else {
		userSessions = multiUserSessions[0]
		userID = keys[0].IntID()
	}
	if firstTime {
		firstName, lastName := splitName(name)
		// create UserSessions kind
		userSessions.Provider = provider
		// userSessions.Firsbase = firebase
		userSessions.User = common.User{
			AuthID:      fbID,
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			PhotoURL:    photoURL,
			Permissions: 0,
		}
		// update PhotoURL given by Google so it's higher resolution
		userSessions.User.PhotoURL = strings.Replace(userSessions.User.PhotoURL, "s96-c", "s250-c", -1)
	}
	userSessions.User.ID = userID
	// create the token
	token := &Token{
		User:   userSessions.User,
		IAT:    getIATTime(),
		JTI:    getNewJTI(),
		Expire: GetExpTime(),
	}
	userSessions.TokenIDs = append(userSessions.TokenIDs, TokenID{JTI: token.JTI, Expire: token.Expire})
	var key common.Key
	if firstTime {
		key = db.IncompleteKey(ctx, kind)
	} else {
		key = db.IDKey(ctx, kind, userID)
	}

	key, err = db.Put(ctx, key, userSessions)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotatef("failed to put for userID(%s)", userID)
	}
	// make sure new users have new ids
	token.User.ID = key.IntID()
	return token, nil
}

// GetFromFBToken gets a User and token from a Firebase token.
func (c *Client) GetFromFBToken(ctx context.Context, fbToken string) (*common.User, string, error) {
	if fbToken == "" {
		return nil, "", errInvalidFBToken.Annotate("FBToken is empty")
	}
	fbTKN, err := fbAuth.VerifyIDToken(fbToken)
	if err != nil || time.Since(time.Unix(fbTKN.IssuedAt, 0)) > 15*time.Minute {
		return nil, "", errInvalidFBToken.WithError(err).Annotate("FBToken is too old")
	}
	claims := fbTKN.Claims
	picture := claims["picture"].(string)
	name := claims["name"].(string)
	email, ok := claims["email"].(string)
	if !ok {

	}
	provider, ok := claims["sign_in_provider"].(string)
	firebase, ok := claims["firebase"].(string)
	c.log.Debugf(ctx, "claims: %+v", claims)
	token, err := createSessionToken(ctx, fbTKN.UID, name, email, picture, provider, firebase)
	if err != nil {
		return nil, "", errors.Wrap("failed to create session token", err)
	}
	// TODO: log
	jwtString, err := token.JWTString()
	if err != nil {
		return nil, "", errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Failed to encode user."}
	}
	return &token.User, jwtString, nil
}

// GetUser gets a User. If the token is fresh, it doesn't make a database call.
func (c *Client) GetUser(ctx context.Context, tkn string) (*common.User, error) {
	// TODO: implement
	return nil, nil
}

// Refresh returns a fresh token and invalidates the old one.
func (c *Client) Refresh(ctx context.Context, tkn string) (string, error) {
	// TODO: implement
	return "", nil
}

// CreateReq is a req for Create.
type CreateReq struct {
	AuthID          string
	Email           string
	FirstName       string
	LastName        string
	PaymentProvider common.PaymentProvider
}

// Create creates a new User.
func (c *Client) Create(ctx context.Context, req *CreateReq) (*common.User, error) {
	return nil, nil
}

// func (c *Client) LinkViaEmail(ctx context.Context, email string, fbID string) (*common.User, error) {
// 	return nil, nil
// }

// func (c *Client) LinkViaUserID(ctx context.Context, email string, userID int64, fbID string) (*common.User, error) {
// 	return nil, nil
// }

// Update updates a user.
func (c *Client) Update(ctx context.Context, userID int64, ops ...func(*common.User)) (*common.User, error) {
	return nil, nil
}

// InvalidateSessions invalidates all sessions for an User.
func (c *Client) InvalidateSessions(ctx context.Context, userID int64) error {
	return nil
}

// GetExpTime returns the time a token should expire from now
func GetExpTime() time.Time {
	return time.Now().UTC().Add(time.Hour * 24 * 60)
}

func getNewJTI() int32 {
	return rand.Int31()
}

func getIATTime() time.Time {
	return time.Now().UTC()
}

func splitName(name string) (string, string) {
	first := "-"
	last := "-"
	nameStripper := strings.NewReplacer(".", "", "Mr ", "", "Ms ", "", "the", "", "Dr ", "")
	nameArray := strings.Split(strings.Title(strings.TrimSpace(nameStripper.Replace(name))), " ")
	if len(nameArray) >= 2 {
		last = nameArray[len(nameArray)-1]
	}
	if len(nameArray) >= 1 {
		first = nameArray[0]
	}
	return first, last
}

// Setup sets up auth.
func Setup(ctx context.Context, standardAppEngine bool, projectID string, httpClient *http.Client, dbC common.DB, jwtSecret string) error {
	rand.Seed(time.Now().Unix())
	var err error
	standAppEngine = standardAppEngine
	projID = projectID
	if !standAppEngine {
		err = setupFBApp(ctx, httpClient)
		if err != nil {
			return err
		}
	}
	if dbC == nil {
		return fmt.Errorf("db cannot be nil for sub")
	}
	db = dbC
	if jwtSecret == "" {
		return fmt.Errorf("jwt secret is empty")
	}
	jwtKey = []byte(jwtSecret)
	return nil
}

func setupFBApp(ctx context.Context, httpClient *http.Client) error {
	var ops []option.ClientOption
	if httpClient != nil {
		ops = append(ops, option.WithHTTPClient(httpClient))
	}
	if appengine.IsDevAppServer() {
		ops = append(ops, option.WithCredentialsFile("../private/gitkit_cert.json"))
	}
	var err error
	fbApp, err := fb.NewApp(ctx, &fb.Config{ProjectID: projID}, ops...)
	if err != nil {
		return err
	}
	fbAuth, err = fbApp.Auth()
	if err != nil {
		return err
	}
	return nil
}
