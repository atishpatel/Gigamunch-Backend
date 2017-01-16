package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/image"

	"google.golang.org/appengine"

	"strconv"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/like"
	"github.com/atishpatel/Gigamunch-Backend/corenew/maps"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/review"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	bucketName          string
	errInternal         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error while uploading file."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "An invalid parameter was used."}
	projectID           string
)

func handleResp(ctx context.Context, w http.ResponseWriter, funcName string, respErr errors.ErrorWithCode, resp interface{}) {
	// encode json resp and log errors
	if projectID == "" {
		projectID = os.Getenv("PROJECTID")
	}
	if respErr.Code != 0 && respErr.Code != errors.CodeInvalidParameter {
		utils.Criticalf(ctx, "%s SERVER: Error %s: %+v", projectID, funcName, respErr)
		w.WriteHeader(http.StatusInternalServerError)
	} else if respErr.Code != 0 {
		utils.Infof(ctx, "CodeInvalidParameter(%d) Error %s: %+v", respErr.Code, funcName, respErr)
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		utils.Criticalf(ctx, "%s SERVER: Error encoding json: %+v", projectID, err)
	}
}

type urlResp struct {
	URL string               `json:"url"`
	Err errors.ErrorWithCode `json:"err"`
}

func handleUpload(w http.ResponseWriter, req *http.Request) {
	resp := new(urlResp)
	ctx := appengine.NewContext(req)

	defer handleResp(ctx, w, "Upload", resp.Err, resp)

	time.Sleep(500 * time.Millisecond) // check if failed to find blob file bug is fixed with this

	// get file
	blobs, _, err := blobstore.ParseUpload(req)
	if err != nil {
		resp.Err = errInvalidParameter.WithMessage("Error parsing multipart form.").WithError(err)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		resp.Err = errInvalidParameter.WithMessage("No file was uploaded.")
		return
	}
	opts := &image.ServingURLOptions{
		Secure: true,
		Crop:   true,
	}
	time.Sleep(500 * time.Millisecond) // check if failed to find blob file bug is fixed with this
	// ctx, _ = context.WithDeadline(ctx, time.Now().Add(60*time.Second))
	url, err := image.ServingURL(ctx, file[0].BlobKey, opts)
	if err != nil {
		deadline, _ := ctx.Deadline()
		resp.Err = errInternal.WithError(err).Wrapf("failed to get image.ServingURL (blobkey: %v) (now:%v context.Deadline:%v)", file[0].BlobKey, time.Now(), deadline)

		return
	}
	resp.URL = url.String()
}

func handleGetUploadURL(w http.ResponseWriter, req *http.Request) {
	resp := new(urlResp)
	ctx := appengine.NewContext(req)
	defer handleResp(ctx, w, "GetUploadURL", resp.Err, resp)
	if bucketName == "" {
		bucketName = config.GetBucketName(ctx)
	}
	// get user
	user, err := getUserFromCookie(req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return
	}
	opts := &blobstore.UploadURLOptions{
		StorageBucket: fmt.Sprintf("%s/%s", bucketName, user.ID),
	}
	uploadURL, err := blobstore.UploadURL(ctx, "/upload", opts)
	if err != nil {
		resp.Err = errInternal.WithError(err).Wrap("error getting blobstore.UploadURL")
		return
	}
	resp.URL = uploadURL.String()
}

type getFeedResp struct {
	Menus []MenuWithItems      `json:"menus"`
	Err   errors.ErrorWithCode `json:"err"`
}

