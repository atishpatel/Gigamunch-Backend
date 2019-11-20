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
		common.GeoPoint{Longitude: -86.5674591, Latitude: 36.272511},
		common.GeoPoint{Longitude: -86.5564718, Latitude: 36.3194122},
		common.GeoPoint{Longitude: -86.5283206, Latitude: 36.2871782},
		common.GeoPoint{Longitude: -86.4239502, Latitude: 36.2863493},
		common.GeoPoint{Longitude: -86.389567, Latitude: 36.3114654},
		common.GeoPoint{Longitude: -86.3970799, Latitude: 36.34098},
		common.GeoPoint{Longitude: -86.3887165, Latitude: 36.4066105},
		common.GeoPoint{Longitude: -86.4406676, Latitude: 36.4151519},
		common.GeoPoint{Longitude: -86.5004994, Latitude: 36.531158},
		common.GeoPoint{Longitude: -86.6299439, Latitude: 36.5664606},
		common.GeoPoint{Longitude: -86.727698, Latitude: 36.4144326},
		common.GeoPoint{Longitude: -86.7319105, Latitude: 36.3264684},
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
		common.GeoPoint{Longitude: -86.5674591, Latitude: 36.272511},
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
