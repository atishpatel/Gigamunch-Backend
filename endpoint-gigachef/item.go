package gigachef

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/item"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

// Item is an item created by the chef that holds the basic info for a post
type Item struct {
	BaseItem                  // embedded
	ID                int     `json:"id"`
	LastUsedDateTime  int     `json:"last_used_datetime"`
	NumPostsCreated   int     `json:"num_posts_created"`
	NumTotalOrders    int     `json:"num_total_orders"`
	AverageItemRating float32 `json:"average_item_rating"`
	NumRatings        int     `json:"num_ratings"`
}

// Set takes a item form the item package and converts it to a endpoint item
func (i *Item) Set(id int, item *item.Item) {
	i.ID = id
	i.Title = item.Title
	i.Subtitle = item.Subtitle
	i.Description = item.Description
	i.Ingredients = item.Ingredients
	i.GeneralTags = item.GeneralTags
	i.DietaryNeedsTags = item.DietaryNeedsTags
	i.Photos = item.Photos
	i.LastUsedDateTime = int(item.LastUsedDateTime.Unix())
}

// Get creates a item.Item version of the endpoint item
func (i *Item) Get() *item.Item {
	item := &item.Item{}
	item.Title = i.Title
	item.Subtitle = i.Subtitle
	item.Description = i.Description
	item.Ingredients = i.Ingredients
	item.GeneralTags = i.GeneralTags
	item.DietaryNeedsTags = i.DietaryNeedsTags
	item.Photos = i.Photos
	item.LastUsedDateTime = time.Unix(int64(i.LastUsedDateTime), 0)
	item.NumPostsCreated = i.NumPostsCreated
	item.NumTotalOrders = i.NumTotalOrders
	item.AverageItemRating = i.AverageItemRating
	item.NumRatings = i.NumRatings
	return item
}

// GetItemReq is the input request needed for GetItem.
type GetItemReq struct {
	GigaToken string `json:"gigatoken"`
	ID        int    `json:"id"`
}

// Gigatoken returns the GigaToken string
func (req *GetItemReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *GetItemReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty")
	}
	if req.ID == 0 {
		return fmt.Errorf("ID is empty")
	}
	return nil
}

// GetItemResp is the output response for GetItem.
type GetItemResp struct {
	Item Item                 `json:"item"`
	Err  errors.ErrorWithCode `json:"err,omitempty"`
}

// GetItem is an endpoint that gets an item.
func (service *Service) GetItem(ctx context.Context, req *GetItemReq) (*GetItemResp, error) {
	resp := new(GetItemResp)
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "GetItem err: ", resp.Err)
		}
	}()
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	i, err := item.GetItem(ctx, user, int64(req.ID))
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Item.Set(req.ID, i)
	return resp, nil
}

// GetItemsReq is the input request needed for GetItems.
type GetItemsReq struct {
	GigaToken  string `json:"gigatoken"`
	StartLimit int    `json:"start_limit"`
	EndLimit   int    `json:"end_limit"`
}

// Gigatoken returns the GigaToken string
func (req *GetItemsReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *GetItemsReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
	}
	if req.StartLimit < 0 || req.EndLimit < 0 {
		return fmt.Errorf("Limit is out of range.")
	}
	if req.EndLimit <= req.StartLimit {
		return fmt.Errorf("EndLimit cannot be less than or equal to StartLimit.")
	}
	return nil
}

// GetItemsResp is the output response for GetItems.
type GetItemsResp struct {
	Items []Item               `json:"items"`
	Err   errors.ErrorWithCode `json:"err,omitempty"`
}

// GetItems is an endpoint that gets a Gigachef's items.
func (service *Service) GetItems(ctx context.Context, req *GetItemsReq) (*GetItemsResp, error) {
	resp := new(GetItemsResp)
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "GetItems err: ", resp.Err)
		}
	}()
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	limit := &types.Limit{Start: req.StartLimit, End: req.EndLimit}
	ids, items, err := item.GetItems(ctx, user, limit)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Items = make([]Item, len(ids))
	for i := range ids {
		resp.Items[i].Set(int(ids[i]), &items[i])
	}
	return resp, nil
}

// SaveItemReq is the input request needed for SaveItem.
type SaveItemReq struct {
	GigaToken string `json:"gigatoken"`
	Item      Item   `json:"item"`
}

// Gigatoken returns the GigaToken string
func (req *SaveItemReq) Gigatoken() string {
	return req.GigaToken
}

// Valid validates a req
func (req *SaveItemReq) Valid() error {
	if req.GigaToken == "" {
		return fmt.Errorf("GigaToken is empty.")
	}
	return nil
}

// SaveItemResp is the output respones for a SaveItem.
type SaveItemResp struct {
	Item Item                 `json:"item"`
	Err  errors.ErrorWithCode `json:"err,omitempty"`
}

// SaveItem is an endpoint that saves an item form a Gigachef.
// If item id is 0, a new item is created.
func (service *Service) SaveItem(ctx context.Context, req *SaveItemReq) (*SaveItemResp, error) {
	resp := new(SaveItemResp)
	defer func() {
		if resp.Err.Code != 0 && resp.Err.Code != errors.CodeInvalidParameter {
			utils.Errorf(ctx, "SaveItem err: ", resp.Err)
		}
	}()
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	itemID, err := item.SaveItem(ctx, user, int64(req.Item.ID), req.Item.Get())
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Item = req.Item
	resp.Item.ID = int(itemID)
	return resp, nil
}
