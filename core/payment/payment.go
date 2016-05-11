package payment

import (
	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/braintree-go"
	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

var (
	errBT               = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with payment processing."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	btConfig            config.BTConfig
)

// Client is the payment client.  A new client should be created for each context.
type Client struct {
	bt *braintree.Braintree
}

// New returns a new Client. A new client should be created for each context.
func New(ctx context.Context) Client {
	return Client{
		bt: getBTClient(ctx),
	}
}

// GenerateToken generates a token with a customerID
// customerID must be 36 long.
func (c Client) GenerateToken(customerID string) (string, error) {
	if len(customerID) != 36 {
		return "", errInvalidParameter.Wrap("customerID is invalid")
	}
	// check if customer exist
	customerGateway := c.bt.Customer()
	_, err := customerGateway.Find(customerID)
	if err != nil {
		// create customer
		c := &braintree.Customer{
			Id: customerID,
		}
		_, err = customerGateway.Create(c)
		if err != nil {
			return "", errBT.WithError(err).Wrap("cannot create a customer")
		}
	}
	// generate token
	clientToken := c.bt.Transaction().ClientToken()
	token, err := clientToken.GenerateWithCustomer(customerID)
	if err != nil {
		return "", errBT.WithError(err).WithMessage("cannot generate token")
	}
	return token, nil
}

// ReleaseSale release a sale with the SaleID
func (c Client) ReleaseSale(id string) (string, error) {
	t, err := c.bt.Transaction().ReleaseFromEscrow(id)
	if err != nil {
		return "", errBT.WithError(err).Wrap("cannot release from escrow")
	}
	if t.EscrowStatus != braintree.EscrowStatus.ReleasePending && t.EscrowStatus != braintree.EscrowStatus.Released {
		return "", errBT.Wrap("invalid escrow status on release: escrow status: " + t.EscrowStatus)
	}
	return t.Id, nil
}

func (c Client) getTransactionStatus(id string) (string, error) {
	t, err := c.bt.Transaction().Find(id)
	if err != nil {
		return "", errBT.WithError(err)
	}
	return t.Status, nil
}

// RefundSale voids a sale with the SaleID
func (c Client) RefundSale(id string) (string, error) {
	status, err := c.getTransactionStatus(id)
	if err != nil {
		return "", errors.Wrap("cannot find sale", err)
	}
	var t *braintree.Transaction
	if status == "authorized" || status == "submitted_for_settlement" {
		t, err = c.bt.Transaction().Void(id)
	} else {
		t, err = c.bt.Transaction().Refund(id)
	}
	if err != nil {
		return "", errBT.WithError(err)
	}
	return t.Id, nil
}

// MakeSale makes an escrow sale
func (c Client) MakeSale(subMerchantID, nonce string, amount, serviceFee float32) (string, error) {
	t, err := c.bt.Transaction().Create(&braintree.Transaction{
		Type:               "sale",
		MerchantAccountId:  subMerchantID,
		PaymentMethodNonce: nonce,
		Amount:             getBTDecimal(amount),
		ServiceFeeAmount:   getBTDecimal(serviceFee),
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
			HoldInEscrow:        true,
		},
	})
	if err != nil {
		return "", errBT.WithError(err).Wrap("cannot create transaction")
	}
	if t.EscrowStatus != braintree.EscrowStatus.HoldPending && t.EscrowStatus != braintree.EscrowStatus.Held {
		return "", errBT.Wrap("invalid transaction escrow status: status: " + t.Status + " escrow status: " + t.EscrowStatus)
	}
	return t.Id, nil
}

// CreateSubMerchantReq is the request to CreateSubMerchant
// ID must be len 32
type CreateSubMerchantReq struct {
	ID                  string
	FirstName, LastName string
	Email               string
	DateOfBirth         string
	AccountNumber       string
	RoutingNumber       string
	Address             types.Address
}

func (req *CreateSubMerchantReq) valid() error {
	if len(req.ID) != 32 {
		return errInvalidParameter.WithMessage("ID must be length 32.")
	}
	return nil
}

// CreateSubMerchant creates a sub-merchant
// TODO change to update?
func (c Client) CreateSubMerchant(req *CreateSubMerchantReq) (string, error) {
	err := req.valid()
	if err != nil {
		return "", err
	}
	account, err := c.bt.MerchantAccount().Create(&braintree.MerchantAccount{
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
	})
	if err != nil {
		return "", errBT.WithError(err).Wrap("cannot create sub-merchant")
	}
	if account.Status != "pending" {
		return "", errBT.WithMessage("Error creating sub-merchant account with status " + account.Status)
	}
	return account.Id, err
}

func getBTAddress(a types.Address) *braintree.Address {
	return &braintree.Address{
		StreetAddress:   a.Street,
		ExtendedAddress: a.APT,
		Locality:        a.City,
		Region:          a.State,
		PostalCode:      a.Zip,
	}
}

func getBTClient(ctx context.Context) *braintree.Braintree {
	if btConfig.BTMerchantID == "" {
		btConfig = config.GetBTConfig(ctx)
	}
	env := braintree.Production
	if btConfig.BTEnvironment == config.BTSandbox {
		env = braintree.Sandbox
	}
	bt := braintree.New(
		env,
		btConfig.BTMerchantID,
		btConfig.BTPublicKey,
		btConfig.BTPrivateKey,
	)
	bt.HttpClient = urlfetch.Client(ctx)
	return bt
}

func getBTDecimal(v float32) *braintree.Decimal {
	return braintree.NewDecimal(int64(v*100), 2)
}
