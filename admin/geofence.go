package main

import (
	"context"
	"net/http"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/geofence"

	pb "github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbadmin"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// TODO: implement to take in request
// UpdateGeofence updates a geofence.
func (s *server) UpdateGeofence(ctx context.Context, w http.ResponseWriter, r *http.Request, log *logging.Client) Response {
	var err error
	req := new(pb.GetLogsReq)

	// decode request
	err = decodeRequest(ctx, r, req)
	if err != nil {
		return failedToDecode(err)
	}
	// end decode request
	geofenceC, err := geofence.NewClient(ctx, log, s.db, s.sqlDB, s.serverInfo)
	if err != nil {
		return errors.Annotate(err, "failed to geofence.NewClient")
	}
	points := []common.GeoPoint{
		common.GeoPoint{Longitude: -86.6292572, Latitude: 36.2509183},
		common.GeoPoint{Longitude: -86.6388705, Latitude: 36.2777685},
		common.GeoPoint{Longitude: -86.6539764, Latitude: 36.2885632},
		common.GeoPoint{Longitude: -86.7157746, Latitude: 36.2702966},
		common.GeoPoint{Longitude: -86.7576599, Latitude: 36.2719574},
		common.GeoPoint{Longitude: -86.8139649, Latitude: 36.2459345},
		common.GeoPoint{Longitude: -86.8805695, Latitude: 36.2304274},
		common.GeoPoint{Longitude: -86.9011731, Latitude: 36.1971752},
		common.GeoPoint{Longitude: -86.8956756, Latitude: 36.1417563},
		common.GeoPoint{Longitude: -86.9739841, Latitude: 36.1215297},
		common.GeoPoint{Longitude: -87.0350646, Latitude: 36.0701923},
		common.GeoPoint{Longitude: -86.9780732, Latitude: 35.9552207},
		common.GeoPoint{Longitude: -86.9408502, Latitude: 35.8638018},
		common.GeoPoint{Longitude: -86.8441772, Latitude: 35.8428646},
		common.GeoPoint{Longitude: -86.7350006, Latitude: 35.8907189},
		common.GeoPoint{Longitude: -86.6629813, Latitude: 35.9431917},
		common.GeoPoint{Longitude: -86.5483139, Latitude: 36.0053122},
		common.GeoPoint{Longitude: -86.5626526, Latitude: 36.0835115},
		common.GeoPoint{Longitude: -86.5956116, Latitude: 36.1323292},
		common.GeoPoint{Longitude: -86.5647125, Latitude: 36.1927545},
		common.GeoPoint{Longitude: -86.6004181, Latitude: 36.2498108},
		common.GeoPoint{Longitude: -86.6292572, Latitude: 36.2509183},
	}

	fence := &geofence.Geofence{
		ID:   common.Nashville.ID(),
		Name: common.Nashville.String(),
		Type: geofence.ServiceZone,
	}
	for _, p := range points {
		fence.Points = append(fence.Points, geofence.Point{GeoPoint: p})
	}
	err = geofenceC.UpdateGeofence(ctx, fence)
	if err != nil {
		return errors.Annotate(err, "failed to geofence.UpdateGeofence")
	}
	resp := &pb.ErrorOnlyResp{}
	return resp
}
