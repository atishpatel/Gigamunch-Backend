package gigachef

import (
	"fmt"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/item"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// Item is an item created by the chef that holds the basic info for a post
type Item struct {
	BaseItem                  // embedded
	ID                string  `json:"id"`
	ID64              int64   `json:"-"`
	LastUsedDateTime  int     `json:"last_used_datetime"`
	NumPostsCreated   int     `json:"num_posts_created"`
	NumTotalOrders    int     `json:"num_total_orders"`
	AverageItemRating float32 `json:"average_item_rating"`
	NumRatings        int     `json:"num_ratings"`
}

// Set takes a item form the item package and converts it to a endpoint item
func (i *Item) Set(id int64, item *item.Item) {
	i.ID = itos(id)
	i.ID64 = id
	i.Title = item.Title
	i.Description = item.Description
	i.Ingredients = item.Ingredients
	i.GeneralTags = item.GeneralTags
	i.DietaryNeedsTags = item.DietaryNeedsTags
	i.Photos = item.Photos
	i.LastUsedDateTime = ttoi(item.LastUsedDateTime)
}

// Get creates a item.Item version of the endpoint item
func (i *Item) Get() *item.Item {
	item := &item.Item{}
	item.Title = i.Title
	item.Description = i.Description
	item.Ingredients = i.Ingredients
	item.GeneralTags = i.GeneralTags
	item.DietaryNeedsTags = i.DietaryNeedsTags
	item.Photos = i.Photos
	item.LastUsedDateTime = itot(i.LastUsedDateTime)
	item.NumPostsCreated = i.NumPostsCreated
	item.NumTotalOrders = i.NumTotalOrders
	item.AverageItemRating = i.AverageItemRating
	item.NumRatings = i.NumRatings
	return item
}

// GetItemReq is the input request needed for GetItem.
type GetItemReq struct {
	Gigatoken string `json:"gigatoken"`
	ID        string `json:"id"`
	ID64      int64  `json:"-"`
}

// gigatoken returns the Gigatoken string
func (req *GetItemReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *GetItemReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty")
	}
	if req.ID == "" || req.ID == "0" {
		return fmt.Errorf("ID is empty")
	}
	var err error
	req.ID64, err = stoi(req.ID)
	if err != nil {
		return fmt.Errorf("Error with ID: %v", err)
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
	defer handleResp(ctx, "GetItem", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	i, err := item.GetItem(ctx, user, req.ID64)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Item.Set(req.ID64, i)
	return resp, nil
}

// GetItemsReq is the input request needed for GetItems.
type GetItemsReq struct {
	Gigatoken  string `json:"gigatoken"`
	StartLimit int    `json:"start_limit"`
	EndLimit   int    `json:"end_limit"`
}

// gigatoken returns the Gigatoken string
func (req *GetItemsReq) gigatoken() string {
	return req.Gigatoken
}

// Valid validates a req
func (req *GetItemsReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
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
	defer handleResp(ctx, "GetItems", resp.Err)
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
		resp.Items[i].Set(ids[i], &items[i])
	}
	return resp, nil
}

// SaveItemReq is the input request needed for SaveItem.
type SaveItemReq struct {
	Gigatoken string `json:"gigatoken"`
	Item      Item   `json:"item"`
}

// gigatoken returns the Gigatoken string
func (req *SaveItemReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *SaveItemReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	if req.Item.ID != "" {
		var err error
		req.Item.ID64, err = stoi(req.Item.ID)
		if err != nil {
			return fmt.Errorf("Error with ID: %v", err)
		}
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
	defer handleResp(ctx, "SaveItem", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	itemID, err := item.SaveItem(ctx, user, req.Item.ID64, req.Item.Get())
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	resp.Item = req.Item
	resp.Item.ID = itos(itemID)
	return resp, nil
}
