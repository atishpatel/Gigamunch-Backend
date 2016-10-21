package payment

import (
	"math"

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
	errInternal         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "There was something went wrong."}
	btConfig            config.BTConfig
)

// GetPricePerServing returns the price per serving from cook price per serving.
func GetPricePerServing(cookPPS float32) float32 {
	return float32(math.Ceil(float64(cookPPS) * 1.2))
}

// GetTaxPrice returns the tax percentage.
func GetTaxPrice(pricePerServing float32, lat, long float64) float32 {
	return 7.5 * pricePerServing
}

// Client is the payment client.  A new client should be created for each context.
type Client struct {
	ctx context.Context
	bt  *braintree.Braintree
}

// New returns a new Client. A new client should be created for each context.
func New(ctx context.Context) Client {
	return Client{
		ctx: ctx,
		bt:  getBTClient(ctx),
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
		return "", errBT.WithError(err).Wrapf("cannot release transaction(%d) from escrow", id)
	}
	if t.EscrowStatus != braintree.EscrowStatus.ReleasePending && t.EscrowStatus != braintree.EscrowStatus.Released {
		return "", errBT.Wrapf("invalid escrow status on release: escrow status:%s transactionID: %s", t.EscrowStatus, t.Id)
	}
	return t.Id, nil
}

// CancelRelease cancels release a sale with the SaleID
func (c Client) CancelRelease(id string) (string, error) {
	t, err := c.bt.Transaction().CancelRelease(id)
	if err != nil {
		return "", errBT.WithError(err).Wrapf("cannot cancel release transaction(%d) from escrow", id)
	}
	if t.EscrowStatus != braintree.EscrowStatus.Held && t.EscrowStatus != braintree.EscrowStatus.HoldPending {
		return "", errBT.Wrapf("invalid escrow status on cancel release: escrow status:%s transactionID: %s", t.EscrowStatus, t.Id)
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
	t := &braintree.Transaction{
		Type:               "sale",
		MerchantAccountId:  subMerchantID,
		PaymentMethodNonce: nonce,
		Amount:             getBTDecimal(amount),
		ServiceFeeAmount:   getBTDecimal(serviceFee),
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
			HoldInEscrow:        true,
		},
	}
	t, err := c.bt.Transaction().Create(t)
	if err != nil {
		return "", errBT.WithError(err).Wrapf("cannot create transaction(%#v)", t)
	}
	if t.EscrowStatus != braintree.EscrowStatus.HoldPending && t.EscrowStatus != braintree.EscrowStatus.Held {
		return "", errBT.Wrap("invalid transaction escrow status: status: " + t.Status + " escrow status: " + t.EscrowStatus)
	}
	return t.Id, nil
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

func getAddress(btA *braintree.Address) *types.Address {
	return &types.Address{
		Street: btA.StreetAddress,
		APT:    btA.ExtendedAddress,
		City:   btA.Locality,
		State:  btA.Region,
		Zip:    btA.PostalCode,
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
