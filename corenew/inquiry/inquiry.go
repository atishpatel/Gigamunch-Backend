package inquiry

import (
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errUnauthorized     = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access."}
)

// Client is a client for orders
type Client struct {
	ctx context.Context
}

// New returns a new Client
func New(ctx context.Context) *Client {
	return &Client{ctx: ctx}
}

// Get gets the inquiry by ID.
func (c *Client) Get(user *types.User, id int64) (*Inquiry, error) {
	inquiry, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get inquiry(%d)", id)
	}
	if !user.IsAdmin() && inquiry.CookID != user.ID && inquiry.EaterID != user.ID {
		return nil, errUnauthorized.WithMessage("User does not have access to Inquiry.")
	}
	return inquiry, nil
}

// GetByCookID gets inquiries by CookID.
func (c *Client) GetByCookID(cookID string, startLimit, endLimit int) ([]*Inquiry, error) {
	inquiries, err := getCookInquirys(c.ctx, cookID, startLimit, endLimit)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get inquiries with cookID(%s)", cookID)
	}
	return inquiries, nil
}

// GetByEaterID gets inquiries by EaterID.
func (c *Client) GetByEaterID(eaterID string, startLimit, endLimit int) ([]*Inquiry, error) {
	inquiries, err := getEaterInquirys(c.ctx, eaterID, startLimit, endLimit)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get inquiries with eaterID(%s)", eaterID)
	}
	return inquiries, nil
}

// GetByTransactionID gets the inquiry by TransactionID.
func (c *Client) GetByTransactionID(transactionID string) (*Inquiry, error) {
	inquiry, err := getByBTTransactionID(c.ctx, transactionID)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get inquiry with transactionID(%s)", transactionID)
	}
	return inquiry, nil
}

func validateMakeParams(itemID int64, nonce string, eaterID string, eaterAddress *types.Address, numServings int32, exchangeMethod types.ExchangeMethod) error {
	if itemID == 0 {
		return errInvalidParameter.WithMessage("ItemID cannot be 0.")
	}
	if nonce == "" {
		return errInvalidParameter.WithMessage("Nonce cannot be empty.")
	}
	if eaterID == "" {
		return errInvalidParameter.WithMessage("EaterID cannot be empty.")
	}
	if eaterAddress == nil {
		return errInvalidParameter.WithMessage("EaterAddress cannot be empty.")
	}
	if numServings == 0 {
		return errInvalidParameter.WithMessage("NumServings cannot be empty.")
	}
	if !exchangeMethod.Valid() {
		return errInvalidParameter.WithMessage("Exchange method is not valid.")
	}
	return nil
}

