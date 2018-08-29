package activity

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

var (
	sqlDB *sqlx.DB
	db    common.DB
)

// Errors
var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
)

// Client is a client for manipulating activity.
type Client struct {
	ctx context.Context
	log *logging.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx: ctx,
		log: log,
	}, nil
}

// Setup sets up the logging package.
func Setup(ctx context.Context, sqlC *sqlx.DB, dbC common.DB) error {
	if sqlC == nil {
		return fmt.Errorf("sqlDB cannot be nil for sub")
	}
	sqlDB = sqlC
	if dbC == nil {
		return fmt.Errorf("db cannot be nil for sub")
	}
	db = dbC
	return nil
}
