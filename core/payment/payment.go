package payment

import (
	"context"
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/core/common"
	"github.com/atishpatel/Gigamunch-Backend/core/logging"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/jmoiron/sqlx"
	braintree "github.com/lionelbarrow/braintree-go"
)

// Errors
var (
	errInternal   = errors.InternalServerError
	errBadRequest = errors.BadRequestError
	errBT         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with payment processing. Your card wasn't charged."}
	btConfig      config.BTConfig
)

// Client is a client for manipulating payments.
type Client struct {
	ctx        context.Context
	log        *logging.Client
	sqlDB      *sqlx.DB
	db         common.DB
	serverInfo *common.ServerInfo
	bt         *braintree.Braintree
}

// NewClient gives you a new client.
func NewClient(ctx context.Context, log *logging.Client, dbC common.DB, sqlC *sqlx.DB, serverInfo *common.ServerInfo) (*Client, error) {
	if log == nil {
		return nil, errInternal.Annotate("failed to get logging client")
	}
	// if sqlC == nil {
	// 	return nil, fmt.Errorf("sqlDB cannot be nil for sub")
	// }
	if dbC == nil {
		return nil, fmt.Errorf("failed to get db")
	}
	if serverInfo == nil {
		return nil, errInternal.Annotate("failed to get server info")
	}
	return &Client{
		ctx:        ctx,
		log:        log,
		sqlDB:      sqlC,
		db:         dbC,
		serverInfo: serverInfo,
		bt:         getBTClient(ctx, serverInfo.IsStandardAppEngine),
	}, nil
}

// CreateCustomer creates a customer.
func (c *Client) CreateCustomer(nonce, email, firstName, lastName string) (string, string, error) {
	customerReq := &braintree.CustomerRequest{
		Email:              email,
		FirstName:          firstName,
		LastName:           lastName,
		PaymentMethodNonce: nonce,
	}
	customer, err := c.bt.Customer().Create(c.ctx, customerReq)
	if err != nil {
		return "", "", errBT.WithError(err).Wrap("failed to bt.Customer.Create")
	}
	return customer.Id, customer.DefaultPaymentMethod().GetToken(), nil
}

// RefundSale voids a sale with the SaleID.
func (c *Client) RefundSale(id string, amount float32) (string, error) {
	status, err := c.getTransactionStatus(id)
	if err != nil {
		return "", errors.Wrap("cannot find sale", err)
	}
	var t *braintree.Transaction
	if status == "authorized" || status == "submitted_for_settlement" {
		t, err = c.bt.Transaction().Void(c.ctx, id)
	} else {
		t, err = c.bt.Transaction().Refund(c.ctx, id, getBTDecimal(amount))
	}
	if err != nil {
		return "", errBT.WithError(err)
	}
	return t.Id, nil
}

func (c *Client) getTransactionStatus(id string) (string, error) {
	t, err := c.bt.Transaction().Find(c.ctx, id)
	if err != nil {
		return "", errBT.WithError(err)
	}
	return string(t.Status), nil
}

func getBTDecimal(v float32) *braintree.Decimal {
	return braintree.NewDecimal(int64(v*100), 2)
}

func getBTClient(ctx context.Context, IsStandardAppEngine bool) *braintree.Braintree {
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
	return bt
}
