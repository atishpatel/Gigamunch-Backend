package geofence

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

var (
	standAppEngine bool
	db             common.DB
	kind           = "Geofence"
)

var (
	errDatastore  = errors.InternalServerError
	errInternal   = errors.InternalServerError
	errBadRequest = errors.BadRequestError
)

// Client is a client for manipulating subscribers.
type Client struct {
	ctx context.Context
	log *logging.Client
}

type Type string

var (
	JoyDriv       Type = "JoyDriv"
	FounderDriver Type = "FounderDriver"
	ServiceZone   Type = "ServiceZone"
)

// Point is a point.
type Point struct {
	common.GeoPoint
}

// Geofence is a polygon related to a geofence.
type Geofence struct {
	ID          string  `json:"id" datastore:",noindex"`
	Name        string  `json:"name" datastore:",index"`
	Type        Type    `json:"type" datastore:",index"`
	DriverID    int64   `json:"driver_id" datastore:",index"`
	DriverEmail string  `json:"driver_email" datastore:",index"`
	DriverName  string  `json:"driver_name" datastore:",noindex"`
	Points      []Point `json:"points" datastore:",noindex"`
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

// AddGeofence adds a geofence zone.
func (c *Client) AddGeofence(ctx context.Context, fence *Geofence) error {
	if fence.ID == "" && fence.DriverID == 0 {
		return errBadRequest.Annotate("id is empty")
	}
	polygon := NewPolygon(fence.Points)
	if !polygon.IsClosed() {
		return errBadRequest.Annotate("polygon is closed")
	}
	var key common.Key
	if fence.ID == "" {
		key = db.NameKey(ctx, kind, fence.ID)
	} else {
		key = db.IDKey(ctx, kind, fence.DriverID)
	}
	_, err := db.Put(ctx, key, fence)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to db.Put")
	}
	return nil
}

// GetDriverZone gets a driver's zone.
func (c *Client) GetDriverZone(ctx context.Context, driverID int64) error {
	if driverID == 0 {
		return errBadRequest.Annotate("invalid driverID")
	}
	fence := new(Geofence)
	_, err := db.QueryFilter(ctx, kind, 0, 1, "DriverID=", driverID, fence)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to db.QueryFilter")
	}
	return nil
}

// InNashvilleZone checks if an address is in Nashville zone.
func (c *Client) InNashvilleZone(ctx context.Context, addr *common.Address) (bool, error) {
	var err error
	if !addr.GeoPoint.Valid() {
		// TODO get geopoint form address
		err = maps.GetGeopoint(ctx, addr)
		if err != nil {
			return false, errors.Annotate(err, "failed to maps.GetGeopoint")
		}
	}
	fence := new(Geofence)
	key := db.NameKey(ctx, kind, common.Nashville.String())
	err = db.Get(ctx, key, fence)
	if err != nil {
		return false, errDatastore.WithError(err).Annotate("failed to db.Get")
	}
	polygon := NewPolygon(fence.Points)
	contains := polygon.Contains(Point{GeoPoint: addr.GeoPoint})
	return contains, nil
}

// Setup sets up the logging package.
func Setup(ctx context.Context, standardAppEngine bool, dbC common.DB) error {
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
	return nil
}

func setup(ctx context.Context) error {
	return nil
}
