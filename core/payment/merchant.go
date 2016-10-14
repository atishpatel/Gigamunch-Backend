package payment

import (
	"fmt"
	"strings"
	"time"

	"github.com/atishpatel/braintree-go"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/gigachef"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
)

const (
	dateOfBirthFormat = "02-01-2006"
)

// GetSubMerchant returns a SubMerchant
func (c Client) GetSubMerchant(subMerchantID string) (*SubMerchantInfo, error) {
	if len(subMerchantID) != 32 {
		return nil, errInvalidParameter.WithMessage("ID must be length 32.")
	}
	ma, err := c.bt.MerchantAccount().Find(subMerchantID)
	if err != nil {
		return nil, errBT.WithError(err).Wrapf("cannot bt.MerchantAccount().Find sub-merchant(%s)", subMerchantID)
	}
	dob, err := time.Parse(dateOfBirthFormat, ma.Individual.DateOfBirth)
	if err != nil {
		return nil, errInternal.WithError(err).Wrapf("failed to parse time from string(%s)", ma.Individual.DateOfBirth)
	}
	sm := &SubMerchantInfo{
		ID:            ma.Id,
		FirstName:     ma.Individual.FirstName,
		LastName:      ma.Individual.LastName,
		Email:         ma.Individual.Email,
		DateOfBirth:   dob,
		AccountNumber: ma.FundingOptions.AccountNumber,
		RoutingNumber: ma.FundingOptions.RoutingNumber,
		Address:       *getAddress(ma.Individual.Address),
	}
	return sm, nil
}

// SubMerchantInfo is the request to UpdateSubMerchant
// ID must be len 32
type SubMerchantInfo struct {
	ID                  string
	FirstName, LastName string
	Email               string
	DateOfBirth         time.Time
	AccountNumber       string
	RoutingNumber       string
	Address             types.Address
}

func (req *SubMerchantInfo) valid() error {
	if len(req.ID) != 32 {
		return errInvalidParameter.WithMessage("ID must be length 32.")
	}
	return nil
}

// UpdateSubMerchant creates or updates sub-merchant info
func (c Client) UpdateSubMerchant(user *types.User, req *SubMerchantInfo) (string, error) {
	err := req.valid()
	if err != nil {
		return "", err
	}
	chefC := gigachef.New(c.ctx)
	return updateSubMerchant(c.ctx, c.bt, chefC, user, req)
}

func updateSubMerchant(ctx context.Context, bt *braintree.Braintree, chefC chefInterface, user *types.User, req *SubMerchantInfo) (string, error) {
	ma := &braintree.MerchantAccount{
		MasterMerchantAccountId: "Gigamunch_marketplace",
		Id:          req.ID,
		TOSAccepted: true,
		Individual: &braintree.MerchantAccountPerson{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			DateOfBirth: req.DateOfBirth.Format(dateOfBirthFormat),
			Address:     getBTAddress(req.Address),
		},
		FundingOptions: &braintree.MerchantAccountFundingOptions{
			Destination:   braintree.FUNDING_DEST_BANK,
			AccountNumber: req.AccountNumber,
			RoutingNumber: req.RoutingNumber,
		},
	}
	_, err := bt.MerchantAccount().Find(req.ID)
	if err == nil {
		ma, err = bt.MerchantAccount().Update(ma)
		if err != nil {
			if strings.Contains(err.Error(), "number is invalid") {
				return "", errInvalidParameter.WithMessage(err.Error())
			}
			return "", errBT.WithError(err).Wrapf("cannot bt.MerchantAccount().Update sub-merchant(%s)", req.ID)
		}
	} else {
		ma, err = bt.MerchantAccount().Create(ma)
		if err != nil {
			if strings.Contains(err.Error(), "number is invalid") {
				return "", errInvalidParameter.WithMessage(err.Error())
			}
			return "", errBT.WithError(err).Wrapf("cannot create sub-merchant(%s)", req.ID)
		}
	}
	if ma.Status != "pending" && ma.Status != "active" {
		return "", errBT.WithMessage("Error creating sub-merchant account with status " + ma.Status)
	}
	if !user.HasSubMerchantID() {
		user.SetSubMerchantID(true)
		err = auth.SaveUser(ctx, user)
		if err != nil {
			return "", errors.Wrap("failed update user to has sub-merchant account", err)
		}
	}
	_, err = chefC.UpdateSubMerchantStatus(ma.Id, ma.Status)
	if err != nil {
		return "", errors.Wrap("failed to chefC.UpdateSubMerchantStatus", err)
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
	if notification == nil || notification.MerchantAccount() == nil {
		return nil, errInternal.Wrapf("payload(%s) \nthere was an error with notification (%+v)\nsubject(%+v) \nMerchantAccount(%+v) \nDisbursement(%+v)", payload, notification, notification.Subject, notification.MerchantAccount(), notification.Disbursement())
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