// Make makes a new Inquiry.
func (c *Client) Make(itemID int64, nonce string, eaterID string, eaterAddress *types.Address, numServings int32, exchangeMethod types.ExchangeMethod, expectedExchangeTime time.Time) (*Inquiry, error) {
	if err := validateMakeParams(itemID, nonce, eaterID, eaterAddress, numServings, exchangeMethod); err != nil {
		return nil, err
	}
	// get item
	itemC := getItemClient(c.ctx)
	item, err := itemC.Get(itemID)
	if err != nil {
		return nil, errors.Wrap("failed to itemC.Get", err)
	}
	// get eater
	eaterC := getEaterClient(c.ctx)
	eater, err := eaterC.Get(eaterID)
	if err != nil {
		return nil, errors.Wrap("failed to eaterC.Get", err)
	}
	// get cook
	cookC := getCookClient(c.ctx)
	cook, err := cookC.Get(item.CookID)
	if err != nil {
		return nil, errors.Wrap("failed to cookC.Get", err)
	}
	ems := types.GetExchangeMethods(cook.Address.GeoPoint, cook.DeliveryRange, cook.DeliveryPrice, eaterAddress.GeoPoint)
	found := false
	var exchangePrice float32
	for _, v := range ems {
		if v.Equal(exchangeMethod) {
			exchangePrice = v.Price
			found = true
		}
	}
	if !found {
		return nil, errInvalidParameter.WithMessage("The selected exchange method is not available.")
	}
	// calculate pricing
	pricePerServing := payment.GetPricePerServing(item.CookPricePerServing)
	totalPricePerServing := pricePerServing * float32(numServings)
	taxPercentage := payment.GetTaxPercentage(cook.Address.Latitude, cook.Address.Longitude)
	totalPrice := (totalPricePerServing + exchangePrice) * (taxPercentage / 100)
	totalTaxPrice := totalPrice - (totalPricePerServing + exchangePrice)
	totalCookPrice := item.CookPricePerServing * float32(numServings)
	totalGigaFee := totalPricePerServing - totalCookPrice
	cookPriceWithDelivery := totalCookPrice
	gigaFeeWithDelivery := totalGigaFee
	if exchangeMethod.Pickup() || exchangeMethod.CookDelivery() {
		cookPriceWithDelivery += exchangePrice
	} else {
		gigaFeeWithDelivery += exchangePrice
	}
	// start payment
	paymentC := getPaymentClient(c.ctx)
	transactionID, err := paymentC.StartSale(cook.BTSubMerchantID, nonce, cookPriceWithDelivery, gigaFeeWithDelivery)
	if err != nil {
		return nil, errors.Wrap("failed to paymentC.StartSale", err)
	}
	// get cheap estimate for distance and duration
	distance := eaterAddress.GreatCircleDistance(cook.Address.GeoPoint)
	duration := eaterAddress.EstimatedDuration(cook.Address.GeoPoint)
	inquiry := &Inquiry{
		CreatedDateTime: time.Now(),
		CookID:          item.CookID,
		EaterID:         eaterID,
		EaterPhotoURL:   eater.PhotoURL,
		EaterName:       eater.Name,
		ItemID:          itemID,
		Item: ItemInfo{
			Name:               item.Title,
			Description:        item.Description,
			Photos:             item.Photos,
			Ingredients:        item.Ingredients,
			DietaryConcerns:    item.DietaryConcerns,
			ServingDescription: item.ServingDescription,
		},
		State:                    State.Pending,
		EaterAction:              EaterAction.Accepted,
		CookAction:               CookAction.Pending,
		BTTransactionID:          transactionID,
		ExpectedExchangeDateTime: expectedExchangeTime,
		Servings:                 numServings,
		PaymentInfo: PaymentInfo{
			CookPricePerServing: item.CookPricePerServing,
			PricePerServing:     pricePerServing,
			CookPrice:           totalCookPrice,
			ExchangePrice:       exchangePrice,
			TaxPrice:            totalTaxPrice,
			ServiceFee:          totalGigaFee,
			TotalPrice:          totalPrice,
		},
		ExchangeMethod: exchangeMethod,
		ExchangePlanInfo: ExchangePlanInfo{
			EaterAddress: *eaterAddress,
			CookAddress:  cook.Address,
			Distance:     distance,
			Duration:     duration,
		},
	}
	_, err = putIncomplete(c.ctx, inquiry)
	if err != nil {
		_, pErr := paymentC.RefundSale(transactionID)
		if pErr != nil {
			utils.Criticalf(c.ctx, "BT Transaction (%s) was not voided! Err: %+v", transactionID, pErr)
		}
		return nil, errDatastore.WithError(err).Wrap("cannot putIncomplete inquiry")
	}
	// Start inquiry bot message
	messageC := getMessageClient(c.ctx)
	cookUI := &message.UserInfo{
		ID:    cook.ID,
		Name:  cook.Name,
		Image: cook.PhotoURL,
	}
	eaterUI := &message.UserInfo{
		ID:    eater.ID,
		Name:  eater.Name,
		Image: eater.PhotoURL,
	}
	var photoURL string
	if len(inquiry.Item.Photos) > 0 {
		photoURL = inquiry.Item.Photos[0]
	}
	inquiryI := &message.InquiryInfo{
		ID:          inquiry.ID,
		State:       inquiry.State,
		CookAction:  inquiry.CookAction,
		EaterAction: inquiry.EaterAction,
		ItemID:      inquiry.ItemID,
		ItemName:    inquiry.Item.Name,
		ItemImage:   photoURL,
	}

	err = messageC.SendInquiryBotMessage(cookUI, eaterUI, inquiryI)
	if err != nil {
		utils.Criticalf(c.ctx, "failed to messageC.SendInquiryBotMessage err: %v", err)
		return inquiry, errors.Wrap("failed to message.SendInquiryBotMessage", err)
	}
	// TODO add task to where cook has (exchangeTime || 12 hours) to reply
	return inquiry, nil
}

// CookAccept updates the Inquiry to Cook Accepted.
func (c *Client) CookAccept(user *types.User, id int64) (*Inquiry, error) {
	inquiry, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get Inquiry(%d)", id)
	}
	if !user.IsAdmin() && user.ID != inquiry.CookID {
		return nil, errUnauthorized.WithMessage("You are not the cook of the Inquiry.")
	}
	if inquiry.State != State.Pending {
		return nil, errInvalidParameter.WithMessage("Inquiry is no long in a pending state.")
	}
	if inquiry.CookAction != CookAction.Accepted {
		inquiry.CookAction = CookAction.Accepted
		if inquiry.EaterAction == EaterAction.Accepted {
			// submit transaction for settlement
			paymentC := getPaymentClient(c.ctx)
			err = paymentC.SubmitForSettlement(inquiry.BTTransactionID)
			if err != nil {
				return nil, errors.Wrap("failed to payment.SubmitForSettlement", err)
			}
			// TODO notify eater and send twilio status message to channel
		}
		err = put(c.ctx, id, inquiry)
		if err != nil {
			return nil, errDatastore.WithError(err).Wrapf("failed to put Inquiry(%d) after submitting for settlement", id)
		}
	}
	return inquiry, nil
}

