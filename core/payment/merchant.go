package payment

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/atishpatel/braintree-go"
	"gitlab.com/atishpatel/Gigamunch-Backend/auth"
	"gitlab.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"gitlab.com/atishpatel/Gigamunch-Backend/errors"
	"gitlab.com/atishpatel/Gigamunch-Backend/types"
)

// UpdateSubMerchantReq is the request to UpdateSubMerchant
// ID must be len 32
type UpdateSubMerchantReq struct {
	User                types.User
	ID                  string
	FirstName, LastName string
	Email               string
	DateOfBirth         string
	AccountNumber       string
	RoutingNumber       string
	Address             types.Address
}

func (req *UpdateSubMerchantReq) valid() error {
	if len(req.ID) != 32 {
		return errInvalidParameter.WithMessage("ID must be length 32.")
	}
	return nil
}

// UpdateSubMerchant creates or updates sub-merchant info
func (c Client) UpdateSubMerchant(req *UpdateSubMerchantReq) (string, error) {
	err := req.valid()
	if err != nil {
		return "", err
	}

	ma := &braintree.MerchantAccount{
		MasterMerchantAccountId: "Gigamunch_marketplace",
		Id:          req.ID,
		TOSAccepted: true,
		Individual: &braintree.MerchantAccountPerson{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			DateOfBirth: req.DateOfBirth,
			Address:     getBTAddress(req.Address),
		},
		FundingOptions: &braintree.MerchantAccountFundingOptions{
			Destination:   braintree.FUNDING_DEST_BANK,
			AccountNumber: req.AccountNumber,
			RoutingNumber: req.RoutingNumber,
		},
	}
	// TODO switch to find?
	if req.User.HasSubMerchantID() {
		ma, err = c.bt.MerchantAccount().Update(ma)
		if err != nil {
			return "", errBT.WithError(err).Wrap("cannot bt.MerchantAccount().Update")
		}
	} else {
		ma, err = c.bt.MerchantAccount().Create(ma)
		if err != nil {
			return "", errBT.WithError(err).Wrap("cannot create sub-merchant")
		}
		req.User.SetSubMerchantID(true)
		err = auth.SaveUser(c.ctx, &req.User)
		if err != nil {
			return "", errors.Wrap("failed update user to has sub-merchant account", err)
		}
	}
	if ma.Status != "pending" {
		return "", errBT.WithMessage("Error creating sub-merchant account with status " + ma.Status)
	}
	return ma.Id, err
}

type chefInterface interface {
	FindBySubMerchantID(string) (*gigachef.Resp, error)
	UpdateSubMerchantStatus(string, string) (*gigachef.Resp, error)
	Notify(string, string, string) error
}

// DisbursementException handles a disbursement exception notification
func (c Client) DisbursementException(signature, payload string) ([]string, error) {
	chefC := gigachef.New(c.ctx)
	return disbursementException(c.ctx, signature, payload, c.bt, chefC)
}

func disbursementException(ctx context.Context, signature, payload string, bt *braintree.Braintree, chefC chefInterface) ([]string, error) {
	notification, err := parseNotification(ctx, signature, payload, braintree.DisbursementExceptionWebhook, bt)
	if err != nil {
		return nil, err
	}
	chef, err := chefC.FindBySubMerchantID(notification.MerchantAccount().Id)
	if err != nil {
		return nil, errors.Wrap("failed to find chef by submerchantID", err)
	}
	disbursement := notification.Disbursement()
	message := fmt.Sprintf("A transaction to your account failed because '%s' please take the following action: '%s'", disbursement.ExceptionMessage, disbursement.FollowUpAction)
	err = chefC.Notify(chef.ID, "There was a problem sending money to you! - Gigamunch", message)
	if err != nil {
		return nil, errors.Wrap("failed to notify chef", err)
	}
	return disbursement.TransactionIds, nil
}

// SubMerchantApproved handles a submerchant approved notification
func (c Client) SubMerchantApproved(signature, payload string) error {
	chefC := gigachef.New(c.ctx)
	return subMerchantApproved(c.ctx, signature, payload, c.bt, chefC)
}

func subMerchantApproved(ctx context.Context, signature, payload string, bt *braintree.Braintree, chefC chefInterface) error {
	notification, err := parseNotification(ctx, signature, payload, braintree.SubMerchantAccountApprovedWebhook, bt)
	if err != nil {
		return err
	}
	merch := notification.MerchantAccount()
	_, err = chefC.UpdateSubMerchantStatus(merch.Id, merch.Status)
	if err != nil {
		return errors.Wrap(fmt.Sprintf("failed to update submerchant(%s) status(%s)", merch.Id, merch.Status), err)
	}
	return nil
}

// SubMerchantDeclined handles a submerchant declined notification
func (c Client) SubMerchantDeclined(signature, payload string) error {
	chefC := gigachef.New(c.ctx)
	return subMerchantDeclined(c.ctx, signature, payload, c.bt, chefC)
}

func subMerchantDeclined(ctx context.Context, signature, payload string, bt *braintree.Braintree, chefC chefInterface) error {
	notification, err := parseNotification(ctx, signature, payload, braintree.SubMerchantAccountDeclinedWebhook, bt)
	if err != nil {
		return err
	}
	merch := notification.MerchantAccount()
	chef, err := chefC.UpdateSubMerchantStatus(merch.Id, merch.Status)
	if err != nil {
		return errors.Wrap(fmt.Sprintf("failed to update submerchant(%s) status(%s)", merch.Id, merch.Status), err)
	}
	err = chefC.Notify(chef.ID, "There was a problem with the approving your bank info - Gigamunch", notification.Subject.APIErrorResponse.ErrorMessage)
	if err != nil {
		return errors.Wrap("failed to notify chef", err)
	}
	return nil
}

func parseNotification(ctx context.Context, signature, payload, wantKind string, bt *braintree.Braintree) (*braintree.WebhookNotification, error) {
	notification, err := bt.WebhookNotification().Parse(signature, payload)
	if err != nil {
		return nil, errBT.WithError(err).Wrap("cannot parse payment notification")
	}

	if notification.Kind != wantKind {
		return nil, errInvalidParameter.Wrapf("%#v notification is not a %s notification", notification, wantKind)
	}
	return notification, nil
}
