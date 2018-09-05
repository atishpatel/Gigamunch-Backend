package lead

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/mail"
	"github.com/atishpatel/Gigamunch-Backend/core/message"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"google.golang.org/appengine"
)

// Errors
var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
)

// Client is a client for manipulating activity.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	db         common.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	if dbC == nil {
		return nil, fmt.Errorf("db cannot be nil for sub")
	}
	if serverInfo == nil {
		return nil, errInternal.Annotate("failed to get server info")
	}
	return &Client{
		ctx:        ctx,
		log:        log,
		db:         dbC,
		serverInfo: serverInfo,
	}, nil
}

// Create creates or updates a business lead.
func (c *Client) Create(email string) error {
	lead := &Lead{
		Email:           email,
		CreatedDatetime: time.Now(),
	}
	key := c.db.NameKey(c.ctx, kind, email)

	_, err := c.db.Put(c.ctx, key, lead)
	if err != nil {
		return errDatastore.WithError(err)
	}
	if !appengine.IsDevAppServer() && !strings.Contains(email, "@test.com") {
		messageC := message.New(c.ctx)
		_ = messageC.SendAdminSMS("6153975516", fmt.Sprintf("New subscriber lead. \nEmail: %s", email))
	}

	mailC, err := mail.NewClient(c.ctx, c.log, c.serverInfo)
	if err != nil {
		c.log.Errorf(c.ctx, "failed to mail.NewClient: %+v", err)
	}
	err = mailC.LeftEmail(email, "", "")
	if err != nil {
		return errors.Annotate(err, "failed to mail.LeftEmail")
	}
	return nil
}
