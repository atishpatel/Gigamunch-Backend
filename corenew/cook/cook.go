package cook

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
)

// Client is a client for the cook package.
type Client struct {
	ctx context.Context
}

// New returns a new cook Client.
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// GetMultiNamesAndPhotos returns an array of names and photos for the CookIDs.
func (c *Client) GetMultiNamesAndPhotos(ids []string) ([]string, []string, error) {
	cooks, err := getMulti(c.ctx, ids)
	if err != nil {
		return nil, nil, errDatastore.WithError(err).Wrapf("failed to getMulti menu ids: %v", ids)
	}
	names := make([]string, len(cooks))
	photos := make([]string, len(cooks))
	for i := range cooks {
		names[i] = cooks[i].Name
		photos[i] = cooks[i].PhotoURL
	}
	return names, photos, nil
}

// Get gets the cook.
func (c *Client) Get(id string) (*Cook, error) {
	cook, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get cook(%s)", id)
	}
	return cook, nil
}

// GetMulti returns an array of Cooks.
func (c *Client) GetMulti(ids []string) (map[string]*Cook, error) {
	cooks := make(map[string]*Cook, len(ids))
	for _, v := range ids {
		cooks[v] = nil
	}
	ids = nil
	unDupedIDs := make([]string, len(cooks))
	index := 0
	for i := range cooks {
		unDupedIDs[index] = i
		index++
	}
	results, err := getMulti(c.ctx, unDupedIDs)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to getMulti menu ids: %v", ids)
	}
	for i, v := range unDupedIDs {
		cooks[v] = &results[i]
	}
	return cooks, nil
}

// GetAddress gets the cook.
func (c *Client) GetAddress(id string) (*types.Address, error) {
	cook, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get cook(%s)", id)
	}
	return &cook.Address, nil
}

// Update updates a cook's info.
func (c *Client) Update(user *types.User, address *types.Address, phoneNumber, bio string,
	deliveryRange int32, weekSchedule []WeekSchedule, instagramID, twitterID, snapchatID string) (*Cook, error) {
	cook, err := get(c.ctx, user.ID)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, errDatastore.WithError(err).Wrapf("failed to get cook(%s)", user.ID)
	}
	if cook.CreatedDatetime.IsZero() {
		// is new cook
		if !appengine.IsDevAppServer() {
			// notify enis
			mC := message.New(c.ctx)
			err = mC.SendSMS("6153975516", fmt.Sprintf("%s just started an application. get on that booty. email: %s", user.Name, user.Email))
			if err != nil {
				utils.Criticalf(c.ctx, "failed to notify enis about cook(%s) submitting application", user.ID)
			}
		}
		cook.CreatedDatetime = time.Now()
	}
	cook.ID = user.ID
	cook.Name = user.Name
	cook.PhotoURL = user.PhotoURL
	cook.Email = user.Email
	cook.ProviderID = user.ProviderID
	if address != nil {
		cook.Address = *address
	}
	if phoneNumber != "" {
		cook.PhoneNumber = phoneNumber
	}
	if cook.BTSubMerchantID == "" {
		tmpID := user.ID
		for len(tmpID) <= 32 {
			tmpID += tmpID
		}
		cook.BTSubMerchantID = tmpID[:32]
	}
	if deliveryRange != 0 {
		cook.DeliveryRange = deliveryRange
	}
	if bio != "" {
		cook.Bio = bio
	}
	if len(weekSchedule) >= 7 {
		cook.WeekSchedule = weekSchedule
	}
	if instagramID != "" {
		cook.InstagramID = instagramID
	}
	if twitterID != "" {
		cook.TwitterID = twitterID
	}
	if snapchatID != "" {
		cook.SnapchatID = snapchatID
	}
	err = put(c.ctx, user.ID, cook)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("cannot put cook(%s)", user.ID)
	}
	return cook, nil
}

// FindBySubMerchantID finds a cook by submerchantID
func (c *Client) FindBySubMerchantID(submerchantID string) (string, *Cook, error) {
	if submerchantID == "" {
		return "", nil, errInvalidParameter.WithMessage("submerchantID is invalid.")
	}
	id, cook, err := getBySubmerchantID(c.ctx, submerchantID)
	if err != nil {
		return "", nil, errDatastore.WithError(err).Wrap("failed to get by submerchantID")
	}
	return id, cook, nil
}

// UpdateSubMerchantStatus updates the chef's SubMerchantStatus status
func (c *Client) UpdateSubMerchantStatus(submerchantID, status string) (string, *Cook, error) {
	if submerchantID == "" {
		return "", nil, errInvalidParameter.WithMessage("submerchantID is invalid.")
	}
	id, cook, err := getBySubmerchantID(c.ctx, submerchantID)
	if err != nil {
		return "", nil, errDatastore.WithError(err).Wrap("failed to get by submerchantID")
	}
	cook.SubMerchantStatus = status
	err = put(c.ctx, id, cook)
	if err != nil {
		return "", nil, errDatastore.WithError(err).Wrap("failed to put chef")
	}
	return id, cook, nil
}

// Notify notifies chef with the message
func (c *Client) Notify(id, subject, msg string) error {
	cook, err := get(c.ctx, id)
	if err != nil {
		return errDatastore.WithError(err).Wrap("failed to get by chef")
	}
	messageC := message.New(c.ctx)
	err = messageC.SendSMS(cook.PhoneNumber, msg)
	if err != nil {
		return errors.Wrap("failed to send sms", err)
	}

	return nil
}
