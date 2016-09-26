package cook

import (
	"golang.org/x/net/context"
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/notification"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
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

// Get gets the cook.
func (c *Client) Get(id string) (*Cook, error) {
	cook, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get cook(%s)", id)
	}
	return cook, nil
}

// Update updates a cook's info.
func (c *Client) Update(user *types.User, address *types.Address, phoneNumber, bio string,
	deliveryRange int32, weekSchedule []WeekSchedule, instagramID, twitterID, snapchatID string) (*Cook, error) {
	cook, err := get(c.ctx, user.ID)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, errDatastore.WithError(err).Wrapf("failed to get cook(%s)", user.ID)
	}
	if cook.CreatedDatetime.IsZero() {
		cook.CreatedDatetime = time.Now()
	}
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
func (c *Client) Notify(id, subject, message string) error {
	cook, err := get(c.ctx, id)
	if err != nil {
		return errDatastore.WithError(err).Wrap("failed to get by chef")
	}
	notifcationC := notification.New(c.ctx)
	err = notifcationC.SendSMS(cook.PhoneNumber, message)
	if err != nil {
		return errors.Wrap("failed to send sms", err)
	}

	return nil
}
