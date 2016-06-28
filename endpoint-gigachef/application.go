package gigachef

import (
	"fmt"
	"time"

	"gitlab.com/atishpatel/Gigamunch-Backend/auth"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/payment"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
	"gitlab.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

type Gigachef struct {
	CreatedDatetime   string         `json:"created_datetime"`
	HasCarInsurance   bool           `json:"has_car_insurance"`
	types.UserDetail                 //embedded
	Bio               string         `json:"bio"`
	PhoneNumber       string         `json:"phone_number"`
	Address           Address        `json:"address"`
	AddressTypes      *types.Address `json:"-"`
	DeliveryRange     int            `json:"delivery_range"`
	SendWeeklySummary bool           `json:"send_weekly_summary"`
	UseEmailOverSMS   bool           `json:"use_email_over_sms"`
	gigachef.Rating                  // embedded
	NumPosts          int            `json:"num_posts"`
	NumOrders         int            `json:"num_orders"`
	NumFollowers      int            `json:"num_followers"`
	KitchenPhotoURLs  []string       `json:"kitchen_photo_urls"`
	SubMerchantStatus string         `json:"sub_merchant_status"`
	Application       bool           `json:"application"`
	KitchenInspection bool           `json:"kitchen_inspection"`
	BackgroundCheck   bool           `json:"background_check"`
	FoodHandlerCard   bool           `json:"food_handler_card"`
	PayoutMethod      bool           `json:"payout_method"`
	Verified          bool           `json:"verified"`
}

func (c *Gigachef) set(chef *gigachef.Resp) {
	c.CreatedDatetime = ttos(chef.CreatedDatetime)
	c.HasCarInsurance = chef.HasCarInsurance
	c.UserDetail = chef.UserDetail
	c.Bio = chef.Bio
	c.PhoneNumber = chef.PhoneNumber
	c.Address.set(&chef.Address)
	c.DeliveryRange = int(chef.DeliveryRange)
	c.SendWeeklySummary = chef.SendWeeklySummary
	c.UseEmailOverSMS = chef.UseEmailOverSMS
	c.Rating = chef.Rating
	c.NumPosts = chef.NumPosts
	c.NumOrders = chef.NumOrders
	c.NumFollowers = chef.NumFollowers
	c.KitchenPhotoURLs = chef.KitchenPhotoURLs
	c.SubMerchantStatus = chef.SubMerchantStatus
	c.Application = chef.Application
	c.KitchenInspection = chef.KitchenInspection
	c.BackgroundCheck = chef.BackgroundCheck
	c.FoodHandlerCard = chef.FoodHandlerCard
	c.PayoutMethod = chef.PayoutMethod
	c.Verified = chef.Verified
}

func (c *Gigachef) valid() error {
	if c.Name == "" {
		return fmt.Errorf("Name cannot be empty.")
	}
	if c.Email == "" {
		return fmt.Errorf("Email cannot be empty.")
	}
	if c.PhoneNumber == "" {
		return fmt.Errorf("PhoneNumber cannot be empty.")
	}
	if c.Address.Street == "" {
		return fmt.Errorf("Street cannot be empty.")
	}
	if c.Address.City == "" {
		return fmt.Errorf("City cannot be empty.")
	}
	if c.Address.State == "" {
		return fmt.Errorf("State cannot be empty.")
	}
	if c.Address.Zip == "" {
		return fmt.Errorf("Zip cannot be empty.")
	}
	var err error
	c.AddressTypes, err = c.Address.get()
	return err
}

// UpdateProfileReq is the input request needed for SubmitApplication.
type UpdateProfileReq struct {
	Gigatoken string   `json:"gigatoken"`
	Gigachef  Gigachef `json:"gigachef"`
}

// gigatoken returns the Gigatoken string
func (req *UpdateProfileReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *UpdateProfileReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	err := req.Gigachef.valid()
	if err != nil {
		return err
	}
	return nil
}

// GigachefResp is the output response with a gigachef and error.
type GigachefResp struct {
	Gigachef Gigachef             `json:"gigachef"`
	Err      errors.ErrorWithCode `json:"err"`
}

// UpdateProfile is an endpoint that submits or updates a chef application.
func (service *Service) UpdateProfile(ctx context.Context, req *UpdateProfileReq) (*GigachefResp, error) {
	resp := new(GigachefResp)
	defer handleResp(ctx, "UpdateProfile", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	if req.Gigachef.Address.Country == "" {
		req.Gigachef.Address.Country = "USA"
	}
	if req.Gigachef.PhotoURL != "" {
		user.PhotoURL = req.Gigachef.PhotoURL
	}
	chefC := gigachef.New(ctx)
	chef, err := chefC.UpdateProfile(user, req.Gigachef.AddressTypes, req.Gigachef.PhoneNumber, req.Gigachef.Bio, int32(req.Gigachef.DeliveryRange))
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to update chef profile")
		return resp, nil
	}
	if !user.IsChef() {
		user.SetChef(true)
	}
	err = auth.SaveUser(ctx, user)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to auth.SaveUser")
		return resp, nil
	}
	resp.Gigachef.set(chef)
	return resp, nil
}

