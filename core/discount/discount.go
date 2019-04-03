package discount

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"

	"github.com/jmoiron/sqlx"
)

const (
	// DateFormat is the expected format of date in activity.
	DateFormat                     = "2006-01-02" // "Jan 2, 2006"
	selectActivityByEmailStatement = "SELECT * FROM activity WHERE date=? AND email=?"
)

var (
	errDatastore             = errors.InternalServerError
	errNoSuchEntityDatastore = errors.NotFoundError
	errInternal              = errors.InternalServerError
	errInvalidParameter      = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	db         common.DB
	sqlDB      *sqlx.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	if sqlC == nil {
		return nil, errInternal.Annotate("failed to get sql client")
	}
	if dbC == nil {
		return nil, fmt.Errorf("failed to get db")
	}
	if serverInfo == nil {
		return nil, errInternal.Annotate("failed to get server info")
	}
	return &Client{
		ctx:        ctx,
		log:        log,
		sqlDB:      sqlC,
		db:         dbC,
		serverInfo: serverInfo,
	}, nil
}
