package gigamuncher

import (
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigamuncher"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"

	"golang.org/x/net/context"
)

// GetAddressesReq is the req for GetAddresses
type GetAddressesReq struct {
	Gigatoken string `json:"gigatoken"`
}

func (req *GetAddressesReq) gigatoken() string {
	return req.Gigatoken
}

func (req *GetAddressesReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// GetAddressesResp is the resp for GetAddresses
type GetAddressesResp struct {
	Addresses []Address            `json:"addresses,omitempty"`
	Err       errors.ErrorWithCode `json:"err"`
}

// GetAddresses is used to like an item
func (service *Service) GetAddresses(ctx context.Context, req *GetAddressesReq) (*GetAddressesResp, error) {
	resp := new(GetAddressesResp)
	defer handleResp(ctx, "GetAddresses", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	muncerC := gigamuncher.New(ctx)
	addresses, err := muncerC.GetAddresses(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to gigamuncher.GetAddresses")
		return resp, nil
	}
	resp.Addresses = make([]Address, len(addresses))
	for i := range addresses {
		resp.Addresses[i].set(&addresses[i].Address, &addresses[i].AddedDateTime, addresses[i].Selected)
	}
	return resp, nil
}

// SelectAddressReq is the req for SelectAddress
type SelectAddressReq struct {
	Gigatoken   string         `json:"gigatoken"`
	Address     Address        `json:"address"`
	AddressType *types.Address `json:"-"`
}

func (req *SelectAddressReq) gigatoken() string {
	return req.Gigatoken
}

func (req *SelectAddressReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	var err error
	req.AddressType, err = req.Address.get()
	if err != nil {
		return fmt.Errorf("Failed to decode address: %#v", err)
	}
	return nil
}

// SelectAddressResp is the resp for SelectAddress
type SelectAddressResp struct {
	Address gigamuncher.Addresses `json:"address"`
	Err     errors.ErrorWithCode  `json:"err"`
}

// SelectAddress is used to like an item
func (service *Service) SelectAddress(ctx context.Context, req *SelectAddressReq) (*SelectAddressResp, error) {
	resp := new(SelectAddressResp)
	defer handleResp(ctx, "SelectAddress", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	muncerC := gigamuncher.New(ctx)
	a, err := muncerC.SelectAddress(user.ID, req.AddressType)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to gigamuncher.SelectAddress")
		return resp, nil
	}
	resp.Address = *a
	return resp, nil
}
