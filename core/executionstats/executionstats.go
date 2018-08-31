package executionstats

import (
	"context"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	sqlDB      *sqlx.DB
	db         common.DB
	serverInfo *common.ServerInfo
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
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

// Get gets a culture ExecutionStats.
func (c *Client) Get(id int64) (*ExecutionStats, error) {
	key := c.db.IDKey(c.ctx, Kind, id)
	exe := new(ExecutionStats)
	err := c.db.Get(c.ctx, key, exe)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to get")
	}
	return exe, nil
}

// GetAll gets all ExecutionStats ordered by created datetime.
func (c *Client) GetAll(start, limit int) ([]*ExecutionStats, error) {
	var exes []*ExecutionStats
	_, err := c.db.QueryOrdered(c.ctx, Kind, start, limit, "-CreatedDatetime", exes)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryOrdered")
	}
	return exes, nil
}

// Update updates an ExecutionStats.
func (c *Client) Update(exe *ExecutionStats) error {
	key := c.db.IDKey(c.ctx, Kind, exe.ID)
	key, err := c.db.Put(c.ctx, key, exe)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to put")
	}
	if exe.ID == 0 {
		exe.ID = key.IntID()
		_, err = c.db.Put(c.ctx, key, exe)
		if err != nil {
			return errDatastore.WithError(err).Annotate("failed to put")
		}
	}
	return nil
}