// GetGigachef is an endpoint that get the chef info.
func (service *Service) GetGigachef(ctx context.Context, req *GigatokenOnlyReq) (*GigachefResp, error) {
	resp := new(GigachefResp)
	defer handleResp(ctx, "GetGigachef", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefC := gigachef.New(ctx)
	chef, err := chefC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrap("failed to get chef")
		return resp, nil
	}
	resp.Gigachef.set(chef)
	return resp, nil
}

// SubMerchantApplication is the submerchant payout info
type SubMerchantApplication struct {
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Email           string         `json:"email"`
	DateOfBirth     int            `json:"date_of_birth"`
	DateOfBirthTime time.Time      `json:"-"`
	AccountNumber   string         `json:"account_number"`
	RoutingNumber   string         `json:"routing_number"`
	Address         Address        `json:"address"`
	AddressTypes    *types.Address `json:"-"`
}

func (sm *SubMerchantApplication) get(smID string) *payment.SubMerchantInfo {
	return &payment.SubMerchantInfo{
		ID:            smID,
		FirstName:     sm.FirstName,
		LastName:      sm.LastName,
		Email:         sm.Email,
		DateOfBirth:   sm.DateOfBirthTime,
		AccountNumber: sm.AccountNumber,
		RoutingNumber: sm.RoutingNumber,
		Address:       *sm.AddressTypes,
	}
}

func (sm *SubMerchantApplication) set(r *payment.SubMerchantInfo) {
	sm.FirstName = r.FirstName
	sm.LastName = r.LastName
	sm.Email = r.Email
	sm.DateOfBirth = ttoi(r.DateOfBirth)
	sm.AccountNumber = r.AccountNumber
	sm.RoutingNumber = r.RoutingNumber
	sm.Address.set(&r.Address)
}

func (sm *SubMerchantApplication) valid() error {
	if sm.Address.Country == "" {
		sm.Address.Country = "USA"
	}
	var err error
	sm.AddressTypes, err = sm.Address.get()
	return err
}

// UpdateSubMerchantReq updates sub-merchant payment info
type UpdateSubMerchantReq struct {
	Gigatoken   string                 `json:"gigatoken"`
	SubMerchant SubMerchantApplication `json:"sub_merchant"`
}

// gigatoken returns the Gigatoken string
func (req *UpdateSubMerchantReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *UpdateSubMerchantReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	req.SubMerchant.DateOfBirthTime = itot(req.SubMerchant.DateOfBirth)
	err := req.SubMerchant.valid()
	if err != nil {
		return err
	}
	return nil
}

// UpdateSubMerchant creates or updates sub-merchant info
func (service *Service) UpdateSubMerchant(ctx context.Context, req *UpdateSubMerchantReq) (*GigachefResp, error) {
	resp := new(GigachefResp)
	defer handleResp(ctx, "UpdateSubMerchant", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefC := gigachef.New(ctx)
	chef, err := chefC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot get chef(%d)", user.ID)
		return resp, nil
	}
	paymentC := payment.New(ctx)
	utils.Debugf(ctx, "submerch: %#v", req.SubMerchant)
	t, err := req.SubMerchant.Address.get()
	utils.Debugf(ctx, "address type: %#v, err: %s", t, err)
	_, err = paymentC.UpdateSubMerchant(user, req.SubMerchant.get(chef.BTSubMerchantID))
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot update sub-merchant(%d)", chef.BTSubMerchantID)
		return resp, nil
	}
	chef.PayoutMethod = true
	resp.Gigachef.set(chef)
	return resp, nil
}

// GetSubMerchantResp is a resp for GetSubMerchant
type GetSubMerchantResp struct {
	SubMerchant SubMerchantApplication `json:"sub_merchant"`
	Err         errors.ErrorWithCode   `json:"err"`
}

// GetSubMerchant gets a submerchant.
func (service *Service) GetSubMerchant(ctx context.Context, req *GigatokenOnlyReq) (*GetSubMerchantResp, error) {
	resp := new(GetSubMerchantResp)
	defer handleResp(ctx, "GetSubMerchant", resp.Err)
	user, err := validateRequestAndGetUser(ctx, req)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err)
		return resp, nil
	}
	chefC := gigachef.New(ctx)
	chef, err := chefC.Get(user.ID)
	if err != nil {
		resp.Err = errors.GetErrorWithCode(err).Wrapf("cannot get chef(%d)", user.ID)
		return resp, nil
	}
	paymentC := payment.New(ctx)
	sm, err := paymentC.GetSubMerchant(chef.BTSubMerchantID)
	if err != nil {
		resp.SubMerchant.Address.set(&chef.Address)
		resp.SubMerchant.Email = chef.Email
		utils.Infof(ctx, "cannot update sub-merchant(%s): err: %v", chef.BTSubMerchantID, err)
		return resp, nil
	}
	resp.SubMerchant.set(sm)
	return resp, nil
}
