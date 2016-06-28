package gigachef

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/utils"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"gitlab.com/atishpatel/Gigamunch-Backend/core/maps"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/notification"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
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

func (c *Client) UpdateProfile(user *types.User, address *types.Address, phoneNumber, bio string, deliveryRange int32) (*Resp, error) {
	var err error
	chef := new(Gigachef)
	err = get(c.ctx, user.ID, chef)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return nil, errDatastore.WithError(err).Wrap("cannot get gigachef")
	}
	chef.Name = user.Name
	chef.PhotoURL = user.PhotoURL
	chef.Email = user.Email
	chef.ProviderID = user.ProviderID
	chef.DeliveryRange = deliveryRange
	chef.Bio = bio
	if !chef.Application {
		chef.Application = true
		if !appengine.IsDevAppServer() {
			// notify enis
			nC := notification.New(c.ctx)
			err = nC.SendSMS("6153975516", fmt.Sprintf("%s just submit an application. get on that booty. email: %s", user.Name, user.Email))
			if err != nil {
				utils.Criticalf(c.ctx, "failed to notify enis about chef(%s) submitting application", user.ID)
			}
		}
	}
	if chef.DeliveryRange == 0 {
		chef.DeliveryRange = 1
	}
	if address != nil {
		err = maps.GetGeopointFromAddress(c.ctx, address)
		if err != nil {
			return nil, errors.Wrap("failed to GetGeopointFromAddress", err)
		}
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
	if chef.CreatedDatetime.IsZero() {
		chef.CreatedDatetime = time.Now()
	}
	err = put(c.ctx, user.ID, chef)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot put gigachef")
	}
	resp := &Resp{
		ID:       user.ID,
		Gigachef: *chef,
	}
	return resp, nil
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
	if chef.DeliveryRange == 0 {
		chef.DeliveryRange = 5
	}
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
	if chef.CreatedDatetime.IsZero() {
		chef.CreatedDatetime = time.Now()
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

// GetMultiInfo returns Gigachef info
func GetMultiInfo(ctx context.Context, ids []string) ([]Gigachef, error) {
	// TODO add not querying same ids
	chefs, err := getMulti(ctx, ids)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrap("cannot get multi gigachefs")
	}
	return chefs, nil
}

// PostInfoResp contains information related to a post
type PostInfoResp struct {
	Address             types.Address `json:"address"`
	DeliveryRange       int32         `json:"delivery_range"`
	BTSubMerchantID     string        `json:"bt_sub_merchant_id"`
	BTSubMerchantStatus string        `json:"bt_sub_merchant_status"`
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
		Address:             chef.Address,
		DeliveryRange:       chef.DeliveryRange,
		BTSubMerchantID:     chef.BTSubMerchantID,
		BTSubMerchantStatus: chef.SubMerchantStatus,
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
	chef.PayoutMethod = true
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