func handleGetFeed(w http.ResponseWriter, req *http.Request) {
	resp := new(getFeedResp)
	ctx := appengine.NewContext(req)

	defer handleResp(ctx, w, "GetFeed", resp.Err, resp)

	var startIndex int32
	var endIndex int32 = 1000
	lat := 36.1627
	long := -86.7816

	itemC := item.New(ctx)
	itemIDs, menuIDs, cookIDs, err := itemC.GetActiveItemIDs(startIndex, endIndex, lat, long)
	if err != nil {
		resp.Err = errors.Wrap("failed to item.GetActiveItemIDs", err)
		return
	}
	// get items
	var items []item.Item
	itemsErrChan := make(chan error, 1)
	go func() {
		var goErr error
		items, goErr = itemC.GetMulti(itemIDs)
		itemsErrChan <- goErr
	}()
	// get menus
	var menus map[int64]*menu.Menu
	menusErrChan := make(chan error, 1)
	go func() {
		var goErr error
		menuC := menu.New(ctx)
		menus, goErr = menuC.GetMulti(menuIDs)
		menusErrChan <- goErr
	}()
	// get cooks
	var cooks map[string]*cook.Cook
	cooksErrChan := make(chan error, 1)
	go func() {
		var goErr error
		cookC := cook.New(ctx)
		cooks, goErr = cookC.GetMulti(cookIDs)
		cooksErrChan <- goErr
	}()
	// get likes
	var numLikes []int32
	likeErrChan := make(chan error, 1)
	go func() {
		// get user if there
		var userID string
		// if req.Gigatoken != "" {
		// 	user, _ := auth.GetUserFromToken(ctx, req.Gigatoken)
		// 	if user != nil {
		// 		userID = user.ID
		// 	}
		// }
		var goErr error
		likeC := like.New(ctx)
		_, numLikes, goErr = likeC.LikesItems(userID, itemIDs)
		likeErrChan <- goErr
	}()
	// handle errors
	err = processErrorChans(itemsErrChan, menusErrChan, cooksErrChan, likeErrChan)
	if err != nil {
		resp.Err = errors.Wrap("failed to item.GetMulti or menu.GetMulti or cook.GetMulti", err)
		return
	}
	// get menu order
	menuOrder := make([]int64, len(menus))
	index := 0
	for i := range items {
		found := false
		for _, v := range menuOrder {
			if v == items[i].MenuID {
				found = true
			}
		}
		if !found {
			menuOrder[index] = items[i].MenuID
			index++
		}
	}
	// set menus
	eaterPoint := types.GeoPoint{Latitude: lat, Longitude: long}
	resp.Menus = make([]MenuWithItems, len(menus))
	for i, v := range menuOrder {
		c := cooks[menus[v].CookID] // cook for this menu
		// get exchangeoptions and distance
		ems := types.GetExchangeMethods(c.Address.GeoPoint, c.DeliveryRange, c.DeliveryPrice, eaterPoint)
		distance := eaterPoint.GreatCircleDistance(c.Address.GeoPoint)
		resp.Menus[i].Menu = *menus[v]
		resp.Menus[i].Cook.Cook = *c
		resp.Menus[i].Cook.Distance = distance
		resp.Menus[i].Cook.ExchangeOptions = make([]ExchangeOption, len(ems))
		for j := range ems {
			resp.Menus[i].Cook.ExchangeOptions[j].ID = ems[j].ID()
			resp.Menus[i].Cook.ExchangeOptions[j].Price = ems[j].Price
			resp.Menus[i].Cook.ExchangeOptions[j].IsDelivery = ems[j].Delivery()
		}
		menuID := menus[v].ID
		for j := range items {
			if items[j].MenuID == menuID {
				resp.Menus[i].Items = append(resp.Menus[i].Items, Item{
					Item:            items[j],
					NumLikes:        numLikes[j],
					PricePerServing: payment.GetPricePerServing(items[j].CookPricePerServing),
				})
			}
		}
	}
}

type getItemResp struct {
	Item    Item                 `json:"item"`
	Reviews []Review             `json:"reviews"`
	Cook    Cook                 `json:"cook"`
	Err     errors.ErrorWithCode `json:"err"`
}

