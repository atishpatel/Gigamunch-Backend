package gigamuncher

import (
	"encoding/json"
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/like"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"golang.org/x/net/context"
)

// LikeReq is the request needed to like or unlike items
type LikeReq struct {
	Gigatoken string      `json:"gigatoken"`
	ItemID    json.Number `json:"item_id"`
	ItemID64  int64       `json:"-"`
}

func (req *LikeReq) gigatoken() string {
	return req.Gigatoken
}

func (req *LikeReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	var err error
	req.ItemID64, err = req.ItemID.Int64()
	if err != nil {
		return err
	}
	if req.ItemID64 == 0 {
		return fmt.Errorf("ItemID is 0.")
	}
	return nil
}

// LikeItem is used to like an item
func (service *Service) LikeItem(ctx context.Context, req *LikeReq) (*ErrorOnlyResp, error) {
	likeC := like.New(ctx)
	return likeOperation(ctx, req, "Like", likeC.Like)
}

// UnlikeItem is used to unlike an item
func (service *Service) UnlikeItem(ctx context.Context, req *LikeReq) (*ErrorOnlyResp, error) {
	likeC := like.New(ctx)
	return likeOperation(ctx, req, "Unlike", likeC.Unlike)
}

func likeOperation(ctx context.Context, req *LikeReq, opName string, fn func(string, int64) error) (*ErrorOnlyResp, error) {
	resp := new(ErrorOnlyResp)
	defer handleResp(ctx, opName, resp.Err)

	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	err = fn(user.ID, req.ItemID64)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("failed to %s item", opName)
		return resp, nil
	}
	return resp, nil
}
