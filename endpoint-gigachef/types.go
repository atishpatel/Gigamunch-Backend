package gigachef

import (
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
)

// BaseItem is the basic stuff in an Item
type BaseItem struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Ingredients      []string `json:"ingredients"`
	GeneralTags      []string `json:"general_tags"`
	CuisineTags      []string `json:"cuisine_tags"`
	DietaryNeedsTags []string `json:"dietary_needs_tags"`
	Photos           []string `json:"photos"`
}

// GigatokenOnlyReq is a request with only a gigatoken input
type GigatokenOnlyReq struct {
	Gigatoken string `json:"gigatoken"`
}

func (req *GigatokenOnlyReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *GigatokenOnlyReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// ErrorOnlyResp is a response with only an error with code
type ErrorOnlyResp struct {
	Err errors.ErrorWithCode `json:"err"`
}

// Address represents a location with GeoPoints and address
type Address struct {
	APT       string `json:"apt"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	Country   string `json:"country"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (a *Address) get() (*types.Address, error) {
	var lat, long float64
	var err error
	if a.Latitude != "" {
		lat, err = stof64(a.Latitude)
		if err != nil {
			return nil, err
		}
	}
	if a.Longitude != "" {
		long, err = stof64(a.Longitude)
		if err != nil {
			return nil, err
		}
	}
	add := &types.Address{
		APT:     a.APT,
		Street:  a.Street,
		City:    a.City,
		State:   a.State,
		Zip:     a.Zip,
		Country: a.Country,
		GeoPoint: types.GeoPoint{
			Latitude:  lat,
			Longitude: long,
		},
	}
	return add, nil
}

func (a *Address) set(add *types.Address) {
	a.APT = add.APT
	a.Street = add.Street
	a.City = add.City
	a.State = add.State
	a.Zip = add.Zip
	a.Country = add.Country
	a.Latitude = ftos64(add.Latitude)
	a.Longitude = ftos64(add.Longitude)
}
