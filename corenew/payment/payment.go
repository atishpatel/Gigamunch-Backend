package payment

import (
	"math"
	"math/rand"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/braintree-go"
	"golang.org/x/net/context"

	"strings"

	"google.golang.org/appengine/urlfetch"
)

var (
	errBT                  = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with payment processing. Your card wasn't charged."}
	errInvalidParameter    = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errInternal            = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "There was something went wrong."}
	btConfig               config.BTConfig
	replaceCharactersArray = [...]string{"!", "#", "$", "%", "&", "'", "*", "+", "/", "=", "?", "^", "`", "{", "|", "}", "~", ".", "@"}
)

// GetPricePerServing returns the price per serving from cook price per serving.
func GetPricePerServing(cookPPS float32) float32 {
	return float32(math.Ceil(float64(cookPPS) * 1.2))
}

// GetTaxPercentage returns the tax percentage.
func GetTaxPercentage(lat, long float64) float32 {
	return 7.5
}

// Client is the payment client.  A new client should be created for each context.
type Client struct {
	ctx context.Context
	bt  *braintree.Braintree
}

// New returns a new Client. A new client should be created for each context.
func New(ctx context.Context) *Client {
	return &Client{
		ctx: ctx,
		bt:  getBTClient(ctx),
	}
}

// GenerateToken generates a token with a customerID
// customerID must be 36 long.
func (c *Client) GenerateToken(customerID string) (string, error) {
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
func (c *Client) ReleaseSale(id string) error {
	t, err := c.bt.Transaction().ReleaseFromEscrow(id)
	if err != nil {
		return errBT.WithError(err).Wrapf("cannot release transaction(%d) from escrow", id)
	}
	if t.EscrowStatus != braintree.EscrowStatus.ReleasePending && t.EscrowStatus != braintree.EscrowStatus.Released {
		return errBT.Wrapf("invalid escrow status on release: escrow status:%s transactionID: %s", t.EscrowStatus, t.Id)
	}
	return nil
}

// CancelRelease cancels release a sale with the SaleID
func (c *Client) CancelRelease(id string) (string, error) {
	t, err := c.bt.Transaction().CancelRelease(id)
	if err != nil {
		return "", errBT.WithError(err).Wrapf("cannot cancel release transaction(%d) from escrow", id)
	}
	if t.EscrowStatus != braintree.EscrowStatus.Held && t.EscrowStatus != braintree.EscrowStatus.HoldPending {
		return "", errBT.Wrapf("invalid escrow status on cancel release: escrow status:%s transactionID: %s", t.EscrowStatus, t.Id)
	}
	return t.Id, nil
}

func (c *Client) getTransactionStatus(id string) (string, error) {
	t, err := c.bt.Transaction().Find(id)
	if err != nil {
		return "", errBT.WithError(err)
	}
	return t.Status, nil
}

// RefundSale voids a sale with the SaleID
func (c *Client) RefundSale(id string) (string, error) {
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

// SubmitForSettlement submits a Sale for settlement with the SaleID.
func (c *Client) SubmitForSettlement(id string) error {
	_, err := c.bt.Transaction().SubmitForSettlement(id)
	if err != nil {
		return errBT.WithError(err)
	}
	return nil
}

// GigamunchToSubmerchant sends money form Gigamunch to a submerchant.
func (c *Client) GigamunchToSubmerchant(subMerchantID string, amount float32) (string, error) {
	t := &braintree.Transaction{
		Type:              "sale",
		MerchantAccountId: subMerchantID,
		CreditCard: &braintree.CreditCard{
			CardholderName: btConfig.CompanyCard.CardholderName,
			Number:         btConfig.CompanyCard.Number,
			ExpirationDate: btConfig.CompanyCard.ExpirationDate,
			CVV:            btConfig.CompanyCard.CVV,
		},
		Amount:           getBTDecimal(amount),
		ServiceFeeAmount: getBTDecimal(0),
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
		},
	}
	t, err := c.bt.Transaction().Create(t)
	if err != nil {
		return "", errBT.WithError(err).Wrapf("cannot create transaction(%#v)", t)
	}
	return t.Id, nil
}

// StartSubscriptionReq is the reqest for StartSubscription.
type StartSubscriptionReq struct {
	CustomerID string
	Nonce      string
	PlanID     string
	StartDate  time.Time
}

func (s *StartSubscriptionReq) valid() error {
	if len(s.CustomerID) < 32 {
		return errInvalidParameter.WithMessage("Invalid CustomerID")
	}
	if s.Nonce == "" {
		return errInvalidParameter.WithMessage("Invalid Payment Info.").Wrap("Payment Nonce cannot be empty.")
	}
	if s.PlanID == "" {
		return errInvalidParameter.WithMessage("Invalid Box Plan.").Wrap("PlanID cannot be empty")
	}
	return nil
}

// StartSubscription creates a subscription with a nonce from a customer.
func (c *Client) StartSubscription(req *StartSubscriptionReq) (string, error) {
	if req == nil {
		return "", errInvalidParameter.Wrap("StarySubscription is nil.")
	}
	err := req.valid()
	if err != nil {
		return "", err
	}
	s := &braintree.Subscription{
		Id:                 req.CustomerID[:25] + randStringBytes(7) + "_sub",
		PlanId:             req.PlanID,
		PaymentMethodNonce: req.Nonce,
	}
	if !req.StartDate.IsZero() {
		s.FirstBillingDate = req.StartDate.Format("2006-01-02")
	}
	s, err = c.bt.Subscription().Create(s)
	if err != nil {
		return "", errBT.WithError(err).Wrapf("cannot create subscription(%#v)", s)
	}
	return s.Id, nil
}

// StartSale makes an escrow sale
func (c *Client) StartSale(subMerchantID, nonce string, amount, serviceFee float32) (string, error) {
	t := &braintree.Transaction{
		Type:               "sale",
		MerchantAccountId:  subMerchantID,
		PaymentMethodNonce: nonce,
		Amount:             getBTDecimal(amount),
		ServiceFeeAmount:   getBTDecimal(serviceFee),
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: false,
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

func GetIDFromEmail(email string) string {
	s := email
	for _, v := range replaceCharactersArray {
		s = strings.Replace(s, v, "", -1)
	}
	for len(s) < 36 {
		s += s
	}
	return s[:36]
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
