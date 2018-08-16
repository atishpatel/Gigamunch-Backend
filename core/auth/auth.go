package auth

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	fb "firebase.google.com/go"
	fba "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

var (
	once   sync.Once
	fbAuth *fba.Client
	jwtKey []byte
)

var (
	errDatastore      = errors.InternalServerError
	errInternal       = errors.InternalServerError
	errInvalidToken   = errors.SignOutError.WithMessage("Invalid token.")
	errInvalidFBToken = errors.SignOutError.WithMessage("Invalid Firebase token.")
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
	}
	once.Do(func() {
		setup(ctx, serverInfo)
	})
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

func setup(ctx context.Context, serverInfo *common.ServerInfo) {
	rand.Seed(time.Now().Unix())
	if !serverInfo.IsStandardAppEngine {
		httpClient := http.DefaultClient
		err := setupFBApp(ctx, httpClient, serverInfo.ProjectID)
		if err != nil {
			// TODO:
			// return err
		}
	}
	// Get JWT Secret
	config := config.GetConfig(ctx)
	jwtKey = []byte(config.JWTSecret)
}

func createSessionToken(ctx context.Context, db common.DB, fbID, name, email, photoURL, provider, firebase string) (*Token, error) {
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
		return nil, "", errInvalidFBToken.WithMessage("User must have email.").Annotate("firebase token does not have email")
	}
	provider, ok := claims["sign_in_provider"].(string)
	firebase, ok := claims["firebase"].(string)
	c.log.Debugf(ctx, "claims: %+v", claims)
	token, err := createSessionToken(ctx, c.db, fbTKN.UID, name, email, picture, provider, firebase)
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
