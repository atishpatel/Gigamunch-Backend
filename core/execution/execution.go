package execution

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	standAppEngine bool
	db             common.DB
	projID         string
)

var (
	errDatastore = errors.InternalServerError
	errInternal  = errors.InternalServerError
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx context.Context
	log *logging.Client
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client) (*Client, error) {
	var err error
	if standAppEngine {
		err = setup(ctx)
		if err != nil {
			return nil, err
		}
	}
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	return &Client{
		ctx: ctx,
		log: log,
	}, nil
}

// Get gets a culture Execution.
func (c *Client) Get(id int64) (*Execution, error) {
	key := db.IDKey(c.ctx, Kind, id)
	exe := new(Execution)
	err := db.Get(c.ctx, key, exe)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to get")
	}
	return exe, nil
}

// GetAll gets all Executions ordered by created datetime.
func (c *Client) GetAll(start, limit int) ([]*Execution, error) {
	var exes []*Execution
	_, err := db.QueryOrdered(c.ctx, Kind, start, limit, "-CreatedDatetime", exes)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryOrdered")
	}
	return exes, nil
}

// Update updates an Execution.
func (c *Client) Update(exe *Execution) (*Execution, error) {
	key := db.IDKey(c.ctx, Kind, exe.ID)
	key, err := db.Put(c.ctx, key, exe)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to put")
	}
	if exe.ID == 0 {
		exe.ID = key.IntID()
		_, err = db.Put(c.ctx, key, exe)
		if err != nil {
			return nil, errDatastore.WithError(err).Annotate("failed to put")
		}
	}
	return exe, nil
}

// Setup sets up the logging package.
func Setup(ctx context.Context, standardAppEngine bool, projectID string, dbC common.DB) error {
	var err error
	standAppEngine = standardAppEngine
	if !standAppEngine {
		err = setup(ctx)
		if err != nil {
			return err
		}
	}
	if dbC == nil {
		return fmt.Errorf("db cannot be nil for sub")
	}
	db = dbC
	projID = projectID
	return nil
}

func setup(ctx context.Context) error {
	return nil
}