// CookDecline updates the Inquiry to Cook Declined.
func (c *Client) CookDecline(user *types.User, id int64) (*Inquiry, error) {
	inquiry, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get Inquiry(%d)", id)
	}
	if !user.IsAdmin() && user.ID != inquiry.CookID {
		return nil, errUnauthorized.WithMessage("You are not the cook of the Inquiry.")
	}
	if inquiry.State != State.Pending {
		return nil, errInvalidParameter.WithMessage("Inquiry is no long in a pending state.")
	}
	if inquiry.CookAction != CookAction.Declined {
		inquiry.CookAction = CookAction.Declined
		// submit transaction for refund
		paymentC := getPaymentClient(c.ctx)
		var refundID string
		refundID, err = paymentC.RefundSale(inquiry.BTTransactionID)
		if err != nil {
			return nil, errors.Wrap("failed to payment.SubmitForSettlement", err)
		}
		inquiry.BTRefundTransactionID = refundID
		// TODO notify eater and send twilio status message to channel
		err = put(c.ctx, id, inquiry)
		if err != nil {
			return nil, errDatastore.WithError(err).Wrapf("failed to put Inquiry(%d) after submitting for settlement", id)
		}
	}
	return inquiry, nil
}

// CookCancel
// func (c *Client) CookCancel() error {

// 	return nil
// }

// EaterCancel updates the Inquiry to Eater Canceled.
func (c *Client) EaterCancel(user *types.User, id int64) (*Inquiry, error) {
	inquiry, err := get(c.ctx, id)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to get Inquiry(%d)", id)
	}
	if !user.IsAdmin() && user.ID != inquiry.EaterID {
		return nil, errUnauthorized.WithMessage("You are not part of the Inquiry.")
	}
	if inquiry.State != State.Pending {
		return nil, errInvalidParameter.WithMessage("Inquiry is no long in a pending state.")
	}
	if inquiry.ExpectedExchangeDateTime.Sub(time.Now()) > time.Duration(12)*time.Hour {
		return nil, errInvalidParameter.WithMessage("The Inquiry can no longer be canceled.")
	}
	inquiry.State = State.Refunded
	inquiry.EaterAction = EaterAction.Canceled
	// submit transaction for refund
	paymentC := getPaymentClient(c.ctx)
	var refundID string
	refundID, err = paymentC.RefundSale(inquiry.BTTransactionID)
	if err != nil {
		return nil, errors.Wrap("failed to payment.SubmitForSettlement", err)
	}
	inquiry.BTRefundTransactionID = refundID
	// TODO notify eater and send twilio status message to channel
	err = put(c.ctx, id, inquiry)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to put Inquiry(%d) after submitting for settlement", id)
	}
	return inquiry, nil
}

// EaterRequestRefund
// TODO Add request refund feature

// CookAutoDecline

// Process

type messageClient interface {
	SendInquiryBotMessage(cookUI *message.UserInfo, eaterUI *message.UserInfo, inquiryI *message.InquiryInfo) error
	UpdateChannel(cookUI *message.UserInfo, eaterUI *message.UserInfo, inquiryI *message.InquiryInfo) error
}

var getMessageClient = func(ctx context.Context) messageClient {
	return message.New(ctx)
}

var getPaymentClient = func(ctx context.Context) paymentClient {
	return payment.New(ctx)
}

type paymentClient interface {
	StartSale(string, string, float32, float32) (string, error)
	SubmitForSettlement(string) error
	ReleaseSale(string) (string, error)
	RefundSale(string) (string, error)
	CancelRelease(string) (string, error)
}

var getItemClient = func(ctx context.Context) itemClient {
	return item.New(ctx)
}

type itemClient interface {
	Get(int64) (*item.Item, error)
}

var getEaterClient = func(ctx context.Context) eaterClient {
	return eater.New(ctx)
}

type eaterClient interface {
	Get(string) (*eater.Eater, error)
}

var getCookClient = func(ctx context.Context) cookClient {
	return cook.New(ctx)
}

type cookClient interface {
	Get(string) (*cook.Cook, error)
}