func handleGetItem(w http.ResponseWriter, req *http.Request) {
	resp := new(getItemResp)
	ctx := appengine.NewContext(req)

	defer handleResp(ctx, w, "GetItem", resp.Err, resp)

	lat := 36.1627
	long := -86.7816
	var itemID int64
	itemIDString := req.URL.Query().Get("id")
	if itemIDString == "" {
		resp.Err = errors.ErrorWithCode{
			Code:    errors.CodeInvalidParameter,
			Message: "Invalid item id.",
			Detail:  fmt.Sprintf("Item ID(%s) is not a valid parameter", itemIDString),
		}
		return
	}
	itemID, err := strconv.ParseInt(itemIDString, 10, 64)
	if err != nil {
		resp.Err = errors.Wrap("failed to strconv.ParseInt", err)
		return
	}

	likeC := like.New(ctx)
	// get item
	itemC := item.New(ctx)
	item, err := itemC.Get(itemID)
	if err != nil {
		resp.Err = errors.Wrap("failed to item.Get", err)
		return
	}
	// get cook
	var c *cook.Cook
	cooksErrChan := make(chan error, 1)
	go func() {
		var goErr error
		cookC := cook.New(ctx)
		c, goErr = cookC.Get(item.CookID)
		cooksErrChan <- goErr
	}()
	// get reviews
	var reviews []*review.Review
	reviewErrChan := make(chan error, 1)
	go func() {
		var goErr error
		reviewC := review.New(ctx)
		reviews, goErr = reviewC.GetByCookID(item.CookID, item.ID, 0, 5)
		reviewErrChan <- goErr
	}()
	// get likes
	var numLikes []int32
	likeErrChan := make(chan error, 1)
	go func() {
		// get user if there
		var userID string
		var goErr error
		_, numLikes, goErr = likeC.LikesItems(userID, []int64{item.ID})
		likeErrChan <- goErr
	}()
	// get cook likes
	var cookLikes int32
	cookLikeErrChan := make(chan error, 1)
	go func() {
		var goErr error
		cookLikes, goErr = likeC.GetNumCookLikes(item.CookID)
		cookLikeErrChan <- goErr
	}()
	// handle errors
	err = processErrorChans(cooksErrChan, reviewErrChan, likeErrChan, cookLikeErrChan)
	if err != nil {
		resp.Err = errors.Wrap("failed to cook.Get or review.GetByCookID or like.LikesItems or like.GetNumCookLikes", err)
		return
	}
	eaterPoint := types.GeoPoint{Latitude: lat, Longitude: long}
	cookPoint := c.Address.GeoPoint
	// get distance
	distance, _, err := maps.GetDistance(ctx, cookPoint, eaterPoint)
	if err != nil {
		resp.Err = errors.Wrap("failed to maps.GetDistance", err)
		return
	}
	// get exchangeoptions
	ems := types.GetExchangeMethods(cookPoint, c.DeliveryRange, c.DeliveryPrice, eaterPoint)
	resp.Cook.ExchangeOptions = make([]ExchangeOption, len(ems))
	for i := range ems {
		resp.Cook.ExchangeOptions[i].ID = ems[i].ID()
		resp.Cook.ExchangeOptions[i].Price = ems[i].Price
		resp.Cook.ExchangeOptions[i].IsDelivery = ems[i].Delivery()
	}
	resp.Item.Item = *item
	resp.Item.NumLikes = numLikes[0]
	resp.Item.PricePerServing = payment.GetPricePerServing(item.CookPricePerServing)
	resp.Cook.Cook = *c
	resp.Cook.Distance = distance
	resp.Cook.Likes = cookLikes
	resp.Reviews = make([]Review, len(reviews))
	for i := range reviews {
		resp.Reviews[i].Review = *reviews[i]
	}
	return
}

type ExchangeOption struct {
	ID         int64   `json:"id,omitempty"`
	IsDelivery bool    `json:"is_delivery,omitempty"`
	Price      float32 `json:"price,omitempty"`
}

type Cook struct {
	cook.Cook
	Distance        float32          `json:"distance"`
	ExchangeOptions []ExchangeOption `json:"exchange_options"`
	Likes           int32            `json:"likes"`
}

type Review struct {
	review.Review
}

// Item is an Item.
type Item struct {
	NumLikes int32 `json:"num_likes"`
	item.Item
	PricePerServing float32 `json:"price_per_serving"`
}

// MenuWithItems is a Menu with all it's Items.
type MenuWithItems struct {
	menu.Menu
	Items []Item `json:"items"`
	Cook  Cook   `json:"cook"` // TODO change to basic cook
}

// processErrorChans returns an error if any of the error channels return an error
func processErrorChans(errs ...<-chan error) error {
	var err error
	for _, v := range errs {
		err = <-v
		if err != nil {
			return err
		}
	}
	return nil
}
