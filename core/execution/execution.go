package execution

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

var (
	errDatastore  = errors.InternalServerError
	errInternal   = errors.InternalServerError
	errBadRequest = errors.BadRequestError
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
	if sqlC == nil {
		return nil, fmt.Errorf("sqlDB cannot be nil")
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
		db:         dbC,
		sqlDB:      sqlC,
		serverInfo: serverInfo,
	}, nil
}

// Get gets a culture Execution.
func (c *Client) Get(idOrDate string) (*Execution, error) {
	var err error
	exe := new(Execution)
	var key common.Key
	if strings.Contains(idOrDate, "-") {
		var exes []*Execution
		_, err = c.db.QueryFilter(c.ctx, Kind, 0, 1, "Date=", idOrDate, &exes)
		if err != nil {
			return nil, errDatastore.WithError(err).Annotate("faield to query")
		}
		if len(exes) != 1 {
			return nil, errBadRequest.Annotatef("num exes found: %d", len(exes))
		}
		exe = exes[0]
	} else {
		id, err := strconv.ParseInt(idOrDate, 10, 64)
		if err != nil {
			return nil, errors.Annotate(err, "failed to parse id")
		}
		key = c.db.IDKey(c.ctx, Kind, id)
		err = c.db.Get(c.ctx, key, exe)
		if err != nil {
			return nil, errDatastore.WithError(err).Annotate("failed to get")
		}
	}
	return exe, nil
}

// GetAll gets all Executions ordered by created datetime.
func (c *Client) GetAll(start, limit int) ([]*Execution, error) {
	var exes []*Execution
	_, err := c.db.QueryOrdered(c.ctx, Kind, start, limit, "-Date", &exes)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryOrdered")
	}
	return exes, nil
}

// Update updates an Execution.
func (c *Client) Update(exe *Execution) (*Execution, error) {
	if exe.CreatedDatetime.IsZero() {
		exe.CreatedDatetime = time.Now()
	}
	key := c.db.IDKey(c.ctx, Kind, exe.ID)
	key, err := c.db.Put(c.ctx, key, exe)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to put")
	}
	if exe.ID == 0 {
		// handle new entry
		exe.ID = key.IntID()
		_, err = c.db.Put(c.ctx, key, exe)
		if err != nil {
			return nil, errDatastore.WithError(err).Annotate("failed to put")
		}
	}
	return exe, nil
}
