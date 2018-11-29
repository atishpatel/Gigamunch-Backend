package execution

import (
	"context"
	"fmt"
	"sort"
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
			return nil, errBadRequest.WithError(err).Annotate("failed to parse id")
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
	_, err := c.db.QueryOrdered(c.ctx, Kind, start, limit, "-CreatedDatetime", &exes)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryOrdered")
	}
	sort.Slice(exes, func(i, j int) bool {
		if exes[i].Date == "" {
			return true
		}
		if exes[j].Date == "" {
			return false
		}
		return exes[i].Date > exes[j].Date
	})
	return exes, nil
}

// GetAfterDate gets all the fulture executions after date.
func (c *Client) GetAfterDate(t time.Time) ([]*Execution, error) {
	var exes []*Execution
	_, err := c.db.QueryFilterOrdered(c.ctx, Kind, 0, 1000, "-Date", "Date>=", t.Format(DateFormat), &exes)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryFilterOrdered")
	}
	return exes, nil
}

// GetBeforeDate gets all the fulture executions eefore date.
func (c *Client) GetBeforeDate(t time.Time) ([]*Execution, error) {
	var exes []*Execution
	_, err := c.db.QueryFilterOrdered(c.ctx, Kind, 0, 1000, "-Date", "Date<=", t.Format(DateFormat), &exes)
	if err != nil {
		return nil, errDatastore.WithError(err).Annotate("failed to QueryFilterOrdered")
	}
	return exes, nil
}

// Update updates an Execution.
func (c *Client) Update(exe *Execution) (*Execution, error) {
	var err error
	if exe.CreatedDatetime.IsZero() {
		exe.CreatedDatetime = time.Now()
	}
	key := c.db.IDKey(c.ctx, Kind, exe.ID)
	var exeOld *Execution
	if exe.ID != 0 {
		exeOld = &Execution{}
		err = c.db.Get(c.ctx, key, exeOld)
		if err != nil {
			return nil, errDatastore.WithError(err).Annotate("failed to get")
		}
	}
	key, err = c.db.Put(c.ctx, key, exe)
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
	if exeOld != nil {
		c.log.ExecutionUpdate(exe.ID, exeOld, exe)
	}
	return exe, nil
}
