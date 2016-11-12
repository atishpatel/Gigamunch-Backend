package inquiry

import (
	"fmt"
	"math"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/corenew/cook"
	"github.com/atishpatel/Gigamunch-Backend/corenew/eater"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/message"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/corenew/tasks"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	datetimeFormat = "01/02 at 03:04 PM"
)

var (
	errDatastore        = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "Error with datastore."}
	errInvalidParameter = errors.ErrorWithCode{Code: errors.CodeInvalidParameter, Message: "Invalid parameter."}
	errUnauthorized     = errors.ErrorWithCode{Code: errors.CodeUnauthorizedAccess, Message: "User does not have access."}
	errInternal         = errors.ErrorWithCode{Code: errors.CodeInternalServerErr, Message: "There is a problem. Try again in a few minutes."}
	fixedTimeZone       = time.FixedZone("CDT", -6*3600)
	domainURL           string
)

// Client is a client for Inquiry.
type Client struct {
	ctx context.Context
}

// New returns a new Client.
func New(ctx context.Context) *Client {
	getDomainString()
	// if fixedTimeZone == nil {
	// TODO figure out a way to get timezone
	// var err error
	// fixedTimeZone, err =  time.LoadLocation("CDT")
	// if err != nil {
	// 	utils.Criticalf(ctx, "failed to get fixedTimeZone. err: %+v", err)
	// 	fixedTimeZone = time.FixedZone("CDT", -6*3600)
	// }
	// }
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

// SetToFulfilledByTransactionIDs sets Inquiries with the transactionIDs to Fulfilled and queues them to reattempt in 12 hours.
func (c *Client) SetToFulfilledByTransactionIDs(transactionIDs []string) ([]int64, error) {
	var inquiryIDs []int64
	taskC := tasks.New(c.ctx)
	at := time.Now().Add(12 * time.Hour)
	for _, transactionID := range transactionIDs {
		inquiry, err := getByBTTransactionID(c.ctx, transactionID)
		if err != nil {
			return nil, errDatastore.WithError(err).Wrapf("failed to get inquiry with transactionID(%s)", transactionID)
		}
		inquiry.State = State.Fulfilled
		err = put(c.ctx, inquiry.ID, inquiry)
		if err != nil {
			return nil, errDatastore.WithError(err).Wrapf("failed to put Inquiry(%d)", inquiry.ID)
		}
		err = taskC.AddProcessInquiry(inquiry.ID, at)
		if err != nil {
			return nil, errors.Wrap("failed to task.AddProcessInquiry", err)
		}
		inquiryIDs = append(inquiryIDs, inquiry.ID)
	}
	return inquiryIDs, nil
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
	totalPricePerServing := toFixed(pricePerServing*float32(numServings), 2)
	taxPercentage := payment.GetTaxPercentage(cook.Address.Latitude, cook.Address.Longitude)
	totalPrice := toFixed((totalPricePerServing+exchangePrice)*((taxPercentage/100)+1), 2)
	totalTaxPrice := totalPrice - (totalPricePerServing + exchangePrice)
	totalCookPrice := item.CookPricePerServing * float32(numServings)
	totalGigaFee := totalPricePerServing - totalCookPrice
	gigaFeeWithDeliveryAndTax := totalGigaFee + totalTaxPrice
	if !exchangeMethod.Pickup() && !exchangeMethod.CookDelivery() {
		gigaFeeWithDeliveryAndTax += exchangePrice
	}
	// Create Transaction
	paymentC := getPaymentClient(c.ctx)
	transactionID, err := paymentC.StartSale(cook.BTSubMerchantID, nonce, totalPrice, gigaFeeWithDeliveryAndTax)
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
		ExpectedExchangeDateTime: expectedExchangeTime,
		Servings:                 numServings,
		BTTransactionID:          transactionID,
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
	// Put Inquiry in datastore
	_, err = putIncomplete(c.ctx, inquiry)
	if err != nil {
		_, pErr := paymentC.RefundSale(transactionID)
		if pErr != nil {
			utils.Criticalf(c.ctx, "BT Transaction (%s) was not voided! Err: %+v", transactionID, pErr)
		}
		return nil, errDatastore.WithError(err).Wrap("cannot putIncomplete inquiry")
	}
	// Send inquiry bot message
	messageC := getMessageClient(c.ctx)
	cookUI, eaterUI := getCookAndEaterUserInfo(cook.ID, cook.Name, cook.PhotoURL, eater.ID, eater.Name, eater.PhotoURL)
	inquiryI := getInquiryInfo(inquiry)
	messageID, err := messageC.SendInquiryBotMessage(cookUI, eaterUI, inquiryI)
	if err != nil {
		utils.Criticalf(c.ctx, "failed to messageC.SendInquiryBotMessage err: %v", err)
		return inquiry, errors.Wrap("failed to message.SendInquiryBotMessage", err)
	}
	inquiry.MessageID = messageID
	// update inquiry with messageID
	err = put(c.ctx, inquiry.ID, inquiry)
	if err != nil {
		utils.Criticalf(c.ctx, "failed to put inquiry(%d) after sending InquiryBotMessage err: %v", inquiry.ID, err)
		return inquiry, errDatastore.WithError(err).Wrapf("failed to put inquiry(%d) after sending InquiryBotMessage.", inquiry.ID)
	}
	// add task to timeout inquiry at MIN(exchangeTime, now + 12 hours)
	at := time.Now().Add(12 * time.Hour)
	if inquiry.ExpectedExchangeDateTime.Before(at) {
		at = inquiry.ExpectedExchangeDateTime
	}
	tasksC := tasks.New(c.ctx)
	err = tasksC.AddProcessInquiry(inquiry.ID, at)
	if err != nil {
		utils.Criticalf(c.ctx, "failed to tasks.AddProcessInquiry inquiry(%d): %+v", inquiry.ID, err)
	}
	subject := fmt.Sprintf("%s just requested %s", eater.Name, inquiry.Item.Name)
	msg := fmt.Sprintf("%s just requested %d servings of '%s'.\n\nPlease accept or decline the request within 12 hours.\n\n%s/cook/inquiries",
		eater.Name,
		inquiry.Servings,
		inquiry.Item.Name,
		domainURL)
	err = cookC.Notify(cook.ID, subject, msg)
	if err != nil {
		utils.Criticalf(c.ctx, "failed to notify cook(%s ID: %s) about inquiry(%d) err: %+v", cook.Name, cook.ID, inquiry.ID, err)
	}
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
			inquiry.State = State.Accepted
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
		err = sendUpdatedActionMessage(c.ctx, inquiry)
		if err != nil {
			return inquiry, errors.Wrap("failed to sendCookUpdateActionMessage", err)
		}
		// notify Gigamunch if it's a Gigamunch Delivery
		if inquiry.ExchangeMethod.GigamunchDelivery() {
			messageC := getMessageClient(c.ctx)
			err = messageC.SendSMS("9316445311,9316446755,6153975516,6155454989",
				fmt.Sprintf("Time to do a GigaDelivery bras!\n\nItem:%s\n\nDate and time:%s\n\nPickup Location:%s\n\nDropoff Location:%s\n\nInquiryID:%d",
					inquiry.Item.Name,
					inquiry.ExpectedExchangeDateTime.In(fixedTimeZone).Format(datetimeFormat),
					inquiry.ExchangePlanInfo.CookAddress.String(),
					inquiry.ExchangePlanInfo.EaterAddress.String(),
					inquiry.ID,
				))
			if err != nil {
				utils.Criticalf(c.ctx, "failed to notify about GigaDelivery for inquiry(%d). Err: %+v", inquiry.ID, err)
			}
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
		inquiry.State = State.Declined
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
		err = sendUpdatedActionMessage(c.ctx, inquiry)
		if err != nil {
			return inquiry, errors.Wrap("failed to sendCookUpdateActionMessage", err)
		}
	}
	return inquiry, nil
}

func sendUpdatedActionMessage(ctx context.Context, inquiry *Inquiry) error {
	cookC := getCookClient(ctx)
	cookName, cookPhotoURL, err := cookC.GetDisplayInfo(inquiry.CookID)
	if err != nil {
		return errors.Wrap("failed to cook.GetDisplayInfo", err)
	}
	eaterC := getEaterClient(ctx)
	eaterName, eaterPhotoURL, err := eaterC.GetDisplayInfo(inquiry.EaterID)
	if err != nil {
		return errors.Wrap("failed to eater.GetDisplayInfo", err)
	}
	cookUI, eaterUI := getCookAndEaterUserInfo(inquiry.CookID, cookName, cookPhotoURL, inquiry.EaterID, eaterName, eaterPhotoURL)
	inquiryI := getInquiryInfo(inquiry)
	messageC := getMessageClient(ctx)
	err = messageC.UpdateInquiryStatus(inquiry.MessageID, cookUI, eaterUI, inquiryI)
	if err != nil {
		return errors.Wrap("failed to message.UpdateInquiryStatus", err)
	}
	return nil
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
	if inquiry.ExpectedExchangeDateTime.Sub(time.Now()) < time.Duration(12)*time.Hour {
		return nil, errInvalidParameter.WithMessage("The Inquiry can no longer be canceled.")
	}
	inquiry.EaterAction = EaterAction.Canceled
	inquiry.State = State.Canceled
	// submit transaction for refund
	paymentC := getPaymentClient(c.ctx)
	var refundID string
	refundID, err = paymentC.RefundSale(inquiry.BTTransactionID)
	if err != nil {
		return nil, errors.Wrap("failed to payment.SubmitForSettlement", err)
	}
	inquiry.BTRefundTransactionID = refundID
	err = put(c.ctx, id, inquiry)
	if err != nil {
		return nil, errDatastore.WithError(err).Wrapf("failed to put Inquiry(%d) after submitting for settlement", id)
	}
	err = sendUpdatedActionMessage(c.ctx, inquiry)
	if err != nil {
		return inquiry, errors.Wrap("failed to sendCookUpdateActionMessage", err)
	}
	msg := fmt.Sprintf("%s just canceled their request for '%s'.", inquiry.EaterName, inquiry.Item.Name)
	cookC := getCookClient(c.ctx)
	err = cookC.Notify(inquiry.CookID, "Order Canceled", msg)
	if err != nil {
		utils.Criticalf(c.ctx, "Failed to notify cook(%d) about inquiry(%d) cancel. Err: %+v", inquiry.CookID, inquiry.ID, err)
	}
	// notify Gigamunch if it's a Gigamunch Delivery
	if inquiry.ExchangeMethod.GigamunchDelivery() {
		messageC := getMessageClient(c.ctx)
		err = messageC.SendSMS("9316445311,9316446755,6153975516,6155454989",
			fmt.Sprintf("GigaDelivery canceled.\n\nItem:%s\n\nDate and time:%s\n\nInquiryID:%d",
				inquiry.Item.Name,
				inquiry.ExpectedExchangeDateTime.In(fixedTimeZone).Format(datetimeFormat),
				inquiry.ID,
			))
		if err != nil {
			utils.Criticalf(c.ctx, "failed to notify about GigaDelivery for inquiry(%d). Err: %+v", inquiry.ID, err)
		}
	}
	return inquiry, nil
}

// SetReviewID sets the inquiry's ReviewID.
func (c *Client) SetReviewID(id, reviewID int64) error {
	inquiry, err := get(c.ctx, id)
	if err != nil {
		return errDatastore.WithError(err).Wrapf("failed to get inquiry(%d)", id)
	}
	inquiry.ReviewID = reviewID
	err = put(c.ctx, id, inquiry)
	if err != nil {
		return errDatastore.WithError(err).Wrapf("failed to put Inquiry(%d)", id)
	}
	return nil
}

// EaterRequestRefund
// TODO Add request refund feature

// CookAutoDecline

// Process processes the Inquiry.
func (c *Client) Process(id int64) error {
	inquiry, err := get(c.ctx, id)
	if err != nil {
		return errDatastore.WithError(err).Wrapf("failed to get inquiry(%d)", id)
	}
	switch inquiry.State {
	// TODO add RefundRequested
	case State.Fulfilled:
		// If 48 hours after, release payment
		cookC := cook.New(c.ctx)
		var isApprovedSubMerchant bool
		isApprovedSubMerchant, err = cookC.IsSubMerchantApproved(inquiry.CookID)
		if err != nil {
			return errors.Wrap("failed to cook.IsSubMerchantApproved", err)
		}
		now := time.Now()
		before48Hours := time.Since(inquiry.ExpectedExchangeDateTime) < (48 * time.Hour)
		if inquiry.Issue || before48Hours || !isApprovedSubMerchant {
			// process in min(12 hours, 48 hours after exchange)
			at := now.Add(12 * time.Hour)
			if before48Hours {
				at = inquiry.ExpectedExchangeDateTime.Add(48 * time.Hour)
			}
			taskC := tasks.New(c.ctx)
			err = taskC.AddProcessInquiry(id, at)
			if err != nil {
				return errors.Wrap("failed to task.AddProcessInquiry", err)
			}
			return nil
		}
		// release sale
		paymentC := getPaymentClient(c.ctx)
		err = paymentC.ReleaseSale(inquiry.BTTransactionID)
		if err != nil {
			return errors.Wrap(fmt.Sprintf("failed to payment.ReleaseSale for inquiry(%d)", id), err)
		}
		inquiry.State = State.Paid
		err = put(c.ctx, id, inquiry)
		if err != nil {
			return errDatastore.WithError(err).Wrapf("failed to put inquiry(%d)", id)
		}
	case State.Accepted:
		// If after ExchangeTime, set to Fulfilled state.
		taskC := tasks.New(c.ctx)
		now := time.Now()
		if now.Before(inquiry.ExpectedExchangeDateTime) {
			err = taskC.AddProcessInquiry(id, inquiry.ExpectedExchangeDateTime)
			if err != nil {
				return errors.Wrap("faield to task.AddProcessInquiry", err)
			}
			return nil
		}
		inquiry.State = State.Fulfilled
		err = put(c.ctx, id, inquiry)
		if err != nil {
			return errDatastore.WithError(err).Wrapf("failed to put inquiry(%d)", id)
		}
		// Process in 48 hours
		err = taskC.AddProcessInquiry(id, now.Add(48*time.Hour))
		if err != nil {
			return errors.Wrap("faield to task.AddProcessInquiry", err)
		}
	case State.Pending:
		// If it's 12 hours past exchange time, timeout request.
		if time.Since(inquiry.CreatedDateTime) < (12*time.Hour) && time.Now().Before(inquiry.ExpectedExchangeDateTime) {
			taskC := tasks.New(c.ctx)
			err = taskC.AddProcessInquiry(id, inquiry.ExpectedExchangeDateTime)
			if err != nil {
				return errors.Wrap("faield to task.AddProcessInquiry", err)
			}
			return nil
		}
		// timeout order
		inquiry.State = State.TimedOut
		err = put(c.ctx, id, inquiry)
		if err != nil {
			return errDatastore.WithError(err).Wrapf("failed to put inquiry(%d)", id)
		}
		err = sendUpdatedActionMessage(c.ctx, inquiry)
		if err != nil {
			return errors.Wrap("failed to sendCookUpdateActionMessage", err)
		}
	case State.Declined:
		fallthrough
	case State.TimedOut:
		fallthrough
	case State.Canceled:
		fallthrough
	case State.Refunded:
		fallthrough
	case State.Paid:
		break
	default:
		return errInternal.Wrapf("Unknown state(%s) for Inquiry(%d) while processing", inquiry.State, id)
	}
	return nil
}

type messageClient interface {
	SendInquiryBotMessage(cookUI *message.UserInfo, eaterUI *message.UserInfo, inquiryI *message.InquiryInfo) (string, error)
	UpdateChannel(cookUI *message.UserInfo, eaterUI *message.UserInfo, inquiryI *message.InquiryInfo) error
	UpdateInquiryStatus(messageID string, cookUI *message.UserInfo, eaterUI *message.UserInfo, inquiryI *message.InquiryInfo) error
	SendSMS(to, message string) error
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
	ReleaseSale(string) error
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
	GetDisplayInfo(id string) (string, string, error)
}

var getCookClient = func(ctx context.Context) cookClient {
	return cook.New(ctx)
}

type cookClient interface {
	Get(id string) (*cook.Cook, error)
	GetDisplayInfo(id string) (string, string, error)
	Notify(id, subject, msg string) error
}

func getCookAndEaterUserInfo(cookID, cookName, cookPhotoURL, eaterID, eaterName, eaterPhotoURL string) (*message.UserInfo, *message.UserInfo) {
	cookUI := &message.UserInfo{
		ID:    cookID,
		Name:  cookName,
		Image: cookPhotoURL,
	}
	eaterUI := &message.UserInfo{
		ID:    eaterID,
		Name:  eaterName,
		Image: eaterPhotoURL,
	}
	return cookUI, eaterUI
}

func getInquiryInfo(inquiry *Inquiry) *message.InquiryInfo {
	var photoURL string
	if len(inquiry.Item.Photos) > 0 {
		photoURL = inquiry.Item.Photos[0]
	}
	return &message.InquiryInfo{
		ID:           inquiry.ID,
		State:        inquiry.State,
		CookAction:   inquiry.CookAction,
		EaterAction:  inquiry.EaterAction,
		ItemID:       inquiry.ItemID,
		ItemName:     inquiry.Item.Name,
		ItemImage:    photoURL,
		Price:        inquiry.PaymentInfo.TotalPrice,
		IsDelivery:   inquiry.ExchangeMethod.Delivery(),
		Servings:     inquiry.Servings,
		ExchangeTime: inquiry.ExpectedExchangeDateTime,
	}
}

func round(num float32) int {
	return int(num + float32(math.Copysign(0.5, float64(num))))
}

func toFixed(num float32, precision int) float32 {
	output := float32(math.Pow(10, float64(precision)))
	return float32(round(num*output)) / output
}

func getDomainString() {
	if domainURL == "" {
		domainURL = os.Getenv("DOMAIN_URL")
	}
}
