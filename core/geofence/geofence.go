package geofence

// ServizeZone geofence should look like:
// type Geofence struct {
// 	ID           : common.location.String(),
// 	Name         : common.location.String(),
// 	Type         : geofence.ServiceZone,
// 	Points       : []Point,
// }
//
// Driver geofence should look like:
// type Geofence struct {
// 	ID           : UserID,
// 	Name         : User Full Name,
// 	Type         : geofence.ServiceZone,
// 	Points       : []Point,
// }

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/core/maps"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
)

var (
	kind = "Geofence"
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

type Type string

var (
	DeliveryDriverNames deliveryDriverNames = deliveryDriverNames{}
)

const (
	DeliveryDriverZone Type = "DeliveryDriver"
	ServiceZone        Type = "ServiceZone"
)

type deliveryDriverNames struct{}

func (d deliveryDriverNames) Founder() string {
	return "Founder"
}

func (d deliveryDriverNames) JoyDriv() string {
	return "JoyDriv"
}

// Point is a point.
type Point struct {
	common.GeoPoint
}

// Geofence is a polygon related to a geofence.
type Geofence struct {
	ID          string  `json:"id" datastore:",noindex"`
	Name        string  `json:"name" datastore:",index"`
	Type        Type    `json:"type" datastore:",index"`
	DriverEmail string  `json:"driver_email" datastore:",index"`
	DriverName  string  `json:"driver_name" datastore:",noindex"`
	Points      []Point `json:"points" datastore:",noindex"`
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

// UpdateGeofence creates or updates a geofence zone.
func (c *Client) UpdateGeofence(ctx context.Context, fence *Geofence) error {
	if fence.ID == "" {
		return errBadRequest.Annotate("id is empty")
	}
	if string(fence.Type) == "" {
		return errBadRequest.Annotate("type no set")
	}
	polygon := NewPolygon(fence.Points)
	if !polygon.IsClosed() {
		return errBadRequest.Annotate("polygon is closed")
	}
	key := c.db.NameKey(ctx, kind, fence.ID)
	_, err := c.db.Put(ctx, key, fence)
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
	_, err := c.db.QueryFilter(ctx, kind, 0, 1, "DriverID=", driverID, fence)
	if err != nil {
		return errDatastore.WithError(err).Annotate("failed to db.QueryFilter")
	}
	return nil
}

// InServizeZone checks if an address is in Service zone.
func (c *Client) InServiceZone(addr *common.Address) (bool, error) {
	var err error
	if !addr.GeoPoint.Valid() {
		err = maps.GetGeopoint(c.ctx, addr)
		if err != nil {
			return false, errors.Annotate(err, "failed to maps.GetGeopoint")
		}
	}
	var zones []*Geofence
	_, err = c.db.QueryFilter(c.ctx, kind, 0, 100, "Type=", ServiceZone, &zones)
	if err != nil {
		return false, errDatastore.WithError(err).Annotate("failed to db.QueryFilter")
	}
	for _, zone := range zones {
		polygon := NewPolygon(zone.Points)
		contains := polygon.Contains(Point{GeoPoint: addr.GeoPoint})
		if contains {
			return true, nil
		}
	}
	return false, nil
}
