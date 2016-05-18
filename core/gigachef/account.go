package gigachef

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/core/notification"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

// Client is a client for gigachef
type Client struct {
	ctx context.Context
}

// New returns a client for gigchef
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Resp is a chef with id
type Resp struct {
	ID string
	Gigachef
}

// Get gets the chef
func (c *Client) Get(id string) (*Resp, error) {
	chef := new(Gigachef)
	err := get(c.ctx, id, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get chef(%s)", id)
	}
	return &Resp{ID: id, Gigachef: *chef}, nil
}

// SaveUserInfo is to save a user's info. Only exposed for account package.
// Please use the account package's func instead of this one.
func SaveUserInfo(ctx context.Context, user *types.User, address *types.Address, phoneNumber string) error {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	chef.Name = user.Name
	chef.PhotoURL = user.PhotoURL
	chef.Email = user.Email
	chef.ProviderID = user.ProviderID
	if address != nil {
		chef.Address = *address
	}
	if phoneNumber != "" {
		chef.PhoneNumber = phoneNumber
	}
	if chef.BTSubMerchantID == "" {
		tmpID := user.ID
		for len(tmpID) <= 32 {
			tmpID += tmpID
		}
		chef.BTSubMerchantID = tmpID[:32]
	}
	err = put(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put gigachef")
	}
	return nil
}

// GetInfo returns nonsensitive Gigachef details
func GetInfo(ctx context.Context, id string) (*Gigachef, error) {
	// TODO switch so it's returns a 'nonsensitive' gigachef info
	chef := new(Gigachef)
	err := get(ctx, id, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	return chef, nil
}

// PostInfoResp contains information related to a post
type PostInfoResp struct {
	Address         types.Address `json:"address"`
	DeliveryRange   int32         `json:"delivery_range"`
	BTSubMerchantID string        `json:"bt_sub_merchant_id"`
}

// GetPostInfo returns info related to a post
// returns: *PostInfoResp, error
func GetPostInfo(ctx context.Context, user *types.User) (*PostInfoResp, error) {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	postInfo := &PostInfoResp{
		Address:         chef.Address,
		DeliveryRange:   chef.DeliveryRange,
		BTSubMerchantID: chef.BTSubMerchantID,
	}
	return postInfo, nil
}

// IncrementNumPost increases NumPosts for a Gigachef by 1
func IncrementNumPost(ctx context.Context, user *types.User) error {
	var err error
	chef := new(Gigachef)
	err = get(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	chef.NumPosts++
	err = put(ctx, user.ID, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("cannot put gigachef")
	}
	return nil
}

// FindBySubMerchantID finds a chef by submerchantID
func (c *Client) FindBySubMerchantID(submerchantID string) (*Resp, error) {
	if submerchantID == "" {
		return nil, errInvalidParameter.WithMessage("submerchantID is invalid.")
	}
	id, chef, err := getBySubmerchantID(c.ctx, submerchantID)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to get by submerchantID")
	}
	resp := &Resp{
		ID:       id,
		Gigachef: *chef,
	}
	return resp, nil
}

// UpdateSubMerchantStatus updates the chef's SubMerchantStatus status
func (c *Client) UpdateSubMerchantStatus(submerchantID, status string) (*Resp, error) {
	if submerchantID == "" {
		return nil, errInvalidParameter.WithMessage("submerchantID is invalid.")
	}
	id, chef, err := getBySubmerchantID(c.ctx, submerchantID)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to get by submerchantID")
	}
	chef.SubMerchantStatus = status
	err = put(c.ctx, id, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("failed to put chef")
	}
	resp := &Resp{
		ID:       id,
		Gigachef: *chef,
	}
	return resp, nil
}

// Notify notifies chef with the message
func (c *Client) Notify(id, subject, message string) error {
	chef := new(Gigachef)
	err := get(c.ctx, id, chef)
	if err != nil {
		return errDatastore.WithError(err).Wrap("failed to get by chef")
	}
	notifcationC := notification.New(c.ctx)
	if !chef.UseEmailOverSMS {
		err = notifcationC.SendSMS(chef.PhoneNumber, message)
		if err == nil {
			return nil
		}
		codeErr := errors.GetErrorWithCode(err)
		if codeErr.Code != errors.CodeInvalidParameter {
			return codeErr.Wrap("failed to send sms")
		}
		message += fmt.Sprintf("\nNote: There was an error while trying to send sms to notify you: %s", codeErr.Message)
	}
	err = notifcationC.SendEmail(chef.Email, subject, message)
	if err != nil {
		return errors.Wrap("failed to send email", err)
	}
	return nil
}
